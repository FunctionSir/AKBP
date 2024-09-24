/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2024-09-14 00:23:17
 * @LastEditTime: 2024-09-24 23:40:34
 * @LastEditors: FunctionSir
 * @Description: Config related.
 * @FilePath: /AKBP/midware-dokuwiki/conf.go
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
	if !conf.HasSection("midware-dokuwiki") {
		LogWarnln("Warning: config file \"" + ConfigFile + "\" has no section \"midware-dokuwiki\", default values about midware-dokuwiki will be applied for everything.")
		return
	}
	section := conf.Section("midware-dokuwiki")
	if section.HasKey("MainDB") {
		Db = section.Key("MainDB").String()
	} else {
		LogWarnln("Seems like you didn't specify the DB you want to use, the default value \"" + DEFAULT_DB + "\" will be applied.")
	}
	if section.HasKey("TmpDir") {
		TmpDir = section.Key("TmpDir").String()
	} else {
		LogWarnln("Seems like you didn't specify the tmp dir you want to use, the gened value \"" + TmpDir + "\" will be applied.")
	}
}
