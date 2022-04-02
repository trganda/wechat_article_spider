package main

import (
	"fmt"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

type AppMsgArgs struct {
	Token  string `json:"token"`
	Lang   string `json:"lang"`
	F      string `json:"f"`
	Ajax   string `json:"ajax"`
	Action string `jsong:"list_ex"`
	Begin  string `json:"begin"`
	Count  string `json:"count"`
	Query  string `json:"query"`
	FakeId string `json:"fakeid"`
	Type   string `json:"type"`
}

func waitCondition(wd selenium.WebDriver) (bool, error) {
	title, err := wd.Title()
	if title == "公众号" {
		return true, nil
	}
	return false, err
}

func login() (map[string]string, error) {

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
		return nil, err
	}
	defer service.Stop()

	args := []string{
		"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.55 Safari/537.36",
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
	materialElem, err := wd.FindElement(selenium.ByXPATH, "//*[@id=\"js_level2_title\"]/li[1]/ul/li[1]/a")
	if err != nil {
		panic(err)
	}

	if err := materialElem.Click(); err != nil {
		panic(err)
	}

	time.Sleep(3000)

	createElem, err := wd.FindElement(selenium.ByXPATH, "//*[@id=\"js_main\"]/div[3]/div[2]/div/div/div/div[1]/div/div/div[2]/ul/li[1]")
	if err != nil {
		panic(err)
	}

	if err := createElem.Click(); err != nil {
		panic(err)
	}

	cookies, err := wd.GetCookies()
	if err != nil {
		return nil, err
	}

	core_cookie := make(map[string]string)
	// ret := AppMsgArgs
	for idx := 0; idx < len(cookies); idx++ {
		cookie := cookies[idx]
		core_cookie[cookie.Name] = cookie.Value
	}

	return core_cookie, err
}

func main() {
	cookies, err := login()
	if err != nil {
		return
	}
	print(cookies)
}
