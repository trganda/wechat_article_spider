package utils

type Config struct {
	WebDriver       Driver      `yaml:"webdriver"`
	AppMsgQueryArgs AppMsgQuery `yaml:"appmsg"`
}

type Driver struct {
	ChromeDriver   string `yaml:"chromedriver"`
	SeleniumServer string `yaml:"seleniumserver"`
	Port           int    `yaml:"port"`
}

type AppMsgQuery struct {
	Query    string `yaml:"query"`
	FakeId   string `yaml:"fakeid"`
	TimeLine string `yaml:"timeline"`
}

// Parameter struct of request
type AppMsgArgs struct {
	Token  string `json:"token"` // csrf token
	Lang   string `json:"lang"`  // language
	F      string `json:"f"`     // format
	Ajax   string `json:"ajax"`  // request type
	Action string `json:"list_ex"`
	Begin  string `json:"begin"`  // begin
	Count  string `json:"count"`  // number pre-request
	Query  string `json:"query"`  // query condition
	FakeId string `json:"fakeid"` // id
	Type   string `json:"type"`   // 9
}

// A brief struct for wechat article
type AppMsgListItem struct {
	Aid        string `json:"aid"`
	AlbumId    string `json:"album_id"`
	AppmsgId   string `json:"appmsgid"`
	Checking   string `json:"checking"`
	Cover      string `json:"cover"`
	CreateTime string `json:"create_time"`
	Digest     string `json:"digest"`
	ItemIdx    string `json:"itemidx"`
	Link       string `json:"link"`
	Title      string `json:"title"`
	UpdateTime string `json:"update_time"`
}

type BaseResp struct {
	ErrMsg string `json:"err_msg"`
	Ret    int    `json:"ret"`
}

type AppMsg struct {
	AppMsgCnt  int              `json:"app_msg_cnt"`
	AppMsgList []AppMsgListItem `json:"app_msg_list"`
	Resp       BaseResp         `json:"base_resp"`
}
