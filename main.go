package main

import (
	"encoding/json"
	"fmt"
	"github.com/tebeka/selenium"
	"io/ioutil"
	"os"
	"wechat_crawer/config"
	"wechat_crawer/crawer"
	"wechat_crawer/utils"
)

func main() {
	config.InitConfig("config.yaml")

	var cookies []selenium.Cookie
	var urlArgs utils.AppMsgArgs
	var err error

	if !utils.FileExist("data/cookies.json") || !utils.FileExist("data/urlargs.json") {

		err = os.MkdirAll("data", os.ModePerm)
		if err != nil {
			return
		}

		cookies, urlArgs, err = crawer.Login()
		if err != nil {
			return
		}

		// Save cookies and urlArgs as json file
		jsonCookies, err := json.MarshalIndent(cookies, "", "\t")
		jsonurlArgs, err := json.MarshalIndent(urlArgs, "", "\t")

		err = ioutil.WriteFile("data/cookies.json", jsonCookies, 0644)
		if err != nil {
			return
		}
		err = ioutil.WriteFile("data/urlargs.json", jsonurlArgs, 0644)
		if err != nil {
			return
		}
	} else {
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
	}

	ret := crawer.CrawArticle(cookies, urlArgs)
	fmt.Println(ret)
}
