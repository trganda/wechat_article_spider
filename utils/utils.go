package utils

import (
	"net/http"
	"os"
	"time"

	"github.com/tebeka/selenium"
)

func ConvertToHttpCookie(sourceCookies []selenium.Cookie) []http.Cookie {
	cookies := []http.Cookie{}

	for idx := 0; idx < len(sourceCookies); idx++ {
		cookie := http.Cookie{
			Name:    sourceCookies[idx].Name,
			Value:   sourceCookies[idx].Value,
			Path:    sourceCookies[idx].Path,
			Domain:  sourceCookies[idx].Domain,
			Expires: time.Unix(int64(sourceCookies[idx].Expiry), 0),
		}
		cookies = append(cookies, cookie)
	}

	return cookies
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}
