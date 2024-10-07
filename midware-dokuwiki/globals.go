/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2024-09-24 22:49:18
 * @LastEditTime: 2024-10-07 17:00:34
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
	LINUX_DEV_SHM       string  = "/dev/shm"
	OSM_ZOOM            string  = "18"
	GEN_BBOX_LAT_OFFSET float64 = 0.0005
	GEN_BBOX_LON_OFFSET float64 = 0.005
)

// DEFAULTS //
const (
	DEFAULT_DB                string = "akbp.db"
	DEFAULT_TMP_DIR           string = "akbp-tmp"
	DEFAULT_TIME_TEMPLATE     string = "2006-01-02 15:04:05"
	DEFAULT_NON_INI_TEMPLATE  string = "DefaultNonIni.template"
	DEFAULT_INI_TYPE_TEMPLATE string = "DefaultIniType.template"
	DEFAULT_ALL_ENTRIE_NS     string = "akbp:all"
	DEFAULT_EVENTS_NS         string = "akbp:events"
	DEFAULT_KMLS_NS           string = "akbp:kmls"
)

// COMMON VARS //

var (
	ConfigFile string = "" // Path of the config file.
)

// CONFIG //

var (
	ConfLoaded      bool   = false                     // To prevent data races, do not change it manually.
	TmpDir          string = ""                        // Tmp dir.
	Db              string = DEFAULT_DB                // DB to use.
	UpdGap          int    = 10                        // Update gap.
	TimeTemplate    string = DEFAULT_TIME_TEMPLATE     // Template for time.
	NonIniTemplate  string = DEFAULT_NON_INI_TEMPLATE  // NonIniTemplate.
	IniTypeTemplate string = DEFAULT_INI_TYPE_TEMPLATE // IniTypeTemplate.
	AllEntriesNs    string = DEFAULT_ALL_ENTRIE_NS     // Ns for all entries.
	EventsNs        string = DEFAULT_EVENTS_NS         // Ns for "by events".
	KmlsNs          string = DEFAULT_KMLS_NS           // Ns for kmls.
	Domain          string = ""                        // Domain.
	UsersDir        string = ""                        // Users dirs.
)
