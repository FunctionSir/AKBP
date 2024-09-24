/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2024-09-22 14:46:17
 * @LastEditTime: 2024-09-24 23:44:13
 * @LastEditors: FunctionSir
 * @Description: The bridge to dokuwiki.
 * @FilePath: /AKBP/midware-dokuwiki/main.go
 */

package main

import (
	"os"
	"runtime"

	"github.com/fatih/color"
)

// You can find global consts and vars in globals.go //

func initialize() {
	for _, x := range os.Args {
		switch x {
		case "--no-color":
			color.NoColor = true
		}
	}
	Hello() // Say hello.
	if len(os.Args) <= 1 {
		LogFatalln("No config file specified.")
	}
	ConfigFile = os.Args[1]
	if !FileExists(ConfigFile) {
		LogFatalln("Config file specified not exists or is a dir.")
	}
}

func main() {
	initialize()
	if runtime.GOOS == "linux" {
		TmpDir = LINUX_DEV_SHM + "/" + DEFAULT_TMP_DIR
	} else {
		TmpDir = DEFAULT_TMP_DIR
	}
	LoadConf()
	if !FileExists(Db) {
		LogFatalln("Database \"" + Db + "\" not exists. Please check your config.")
	}
	if len(TmpDir) < len(LINUX_DEV_SHM) || TmpDir[0:len(LINUX_DEV_SHM)] != LINUX_DEV_SHM {
		LogWarnln("A large amount of write might will be preformed. Please make sure you are using a ram disk, or a HDD, or a SSD but without any valuable data in it for tmp dir.")
	}

}
