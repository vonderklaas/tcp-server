package helpers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"
)

func StrToInt(str string) (int, error) {
	return strconv.Atoi(str)
}

func IntToStr(i int) string {
	return strconv.Itoa(i)
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func HttpGetWithTimeout(url string, timeout time.Duration) (*http.Response, error) {
	client := http.Client{
		Timeout: timeout,
	}
	return client.Get(url)
}

func MarshalJson(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
