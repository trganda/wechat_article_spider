package crawler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"strings"
	"wechat_crawler/config"
	"wechat_crawler/utils"

	"github.com/PuerkitoBio/goquery"

	"github.com/gocolly/colly/v2"
)

func DumpItem(item utils.AppMsgListItem, rootPrefix string) {
	DumpPage(item.Link, item.Title, rootPrefix)
}

func DumpPage(urlPath string, title string, rootPrefix string) {
	// [scheme:][//[userinfo@]host][/]path[?query][#fragment]
	parsedUrl, err := url.Parse(urlPath)
	if err != nil {
		log.Fatalf("parse url %s failed, please check it. error: %s\n", urlPath, err.Error())
	}

	// replace the '/' or '\' of title with '-'
	title = strings.TrimSpace(title)
	title = strings.ReplaceAll(title, "/", "-")
	title = strings.ReplaceAll(title, `\`, "-")
	pathPrefix := path.Join(rootPrefix, title)

	// Mkdir with title
	if !utils.FileExist(pathPrefix) {
		err := os.Mkdir(pathPrefix, os.ModePerm)
		if err != nil {
			log.Fatalf("mkdir data failed. error: %s\n", err.Error())
			return
		}
	}

	fmt.Println(parsedUrl.Path)

	c := colly.NewCollector(
		colly.UserAgent(config.UserAgent),
	)

	// Find and visit all links
	c.OnHTML("a[href]", func(element *colly.HTMLElement) {
		href := element.Attr("href")

		if strings.HasSuffix(href, ".js") ||
			strings.HasSuffix(href, ".css") || strings.HasSuffix(href, ".jpeg") {
			element.Request.Visit(href)
		}
	})

	c.OnHTML("script[src]", func(element *colly.HTMLElement) {
		href := element.Attr("src")
		if strings.HasSuffix(href, ".js") ||
			strings.HasSuffix(href, ".css") || strings.HasSuffix(href, ".jpeg") {
			element.Request.Visit(href)
		}
	})

	c.OnHTML("link[href]", func(element *colly.HTMLElement) {
		href := element.Attr("href")

		if strings.HasSuffix(href, ".js") ||
			strings.HasSuffix(href, ".css") || strings.HasSuffix(href, ".jpeg") {
			element.Request.Visit(href)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		// Check the main page
		if r.Request.URL.Path == parsedUrl.Path && r.Request.URL.Host == parsedUrl.Host &&
			r.Request.URL.RawQuery == parsedUrl.RawQuery && strings.Contains(r.Headers.Get("Content-Type"), "text/html") {

			// Load the HTML document
			bodyReader := bytes.NewReader(r.Body)
			document, err := goquery.NewDocumentFromReader(bodyReader)
			if err != nil {
				log.Printf("load document from %s failed. error: %s\n", r.Request.URL.String(), err.Error())
				log.Println("skipped.")
				err = r.Save(path.Join(pathPrefix, title+".html"))
				if err != nil {
					return
				}
			} else {
				// Modify the path to relative path
				document.Find("a[href],link[href]").Each(func(i int, selection *goquery.Selection) {
					href, _ := selection.Attr("href")

					t, _ := url.Parse(href)
					selection.SetAttr("href", strings.TrimLeft(t.Path, "/"))
				})

				document.Find("script[src]").Each(func(i int, selection *goquery.Selection) {
					href, _ := selection.Attr("src")

					t, _ := url.Parse(href)
					selection.SetAttr("src", strings.TrimLeft(t.Path, "/"))
				})

				html, err := document.Html()
				if err != nil {
					return
				}
				err = ioutil.WriteFile(path.Join(pathPrefix, title+".html"), []byte(html), 0644)
				if err != nil {
					log.Fatalf("dump %s failed, error %s", path.Join(pathPrefix, title+".html"), err)
				}
			}

			return
		}

		// Create directories for each url
		requestUrl := r.Request.URL

		urlPath := path.Clean(path.Dir(requestUrl.Path))
		urlFileName := path.Base(requestUrl.Path)

		err := os.MkdirAll(path.Join(pathPrefix, urlPath), os.ModePerm)
		if err != nil {
			return
		}

		err = r.Save(path.Join(pathPrefix, urlPath, urlFileName))
		if err != nil {
			return
		}
	})

	err = c.Visit(urlPath)
	if err != nil {
		return
	}
}
