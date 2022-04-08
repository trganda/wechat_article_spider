package crawer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"wechat_crawer/config"
	"wechat_crawer/utils"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

// this will be used to capture the request id for matching network events
var requestID network.RequestID

func waitCondition(wd selenium.WebDriver) (bool, error) {
	title, err := wd.Title()
	if title == "公众号" {
		return true, nil
	}
	return false, err
}

// FilterCondition Filter the msg item before *timeline*
func FilterCondition(item utils.AppMsgListItem) (bool, error) {
	createTime := time.Unix(item.CreateTime, 0)
	filterTime, err := time.ParseInLocation(config.TimeFormat,
		config.Cfg.AppMsgQueryArgs.TimeLine, time.Local)

	if err != nil {
		log.Fatalf("convert configed timeline failed. error: %s\n", err.Error())
		return false, err
	}

	if createTime.After(filterTime) || createTime.Equal(filterTime) {
		return true, nil
	}

	return false, nil
}

func DefaultFilterCondition(item utils.AppMsgListItem) (bool, error) {
	return true, nil
}

func Login() ([]selenium.Cookie, utils.AppMsgArgs, error) {

	// Start a Selenium WebDriver server instance (if one is not already
	// running).
	var (
		// These paths will be different on your system.
		googleDriverPath = config.Cfg.WebDriver.ChromeDriver
		port             = config.Cfg.WebDriver.Port
	)

	selenium.SetDebug(true)

	var opts []selenium.ServiceOption

	service, err := selenium.NewChromeDriverService(googleDriverPath, port, opts...)
	if err != nil {
		return nil, utils.AppMsgArgs{}, err
	}
	defer service.Stop()

	args := []string{
		"--user-agent=" + config.Cfg.WebDriver.Headers.UserAgent,
	}

	chromeCaps := chrome.Capabilities{Args: args}

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "chrome"}
	caps.AddChrome(chromeCaps)

	wd, _ := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))

	if err := wd.Get("https://mp.weixin.qq.com/"); err != nil {
		panic(err)
	}

	err = wd.Wait(waitCondition)
	if err != nil {
		return nil, utils.AppMsgArgs{}, err
	}

	time.Sleep(3000)

	// Catch csrf token from url
	currentURL, _ := wd.CurrentURL()
	urlArgs, _ := url.ParseQuery(strings.Split(currentURL, "?")[1])

	appmsgArgs := utils.AppMsgArgs{
		Token:  urlArgs["token"][0],
		Lang:   "zh_CN",
		F:      "json",
		Ajax:   "1",
		Action: "list_ex",
		Begin:  "0",
		Count:  "10",
		Query:  config.Cfg.AppMsgQueryArgs.Query,
		FakeId: config.Cfg.AppMsgQueryArgs.FakeId,
		Type:   "9",
	}

	cookies, err := wd.GetCookies()
	if err != nil {
		return nil, utils.AppMsgArgs{}, err
	}

	return cookies, appmsgArgs, err
}

