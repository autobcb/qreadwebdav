package config

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"strings"
)

var Port = 8080
var Url = ""
var Users = []string{}

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	p, err := cfg.Section("web").Key("port").Int()
	if err != nil {
		Port = p
	}
	Url = cfg.Section("web").Key("url").String()
	users := cfg.Section("web").Key("users").String()
	if Url == "" || !strings.HasPrefix(Url, "http") {
		println("url must start with http:// or https://")
		os.Exit(0)
	}
	println("Port:", Port)
	println("Url:", Url)
	println("Users:", users)
	s := strings.Split(users, ",")
	for _, u := range s {
		if u != "" {
			Users = append(Users, u)
		}
	}
}
