package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"wechat_crawer/crawer"
	"wechat_crawer/utils"

	"github.com/tebeka/selenium"
)

func main() {
	var cookies []selenium.Cookie
	var urlArgs utils.AppMsgArgs

	if !utils.FileExist("data/cookies.json") || !utils.FileExist("data/urlargs.json") {
		cookies, urlArgs, err := crawer.Login()
		if err != nil {
			return
		}

		// Save cookies and urlArgs as json file
		jsonCookies, _ := json.MarshalIndent(cookies, "", "\t")
		jsonurlArgs, _ := json.MarshalIndent(urlArgs, "", "\t")

		ioutil.WriteFile("data/cookies.json", jsonCookies, 0644)
		ioutil.WriteFile("data/urlargs.json", jsonurlArgs, 0644)
	}

	buf, _ := ioutil.ReadFile("data/cookies.json")

	err := json.Unmarshal(buf, &cookies)
	if err != nil {
		fmt.Println(err)
	}

	buf, _ = ioutil.ReadFile("data/urlargs.json")

	err = json.Unmarshal(buf, &urlArgs)
	if err != nil {
		fmt.Println(err)
	}

	crawer.CrawArticle(cookies, urlArgs)

}
