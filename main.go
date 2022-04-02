package main

import (
	"fmt"

	"wechat_crawer/crawer"
)

func main() {
	cookies, urlArgs, err := crawer.Login()
	if err != nil {
		return
	}

	fmt.Println(cookies)
	fmt.Println(urlArgs)
}
