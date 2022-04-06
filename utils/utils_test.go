package utils

import (
	"github.com/tebeka/selenium"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestConvertToHttpCookie(t *testing.T) {
	seleniumCookie := selenium.Cookie{
		Name:   "cookie",
		Value:  "cookie",
		Path:   "/",
		Domain: "example.com",
		Secure: true,
		Expiry: 1649400579,
	}

	httpCookie := http.Cookie{
		Name:    "cookie",
		Value:   "cookie",
		Path:    "/",
		Domain:  "example.com",
		Secure:  true,
		Expires: time.Unix(1649400579, 0),
	}

	if ret := ConvertToHttpCookie(seleniumCookie); !reflect.DeepEqual(ret, httpCookie) {
		t.Errorf("convert %v to %v, but %v got.\n", seleniumCookie, httpCookie, ret)
	}
}

func TestConvertToSeleniumCookie(t *testing.T) {
	seleniumCookie := selenium.Cookie{
		Name:   "cookie",
		Value:  "cookie",
		Path:   "/",
		Domain: "example.com",
		Secure: true,
		Expiry: 1649400579,
	}

	httpCookie := http.Cookie{
		Name:    "cookie",
		Value:   "cookie",
		Path:    "/",
		Domain:  "example.com",
		Secure:  true,
		Expires: time.Unix(1649400579, 0),
	}

	if ret := ConvertToSeleniumCookie(&httpCookie); !reflect.DeepEqual(ret, seleniumCookie) {
		t.Errorf("convert %v to %v, but %v got.\n", httpCookie, seleniumCookie, ret)
	}
}

func TestJsonMarshalwithNoHTMLEscape(t *testing.T) {
	source := AppMsgListItems{
		[]AppMsgListItem{
			AppMsgListItem{
				Aid:        "2651958334_1",
				AlbumId:    "0",
				AppmsgId:   0,
				Checking:   0,
				Cover:      "https://example.com",
				CreateTime: 0,
				Digest:     "",
				ItemIdx:    0,
				Link:       "https://example.com/index.php?a=1&b=2",
				Title:      "<em>",
				UpdateTime: 0,
			},
		},
	}

	target := "{\n  \"app_msg_list\": [\n    " +
		"{\n      " +
		"\"aid\": \"2651958334_1\",\n      " +
		"\"album_id\": \"0\",\n      " +
		"\"appmsgid\": 0,\n      " +
		"\"checking\": 0,\n      " +
		"\"cover\": \"https://example.com\",\n      " +
		"\"create_time\": 0,\n      " +
		"\"digest\": \"\",\n      " +
		"\"itemidx\": 0,\n      " +
		"\"link\": \"https://example.com/index.php?a=1&b=2\",\n      " +
		"\"title\": \"<em>\",\n      " +
		"\"update_time\": 0\n    " +
		"}\n  ]\n}\n"

	ret, err := JsonMarshalwithNoHTMLEscape(source)
	if err != nil || string(ret) != target {
		t.Errorf("format %v to %v json, buf got %v\n", source, target, ret)
	}

}
