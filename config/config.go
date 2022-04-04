package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
	"wechat_crawer/utils"
)

var (
	GITHUB     string = ""
	TAG        string = ""
	GOVERSIONI string = "1.17.3"

	ChromeDriver       string = "vendors/chromedriver"
	SeleniumServerPath string = "vendors/selenium-server-4.1.3.jar"
	Port               int    = 9515

	Query    string    = "每日安全动态推送"
	FakeId   string    = "MzA5NDYyNDI0MA=="
	TimeLine time.Time = time.Now()

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
	}

	return nil
}