func Logining() {

	ctx, cancel := chromedp.NewExecAllocator(
		context.Background(),

		// set headless false
		append(
			chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", false),
		)...,
	)
	defer cancel()

	// create context
	ctx, cancel = chromedp.NewContext(
		ctx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// create a timeout as a safety net to prevent any infinite wait loops
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// set up a channel so we can block later while we monitor the download
	// progress
	done := make(chan bool)
	//domUpdated := make(chan bool)

	// set the download url as the chromedp github user avatar
	urlstr := "https://mp.weixin.qq.com/"

	// set up a listener to watch the network events and close the channel when
	// complete the request id matching is important both to filter out
	// unwanted network events and to reference the downloaded file later
	chromedp.ListenTarget(ctx, func(v interface{}) {
		switch ev := v.(type) {
		case *network.EventRequestWillBeSent:
			log.Printf("EventRequestWillBeSent: %v: %v", ev.RequestID, ev.Request.URL)
			if ev.Request.URL == urlstr {
				requestID = ev.RequestID
			}
		case *network.EventLoadingFinished:
			log.Printf("EventLoadingFinished: %v", ev.RequestID)
			if ev.RequestID == requestID {
				close(done)
			}
		case *dom.EventDocumentUpdated:
			log.Printf("EventDocumentUpdated")
		}
	})

	// all we need to do here is navigate to the download url
	if err := chromedp.Run(ctx,
		chromedp.Navigate(urlstr),
		getQRCode(),
	); err != nil {
		log.Fatal(err)
	}

	// This will block until the chromedp listener closes the channel
	<-done

	if err := chromedp.Run(ctx, getQRCode()); err != nil {
		log.Fatal(err)
	}

	// get the downloaded bytes for the request id
	//var buf []byte
	//if err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
	//	var err error
	//	buf, err = network.GetResponseBody(requestID).Do(ctx)
	//	return err
	//})); err != nil {
	//	log.Fatal(err)
	//}

	// write the file to disk - since we hold the bytes we dictate the name and
	// location
	//if err := ioutil.WriteFile("download.png", buf, 0644); err != nil {
	//	log.Fatal(err)
	//}
	//log.Print("wrote download.png")
}

func getQRCode() chromedp.ActionFunc {
	return func(ctx context.Context) error {
		var qrCodeURL string
		var status bool

		chromedp.WaitVisible("#header > div.banner > div > div > " +
			"div.login__type__container.login__type__container__scan > img")
		// get location of qr image
		chromedp.AttributeValue("#header > div.banner > div > div > "+
			"div.login__type__container.login__type__container__scan > img", "src",
			&qrCodeURL, &status)

		if status == false {
			log.Println("get attribute (src) of qr img failed.")
			return nil
		}

		// download the image file
		chromedp.Navigate(qrCodeURL)

		buf, err := network.GetResponseBody(requestID).Do(ctx)
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile("qrcode.png", buf, 0644); err != nil {
			return err
		}
		log.Println("download to qrcode.png.")

		return nil
	}

}

func CrawArticlewithCondition(cookies []selenium.Cookie,
	getArgs utils.AppMsgArgs, condition utils.Condition) utils.AppMsgListItems {

	var ret bool
	var err error

	var appMsg utils.AppMsg
	var appMsgList utils.AppMsgListItems

	for true {
		time.Sleep(utils.RandDuration())
		jsonData, respCookies := CrawArticle(cookies, getArgs)
		cookies = respCookies
		// Forward search
		getArgs.Begin = getArgs.Begin + getArgs.Count

		err = json.Unmarshal(jsonData, &appMsg)
		if err != nil {
			log.Fatalf("unmarshal the response json data to utils.AppMsg failed. error: %s\n", err.Error())
			return appMsgList
		}

		if appMsg.Resp.Ret != 0 {
			log.Printf("response with error: %s\n", appMsg.Resp.ErrMsg)
			return appMsgList
		}

		for idx := 0; idx < len(appMsg.AppMsgList); idx++ {
			ret, err = condition(appMsg.AppMsgList[idx])
			if err != nil || !ret {
				break
			} else if ret {
				appMsgList.Items = append(appMsgList.Items, appMsg.AppMsgList[idx])
			}
		}

		if err != nil || !ret {
			break
		}
	}

	return appMsgList
}

func CrawArticle(cookies []selenium.Cookie, getArgs utils.AppMsgArgs) ([]byte, []selenium.Cookie) {
	client := &http.Client{}

	targetUrl := "https://mp.weixin.qq.com/cgi-bin/appmsg"

	// Create a self defined request
	request, _ := http.NewRequest("GET", targetUrl, nil)
	request.Header.Set("User-Agent", config.Cfg.WebDriver.Headers.UserAgent)
	request.Header.Set("X-Requested-With", "XMLHttpRequest")

	// Convert getArgs to url
	para := request.URL.Query()
	para.Add("action", getArgs.Action)
	para.Add("begin", getArgs.Begin)
	para.Add("count", getArgs.Count)
	para.Add("fakeid", getArgs.FakeId)
	para.Add("type", getArgs.Type)
	para.Add("query", getArgs.Query)
	para.Add("token", getArgs.Token)
	para.Add("lang", getArgs.Lang)
	para.Add("f", getArgs.F)
	para.Add("ajax", getArgs.Ajax)

	request.URL.RawQuery = para.Encode()

	httpCookies := utils.ConvertToHttpCookies(cookies)

	for idx := 0; idx < len(httpCookies); idx++ {
		request.AddCookie(&httpCookies[idx])
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalf("sending rqeust error: %s\n", err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("reading response body error: %s\n", err.Error())
	}

	respCookies := resp.Cookies()
	updatedCookies := cookies

	// Update cookies with response cookies
	for idx := 0; idx < len(respCookies); idx++ {
		oldIdx := utils.IdxofCookieswithName(cookies, respCookies[idx].Name)
		if oldIdx > -1 {
			updatedCookies[oldIdx] = utils.ConvertToSeleniumCookie(respCookies[idx])
		} else {
			updatedCookies = append(updatedCookies, utils.ConvertToSeleniumCookie(respCookies[idx]))
		}
	}

	return body, updatedCookies
}
