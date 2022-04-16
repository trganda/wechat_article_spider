package utils

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

func ConvertToHttpCookies(sourceCookies Cookies) []http.Cookie {
	var cookies []http.Cookie

	for idx := 0; idx < len(sourceCookies.Cookies); idx++ {
		cookie := ConvertToHttpCookie(sourceCookies.Cookies[idx])
		cookies = append(cookies, cookie)
	}

	return cookies
}

func ConvertToHttpCookie(sourceCookie *Cookie) http.Cookie {
	var cookie http.Cookie

	cookie = http.Cookie{
		Name:    sourceCookie.Name,
		Value:   sourceCookie.Value,
		Path:    sourceCookie.Path,
		Domain:  sourceCookie.Domain,
		Secure:  sourceCookie.Secure,
		Expires: time.Unix(int64(sourceCookie.Expiry), 0),
	}

	return cookie
}

func ConvertToSeleniumCookie(sourceCookies *http.Cookie) Cookie {
	var cookie Cookie

	cookie = Cookie{
		Name:   sourceCookies.Name,
		Value:  sourceCookies.Value,
		Path:   sourceCookies.Path,
		Domain: sourceCookies.Domain,
		Secure: sourceCookies.Secure,
		Expiry: uint(sourceCookies.Expires.Unix()),
	}

	return cookie
}

func IdxofCookieswithName(cookies Cookies, name string) int {
	for idx := 0; idx < len(cookies.Cookies); idx++ {
		if cookies.Cookies[idx].Name == name {
			return idx
		}
	}
	return -1
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

// RandDuration create a time.Duration between [5 10] second.
func RandDuration() time.Duration {
	rand.Seed(time.Now().Unix())
	randNumber := rand.Intn(5)

	randNumber += 5

	duration, err := time.ParseDuration(strconv.Itoa(randNumber) + "s")
	if err != nil {
		return 0
	}

	return duration
}

func JsonMarshalwithNoHTMLEscape(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})

	jsonEncoder := json.NewEncoder(buf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.SetIndent("", "  ")
	err := jsonEncoder.Encode(data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
