package utils

type Config struct {
	ChromeDP        Driver      `yaml:"chromedp"`
	AppMsgQueryArgs AppMsgQuery `yaml:"appmsg"`
}

type Driver struct {
	Headless bool           `yaml:"headless"`
	Headers  BrowserHeaders `yaml:"headers"`
}

type BrowserHeaders struct {
	UserAgent string `yaml:"user-agent"`
}

type AppMsgQuery struct {
	Query      string `yaml:"query"`
	FakeId     string `yaml:"fakeid"`
	TimeLine   string `yaml:"timeline"`
	DumpFormat string `yaml:"dumpformat"`
}

// AppMsgArgs Parameter struct of request
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

// AppMsgListItem A brief struct for wechat article
type AppMsgListItem struct {
	Aid        string `json:"aid"`
	AlbumId    string `json:"album_id"`
	AppmsgId   uint64 `json:"appmsgid"`
	Checking   uint   `json:"checking"`
	Cover      string `json:"cover"`
	CreateTime int64  `json:"create_time"`
	Digest     string `json:"digest"`
	ItemIdx    uint   `json:"itemidx"`
	Link       string `json:"link"`
	Title      string `json:"title"`
	UpdateTime int64  `json:"update_time"`
}

type AppMsgListItems struct {
	Items []AppMsgListItem `json:"app_msg_list"`
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

type Cookies struct {
	Cookies []*Cookie `json:"cookies"`
}

type Cookie struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Path   string `json:"path"`
	Domain string `json:"domain"`
	Secure bool   `json:"secure"`
	Expiry uint   `json:"expiry"`
}

type Condition func(item AppMsgListItem) (bool, error)
