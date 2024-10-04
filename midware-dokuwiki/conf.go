/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2024-09-14 00:23:17
 * @LastEditTime: 2024-10-04 22:38:01
 * @LastEditors: FunctionSir
 * @Description: Config related.
 * @FilePath: /AKBP/midware-dokuwiki/conf.go
 */

package main

import (
	"strconv"
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
	if section.HasKey("UpdateGap") {
		UpdGap, err = strconv.Atoi(section.Key("UpdateGap").String())
		if err != nil || UpdGap <= 0 {
			LogFatalln("Can not convert value of key \"UpdateGap\" from string to a positive int. Is this key has a wrong value?")
		}
	} else {
		LogWarnln("Seems like you didn't specify the update gap, use default value 10.")
	}
	if section.HasKey("NonIniTemplate") {
		NonIniTemplate = section.Key("NonIniTemplate").String()
	} else {
		LogWarnln("Seems like you didn't specify template for non-ini messages, use the default one.")
	}
	if section.HasKey("IniTypeTemplate") {
		IniTypeTemplate = section.Key("IniTypeTemplate").String()
	} else {
		LogWarnln("Seems like you didn't specify template for ini messages, use the default one.")
	}
	if section.HasKey("NsForAll") {
		AllEntriesNs = section.Key("NsForAll").String()
	} else {
		LogWarnln("Seems like you didn't specify the namespace linked to tmpdir/dw-ns/all, use the default one.")

	}
	if section.HasKey("NsForEvents") {
		EventsNs = section.Key("NsForEvents").String()
	} else {
		LogWarnln("Seems like you didn't specify the namespace linked to tmpdir/dw-ns/evnets, use the default one.")

	}
	if section.HasKey("NsForKmls") {
		KmlsNs = section.Key("NsForKmls").String()
	} else {
		LogWarnln("Seems like you didn't specify the namespace linked to tmpdir/dw-m/kmls, use the default one.")
	}
}
