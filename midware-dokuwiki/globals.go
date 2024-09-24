/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2024-09-24 22:49:18
 * @LastEditTime: 2024-09-24 23:32:25
 * @LastEditors: FunctionSir
 * @Description: -
 * @FilePath: /AKBP/midware-dokuwiki/globals.go
 */

package main

// ABOUT //

const (
	VER      string = "0.0.1"           // Version.
	CODENAME string = "NunotabaShinobu" // Codename of this version.
)

// COMMON CONSTS //
const (
	LINUX_DEV_SHM string = "/dev/shm"
)

// DEFAULTS //
const (
	DEFAULT_DB      string = "akbp.db"
	DEFAULT_TMP_DIR string = "akbp-tmp"
)

// COMMON VARS //

var (
	ConfigFile string = "" // Path of config file.
)

// CONFIG //

var (
	ConfLoaded bool   = false      // To prevent data races, do not change it manually.
	TmpDir     string = ""         // Tmp dir.
	Db         string = DEFAULT_DB // DB to use.
)
