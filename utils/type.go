package utils

type AppMsgArgs struct {
	Token  string `json:"token"` // csrf token
	Lang   string `json:"lang"`  // language
	F      string `json:"f"`     // format
	Ajax   string `json:"ajax"`  // request type
	Action string `jsong:"list_ex"`
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
