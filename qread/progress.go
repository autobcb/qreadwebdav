package qread

import (
	"errors"
	"fmt"
	"github.com/nahid/gohttp"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"webdav/config"
	"webdav/utils/mjson"
	"webdav/utils/strkit"
)

func ChcekProgress(method, path, username, password string) {
	var accessToken = readaccessToken(username)
	if accessToken == "" {
		if !Login(username, password) {
			return
		}
		accessToken = readaccessToken(username)
	}
	if method == "get" {
		txt := readtxt(path)
		if txt != "" {
			data, err := mjson.ParseHasErr(txt)
			if err == nil {
				var author = strkit.ToString(data["author"])
				var name = strkit.ToString(data["name"])
				var durChapterTitle = strkit.ToString(data["durchapter_title"])
				var durChapterIndex int64 = 0
				var durChapterTime int64 = 0
				var durChapterPos int64 = 0
				d1, err := strconv.ParseInt(strkit.ToString(data["durChapterIndex"]), 10, 64)
				if err == nil {
					durChapterIndex = d1
				}
				d2, err := strconv.ParseInt(strkit.ToString(data["durChapterTime"]), 10, 64)
				if err == nil {
					durChapterTime = d2
				}
				d3, err := strconv.ParseInt(strkit.ToString(data["durChapterPos"]), 10, 64)
				if err == nil {
					durChapterPos = d3
				}

				books, err := Getbooks(accessToken, name, 0)
				if err != nil && err.Error() == "NEED_LOGIN" {
					delaccessToken(username)
					ChcekProgress(method, path, username, password)
					return
				}
				if err != nil {
					println(err.Error())
				}
				if err == nil && len(books) > 0 {
					println("开始寻找书本")
					for _, l := range books {
						book, err := mjson.ParseHasErr(mjson.ToJson(l))
						if err != nil {
							continue
						}
						//相同书如果使用阅读时间最前的进度
						if strkit.ToString(book["name"]) == name && strkit.ToString(book["author"]) == author {
							println("找到书本:" + name + " <UNK>:" + author)
							durChapterTime1, err := strconv.ParseInt(strkit.ToString(book["durChapterTime"]), 10, 64)
							if err == nil && durChapterTime1 > durChapterTime {
								durChapterIndex1, err := strconv.ParseInt(strkit.ToString(book["durChapterIndex"]), 10, 64)
								if err == nil && durChapterIndex1 != durChapterIndex {
									println("找到书本:" + name + "轻阅读阅读时间大于阅读")
									durChapterTime = durChapterTime1
									durChapterIndex = durChapterIndex1
									durChapterPos = 0
									durChapterTitle = strkit.ToString(book["durChapterTitle"])
								}
							}
							break
						}
					}
				}

				data["durChapterIndex"] = durChapterIndex
				data["durChapterTime"] = durChapterTime
				data["durChapterPos"] = durChapterPos
				data["durChapterTitle"] = durChapterTitle
				os.Remove(path)
				writetxt(path, mjson.ToJson(data))
			}
		}
	} else if method == "put" {
		txt := readtxt(path)
		if txt != "" {
			data, err := mjson.ParseHasErr(txt)
			if err == nil {
				var author = strkit.ToString(data["author"])
				var name = strkit.ToString(data["name"])
				var durChapterTitle = strkit.ToString(data["durchapter_title"])
				var durChapterIndex int64 = 0
				//var durChapterTime int64 = 0
				d1, err := strconv.ParseInt(strkit.ToString(data["durChapterIndex"]), 10, 64)
				if err == nil {
					durChapterIndex = d1
				}

				books, err := Getbooks(accessToken, name, 0)
				if err != nil && err.Error() == "NEED_LOGIN" {
					delaccessToken(username)
					ChcekProgress(method, path, username, password)
					return
				}
				if err != nil {
					println(err.Error())
				}
				if err == nil && len(books) > 0 {
					println("开始寻找书本")
					for _, l := range books {
						book, err := mjson.ParseHasErr(mjson.ToJson(l))
						if err != nil {
							continue
						}
						//相同书如果使用阅读时间最前的进度
						if strkit.ToString(book["name"]) == name && strkit.ToString(book["author"]) == author {
							println("找到书本:" + name + " <UNK>:" + author + "，即将上传进度")
							var bookUrl = strkit.ToString(book["bookUrl"])
							saveBookProgress(accessToken, bookUrl, durChapterTitle, strkit.ToString(durChapterIndex), 0)
							break
						}
					}
				}
			}
		}
	}
}

func readtxt(path string) string {
	file, err := os.Open(path)
	if err == nil {
		defer file.Close()
		bytes, _ := ioutil.ReadAll(file)
		return string(bytes)
	}
	return ""
}

func writetxt(path, txt string) {
	file, err := os.Create(path)
	if err != nil {
		println(err.Error())
	} else {
		defer file.Close()
		_, err := file.WriteString(txt)
		if err != nil {
			println(err.Error())
		}
	}
}

func Getbooks(accessToken string, name string, time int) (list []interface{}, err error) {
	defer func() {
		if nerr := recover(); nerr != nil {
			err = fmt.Errorf("%v", nerr)
		}
	}()
	var url = config.Url + "/api/5/getBookshelf?accessToken=" + accessToken + "&name=" + url.QueryEscape(name)
	req := gohttp.NewRequest()

	resp, err := req.
		Get(url)
	if err != nil {
		if time < 5 {
			return Getbooks(accessToken, name, time+1)
		}
	}
	if resp.GetStatusCode() == 200 {
		body, err := resp.GetBodyAsString()
		if err != nil {
			if time > 5 {
				return []interface{}{}, err
			}
			return Getbooks(accessToken, name, time+1)
		}
		data, err := mjson.ParseHasErr(body)
		if err != nil {
			if time > 5 {
				return []interface{}{}, err
			}
			return Getbooks(accessToken, name, time+1)
		}
		if strings.Contains(strkit.ToString(data["errorMsg"]), "NEED_LOGIN") {
			return []interface{}{}, errors.New("NEED_LOGIN")
		}
		if data["isSuccess"] == true {
			return data["data"].([]interface{}), nil
		}
	}
	if time > 5 {
		return []interface{}{}, errors.New("请求失败")
	}
	return Getbooks(accessToken, name, time+1)
}

func saveBookProgress(accessToken, url, title, index string, time int) {
	var url1 = config.Url + "/api/5/saveBookProgress?accessToken=" + accessToken
	req := gohttp.NewRequest()

	resp, err := req.FormData(map[string]string{
		"url":   url,
		"title": title,
		"index": index,
		"pos":   "0",
	}).Post(url1)
	if err != nil {
		if time < 5 {
			saveBookProgress(accessToken, url, title, index, time+1)
			return
		}
	}
	if resp.GetStatusCode() == 200 {
		body, err := resp.GetBodyAsString()
		if err != nil {
			if time < 5 {
				saveBookProgress(accessToken, url, title, index, time+1)
			}
			return
		}
		data, err := mjson.ParseHasErr(body)
		if err != nil {
			if time < 5 {
				saveBookProgress(accessToken, url, title, index, time+1)
			}
			return
		}
		if data["isSuccess"] == true {
			return
		}
	}
	if time < 5 {
		saveBookProgress(accessToken, url, title, index, time+1)
	}
}
