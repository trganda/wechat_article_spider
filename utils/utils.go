package utils

import (
	"math/rand"
	"net/http"
	"os"
	"strconv"
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
