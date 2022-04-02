package crawer

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
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
	const (
		// These paths will be different on your system.
		seleniumPath     = "vendors/selenium-server-4.1.3.jar"
		googleDriverPath = "vendors/chromedriver.exe"
		port             = 9515
	)

	selenium.SetDebug(true)

	opts := []selenium.ServiceOption{}

	service, err := selenium.NewChromeDriverService(googleDriverPath, port, opts...)
	if err != nil {
		return nil, utils.AppMsgArgs{}, err
	}
	defer service.Stop()

	args := []string{
		"--user-agent=Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.60 Safari/537.36",
	}

	chromeCaps := chrome.Capabilities{Args: args}

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "chrome"}
	caps.AddChrome(chromeCaps)

	wd, _ := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))

	if err := wd.Get("https://mp.weixin.qq.com/"); err != nil {
		panic(err)
	}

	wd.Wait(waitCondition)

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

	// buf := strings.NewReader(para.Encode())

	// Create a self defined request
	request, _ := http.NewRequest("GET", targetUrl, nil)
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.60 Safari/537.36")
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
