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
