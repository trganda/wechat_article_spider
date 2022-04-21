package crawler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
	"wechat_crawler/config"
	"wechat_crawler/utils"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	goQrcode "github.com/skip2/go-qrcode"
)

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

func Login() (utils.Cookies, utils.AppMsgArgs, error) {

	ctx, cancel := chromedp.NewExecAllocator(
		context.Background(),

		// set headless false
		append(
			chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", config.Cfg.ChromeDP.Headless),
			chromedp.UserAgent(config.Cfg.ChromeDP.Headers.UserAgent),
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
	qrDone := make(chan bool)
	loginDone := make(chan bool)

	// this will be used to capture the request id for matching network events
	var qrRequestID network.RequestID
	var homeRequestID network.RequestID

	// csrf token
	var token string

	// set the download url as the chromedp github user avatar
	homeUrl := "https://mp.weixin.qq.com/"
	qrUrl, _ := url.Parse("https://mp.weixin.qq.com/cgi-bin/scanloginqrcode?action=getqrcode&random=1649474933069")
	loginUrl, _ := url.Parse("https://mp.weixin.qq.com/cgi-bin/home?t=home/index&lang=zh_CN&token=8")

	// set up a listener to watch the network events and close the channel when
	// complete the request id matching is important both to filter out
	// unwanted network events and to reference the downloaded file later
	chromedp.ListenTarget(ctx, func(v interface{}) {
		switch ev := v.(type) {
		case *network.EventRequestWillBeSent:
			log.Printf("EventRequestWillBeSent: %v: %v", ev.RequestID, ev.Request.URL)
			currentUrl, _ := url.Parse(ev.Request.URL)
			if currentUrl.Host == qrUrl.Host && currentUrl.Query().Get("action") == qrUrl.Query().Get("action") {
				qrRequestID = ev.RequestID
			} else if currentUrl.Host == qrUrl.Host &&
				currentUrl.Path == loginUrl.Path &&
				currentUrl.Query().Get("t") == loginUrl.Query().Get("t") {
				homeRequestID = ev.RequestID
				token = currentUrl.Query().Get("token")
			}
		case *network.EventLoadingFinished:
			log.Printf("EventLoadingFinished: %v", ev.RequestID)
			if ev.RequestID == qrRequestID {
				close(qrDone)
			} else if ev.RequestID == homeRequestID {
				close(loginDone)
			}
		}
	})

	// navigate to the home url
	if err := chromedp.Run(ctx,
		chromedp.Navigate(homeUrl),
	); err != nil {
		log.Fatal(err)
	}

	// This will block until the qr image has been requested
	<-qrDone

	if err := chromedp.Run(ctx,
		getQRCode(qrRequestID),
	); err != nil {
		log.Fatal(err)
	}

	// This will block until login
	<-loginDone

	var cookies utils.Cookies
	if err := chromedp.Run(ctx,
		waitLogin(&cookies),
	); err != nil {
		log.Fatal(err)
	}

	appmsgArgs := utils.AppMsgArgs{
		Token:  token,
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

	return cookies, appmsgArgs, nil
}

func getQRCode(requestID network.RequestID) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		buf, err := network.GetResponseBody(requestID).Do(ctx)
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile("qrcode.jpg", buf, 0644); err != nil {
			return err
		}
		log.Println("download to qrcode.jpg.")

		PrintQRCode(buf)

		return nil
	}
}

func waitLogin(cookies *utils.Cookies) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		urlCookies, _ := network.GetCookies().Do(ctx)
		cookiesData := network.GetAllCookiesReturns{Cookies: urlCookies}

		for idx := 0; idx < len(cookiesData.Cookies); idx++ {
			curCookie := cookiesData.Cookies[idx]
			cookie := utils.Cookie{
				Name:   curCookie.Name,
				Value:  curCookie.Value,
				Path:   curCookie.Path,
				Domain: curCookie.Domain,
				Secure: curCookie.Secure,
				Expiry: uint(curCookie.Expires),
			}
			cookies.Cookies = append(cookies.Cookies, &cookie)
		}

		return nil
	}
}

func PrintQRCode(buf []byte) {
	img, s, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		log.Fatal(err)
	}

	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		log.Fatal(err)
	}

	res, err := qrcode.NewQRCodeReader().Decode(bmp, nil)
	if err != nil {
		log.Fatal(err)
	}

	qr, err := goQrcode.New(res.String(), goQrcode.High)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(qr.ToSmallString(false))

	log.Print(s)
}

func CrawArticlewithCondition(cookies utils.Cookies,
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

func CrawArticle(cookies utils.Cookies, getArgs utils.AppMsgArgs) ([]byte, utils.Cookies) {
	client := &http.Client{}

	targetUrl := "https://mp.weixin.qq.com/cgi-bin/appmsg"

	// Create a self defined request
	request, _ := http.NewRequest("GET", targetUrl, nil)
	request.Header.Set("User-Agent", config.Cfg.ChromeDP.Headers.UserAgent)
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
		ctCookie := utils.ConvertToCookie(respCookies[idx])
		if oldIdx > -1 {
			*updatedCookies.Cookies[oldIdx] = ctCookie
		} else {
			updatedCookies.Cookies = append(updatedCookies.Cookies, &ctCookie)
		}
	}

	return body, updatedCookies
}
