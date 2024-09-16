/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2024-09-12 22:19:36
 * @LastEditTime: 2024-09-16 22:42:09
 * @LastEditors: FunctionSir
 * @Description: Global consts and vars.
 * @FilePath: /AKBP/server/globals.go
 */

package main

// ABOUT //

const (
	VER      string = "0.0.1"           // Version.
	CODENAME string = "NunotabaShinobu" // Codename of this version.
)

// DEFAULTS //

const (
	DEFAULT_ADDR string = "127.0.0.1:4060" // Default listening address.
	DEFAULT_DB   string = "akbp.db"        // Default DB.
	DEFAULT_SID  string = ""               // Default server ID.
)

// COMMON CONSTS //

const (
	ADD_TTL int = 64
)

// COMMON VARS //

var (
	DebugMode  bool   = false       // Switch to true to enable debug for Gin.
	ConfigFile string = ""          // Path of config file.
	ServerId   string = DEFAULT_SID // Server ID.
)

// CONFIG //

var (
	ConfLoaded   bool              = false        // To prevent data races, do not change it manually.
	Addr         string            = DEFAULT_ADDR // Listening address.
	Db           string            = DEFAULT_DB   // DB to use.
	KnownServers map[string]string = nil
	KeysForAuth  map[string]string = nil
)

// FOR SQL //

const (
	TABLE_BEACONS string = "BEACONS"
	TABLE_SERVERS string = "SERVERS"
	STR_BEACON    string = "BEACON"
	STR_SERVER    string = "SERVER"
)

// HEADERS //

const (
	KEY_AKBP_MSG_TYPE          string = "X-Akbp-Msg-Type"
	KEY_AKBP_BEACON_ID         string = "X-Akbp-Beacon-Id"
	KEY_AKBP_AUTH              string = "X-Akbp-Auth"
	KEY_AKBP_TIMESTAMP         string = "X-Akbp-Timestamp"
	KEY_AKBP_EVENT_ID          string = "X-Akbp-Event-Id"
	KEY_AKBP_DO_NOT_FORWARD_TO string = "X-Akbp-Do-Not-Forward-To"
	KEY_AKBP_SERVER_ID         string = "X-Akbp-Server-Id"
	KEY_AKBP_TTL               string = "X-Akbp-Ttl"
	KEY_CONTENT_TYPE           string = "Content-Type"
)

// ERRORS //

const (
	ERR_UNKNOWN_MSG_TYPE             string = "400 BadRequest UnknownMsgType"
	ERR_BAD_TIMESTAMP                string = "400 BadRequest BadTimestamp"
	ERR_NO_VALID_EID                 string = "400 BadRequest NoValidEventId"
	ERR_NO_VALID_AUTH_HEADER         string = "401 Unauthorized NoValidAuthHeader"
	ERR_WRONG_BEACON_ID_OR_AUTH_INFO string = "403 Forbidden WrongBeaconIdOrAuthInfo"
)

// MAKE A CUP OF COFFEE //

const (
	ALWAYS_TEAPOT bool   = false
	FORTUNE       string = "/usr/bin/fortune"
	I_AM_A_TEAPOT string = "418 I'm a teapot"
)
