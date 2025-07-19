package qread

import "os"

const cachepath = "cache"

func init() {
	Checkpath(cachepath)
}

func Checkpath(path string) {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if !os.IsExist(err) {
			os.Mkdir(path, os.ModePerm)
		}
	}
	//return true
}
