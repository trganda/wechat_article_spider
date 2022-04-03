package config

import (
	"time"
	"wechat_crawer/utils"

	"github.com/spf13/viper"
)

var (
	GITHUB     string = ""
	TAG        string = ""
	GOVERSIONI string = "1.17.3"

	ChromeDriver       string = "verdors/chromedriver"
	SeleniumServerPath string = "verdors/selenium-server-4.1.3.jar"
	Port               int    = 9515

	Query    string    = "每日安全动态推送"
	FakeId   string    = "MzA5NDYyNDI0MA=="
	TimeLine time.Time = time.Now()
)

var Cfg *utils.Config

func InitConfig(path string) error {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	if utils.FileExist(path) {
		err := viper.ReadInConfig()
		if err != nil {
			return err
		}

		viper.Unmarshal(&Cfg)
	}

	return nil
}
