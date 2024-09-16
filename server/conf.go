/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2024-09-14 00:23:17
 * @LastEditTime: 2024-09-16 21:18:51
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
	if section.HasKey("ServerID") {
		ServerId = section.Key("ServerID").String()
	}
	if conf.HasSection("known-servers") {
		KnownServers = make(map[string]string)
		KeysForAuth = make(map[string]string)
		sec_known_servers := conf.Section("known-servers")
		for _, key := range sec_known_servers.KeyStrings() {
			tmp := strings.Split(sec_known_servers.Key(key).String(), "$")
			if len(tmp) == 2 {
				KnownServers[key] = strings.TrimSpace(tmp[0])
				KeysForAuth[key] = strings.TrimSpace(tmp[1])
			} else {
				LogFatalln("Wrong value \"" + sec_known_servers.Key(key).String() + "\" in section known-servers of file \"" + ConfigFile + "\".")
			}
		}
	}
}
