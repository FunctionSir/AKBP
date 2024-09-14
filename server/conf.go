/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2024-09-14 00:23:17
 * @LastEditTime: 2024-09-14 19:28:08
 * @LastEditors: FunctionSir
 * @Description: Config related.
 * @FilePath: /AKBP/server/conf.go
 */

package main

import (
	"strings"

	"gopkg.in/ini.v1"
)

func LoadConf() {
	if ConfLoaded {
		LogFatalln("Error: reload config might cause data races.")
	}
	conf, err := ini.Load(ConfigFile)
	if err != nil {
		LogFatalln("Error occurred when loading conf file \"" + ConfigFile + "\": " + strings.Trim(err.Error(), "\n") + ".")
	}
	if !conf.HasSection("server") {
		LogWarnln("Warning: config file \"" + ConfigFile + "\" has no section \"server\", default values about server will be applied for everything.")
		return
	}
	section := conf.Section("server")
	if section.HasKey("Addr") {
		Addr = section.Key("Addr").String()
	}
}
