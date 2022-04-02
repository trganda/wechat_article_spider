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
