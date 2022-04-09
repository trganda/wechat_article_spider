package crawer

import (
	"io/ioutil"
	"testing"
)

func TestLogining(t *testing.T) {
	Logining()
}

func TestPrintQRCode(t *testing.T) {
	file, err := ioutil.ReadFile("qrcode.jpg")
	if err != nil {
		return
	}
	PrintQRCode(file)
}
