package qread

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/nahid/gohttp"
	"io/ioutil"
	"os"
	"path/filepath"
	"webdav/config"
	"webdav/utils/mjson"
)

func Login(username, password string) bool {
	if readaccessToken(username) != "" {
		return true
	}
	req := gohttp.NewRequest()

	resp, err := req.
		FormData(map[string]string{"username": username, "password": password, "model": "gowebdav"}).
		Post(config.Url + "/api/5/login")
	if err != nil {
		return false
	}
	if resp.GetStatusCode() == 200 {
		body, err := resp.GetBodyAsString()
		if err != nil {
			return false
		}
		data, err := mjson.ParseHasErr(body)
		if err != nil {
			return false
		}
		if data["isSuccess"] == true {
			accessToken := data["data"].(map[string]interface{})["accessToken"].(string)
			writeaccessToken(username, accessToken)
			return true
		}
	}

	return false
}
func Md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func writeaccessToken(user, accessToken string) error {
	cacheFileName := Md5V(user) + ".json"
	cacheFile := filepath.Join(cachepath, cacheFileName)
	file, err := os.Create(cacheFile)
	if err != nil {
		return fmt.Errorf("创建缓存文件失败: %w", err)
	}
	defer file.Close()
	file.WriteString(accessToken)
	return nil
}

func delaccessToken(user string) {
	cacheFileName := Md5V(user) + ".json"

	cacheFile := filepath.Join(cachepath, cacheFileName)
	file, err := os.Open(cacheFile)
	if err == nil {
		defer file.Close()
		os.Remove(cacheFile)
	}
}

func readaccessToken(user string) string {
	cacheFileName := Md5V(user) + ".json"

	cacheFile := filepath.Join(cachepath, cacheFileName)
	file, err := os.Open(cacheFile)
	if err != nil {
		return ""
	}
	defer file.Close()
	bytes, _ := ioutil.ReadAll(file)
	return string(bytes)
}
