package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"time"
	"wechat_crawer/utils"
)

var (
	GITHUB     string = "https://github.com/trganda/wechat_article_spider"
	TAG        string = ""
	GOVERSIONI string = "1.18"

	ChromeDriver       string = "vendors/chromedriver"
	SeleniumServerPath string = "vendors/selenium-server-4.1.3.jar"
	Port               int    = 9515

	UserAgent string = "Mozilla/5.0 (Windows NT 6.1; Win64; x64) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.60 Safari/537.36"

	Query    string = "每日安全动态推送"
	FakeId   string = "MzA5NDYyNDI0MA=="
	TimeLine string = time.Now().Format(TimeFormat)

	TimeFormat string = "2006-01-02T15:04:05"
)

var Cfg *utils.Config

func InitConfig(path string) error {
	if utils.FileExist(path) {
		yamlFile, _ := ioutil.ReadFile(path)
		err := yaml.Unmarshal(yamlFile, &Cfg)
		if err != nil {
			return err
		}
	} else {
		err := ioutil.WriteFile("config.yaml", configToYaml(), 0644)
		if err != nil {
			return err
		}
		log.Println("the initialization config file has been generated.")
		tCfg := CreateDefaultConfig()
		Cfg = &tCfg
	}

	return nil
}

func CreateDefaultConfig() utils.Config {
	return utils.Config{
		WebDriver: utils.Driver{
			ChromeDriver:   ChromeDriver,
			SeleniumServer: SeleniumServerPath,
			Port:           Port,
			Headers: utils.BrowserHeaders{
				UserAgent: UserAgent,
			},
		},
		AppMsgQueryArgs: utils.AppMsgQuery{
			Query:    Query,
			FakeId:   FakeId,
			TimeLine: TimeLine,
		},
	}
}

func configToYaml() []byte {
	buf, err := yaml.Marshal(CreateDefaultConfig())
	if err != nil {
		log.Fatalf("convert default config to yaml failed. error: %s\n", err.Error())
	}
	return buf
}
