package crawler

import (
	"testing"
	"wechat_crawler/utils"
)

func TestDumpPage(t *testing.T) {
	item := utils.AppMsgListItem{
		Aid:        "2651958334_1",
		AlbumId:    "0",
		AppmsgId:   0,
		Checking:   0,
		Cover:      "https://example.com",
		CreateTime: 0,
		Digest:     "",
		ItemIdx:    0,
		Link: "http://mp.weixin.qq.com/s?__biz=MzA5NDYyNDI0MA==&" +
			"mid=2651958334&idx=1&" +
			"sn=d76bb032833c91480eb3e376cc49ca06&" +
			"chksm=8baecca1bcd945b71b1f5a8294200e99d9885c66afc03c8280394404775314bef825b456f2ee#rd",
		Title:      "每日安全动态推送(04-06)",
		UpdateTime: 0,
	}

	DumpItem(item, "")
}
