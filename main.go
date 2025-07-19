package main

import (
	"golang.org/x/net/webdav"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"webdav/config"
	_ "webdav/config"
	"webdav/qread"
)

const path = "webdav"

func main() {
	qread.Checkpath(path)
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// 获取用户名/密码
		username, password, ok := req.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// 验证用户名/密码
		if !Check(username, password) {
			http.Error(w, "WebDAV: need authorized!", http.StatusUnauthorized)
			return
		}
		cacheFile := filepath.Join(path, qread.Md5V(username))
		qread.Checkpath(cacheFile)
		fs := &webdav.Handler{
			FileSystem: webdav.Dir(cacheFile),
			LockSystem: webdav.NewMemLS(),
		}
		var method = strings.ToLower(req.Method)
		//处理读取进度和写入进度
		if strings.HasSuffix(req.URL.Path, ".json") && strings.Contains(req.URL.Path, "bookProgress") && strings.ToLower(method) == "get" {
			qread.ChcekProgress(method, filepath.Join(cacheFile, req.URL.Path), username, password)
		}
		fs.ServeHTTP(w, req)
		if strings.HasSuffix(req.URL.Path, ".json") && strings.Contains(req.URL.Path, "bookProgress") && strings.ToLower(method) == "put" {
			qread.ChcekProgress(method, filepath.Join(cacheFile, req.URL.Path), username, password)
		}
	})
	err := http.ListenAndServe(":"+strconv.Itoa(config.Port), nil)
	if err != nil {
		println(err.Error())
	}
}

func Check(username string, password string) bool {
	if len(config.Users) != 0 {
		if !inArray(config.Users, username) {
			return false
		}
	}
	return qread.Login(username, password)
}

func inArray(arr []string, target string) bool {
	for _, v := range arr {
		if v == target {
			return true
		}
	}
	return false
}
