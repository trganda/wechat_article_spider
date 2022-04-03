package crawer

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"wechat_crawer/config"
	"wechat_crawer/utils"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

func waitCondition(wd selenium.WebDriver) (bool, error) {
	title, err := wd.Title()
	if title == "公众号" {
		return true, nil
	}
	return false, err
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
		Query:  "",
		FakeId: "MzA5NDYyNDI0MA==",
		Type:   "9",
	}

	cookies, err := wd.GetCookies()
	if err != nil {
		return nil, utils.AppMsgArgs{}, err
	}

	return cookies, appmsgArgs, err
}

func CrawArticle(cookies []selenium.Cookie, getArgs utils.AppMsgArgs) string {
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

	httpCookies := utils.ConvertToHttpCookie(cookies)

	for idx := 0; idx < len(httpCookies); idx++ {
		request.AddCookie(&httpCookies[idx])
	}

	resp, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(body)
}
