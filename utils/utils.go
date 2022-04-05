package utils

import (
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/tebeka/selenium"
)

func ConvertToHttpCookies(sourceCookies []selenium.Cookie) []http.Cookie {
	var cookies []http.Cookie

	for idx := 0; idx < len(sourceCookies); idx++ {
		cookie := http.Cookie{
			Name:    sourceCookies[idx].Name,
			Value:   sourceCookies[idx].Value,
			Path:    sourceCookies[idx].Path,
			Domain:  sourceCookies[idx].Domain,
			Secure:  sourceCookies[idx].Secure,
			Expires: time.Unix(int64(sourceCookies[idx].Expiry), 0),
		}
		cookies = append(cookies, cookie)
	}

	return cookies
}

func ConvertToSeleniumCookie(sourceCookies *http.Cookie) selenium.Cookie {
	var cookie selenium.Cookie

	cookie = selenium.Cookie{
		Name:   sourceCookies.Name,
		Value:  sourceCookies.Value,
		Path:   sourceCookies.Path,
		Domain: sourceCookies.Domain,
		Secure: sourceCookies.Secure,
		Expiry: uint(sourceCookies.Expires.Unix()),
	}

	return cookie
}

func IdxofCookieswithName(cookies []selenium.Cookie, name string) int {
	for idx := 0; idx < len(cookies); idx++ {
		if cookies[idx].Name == name {
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
