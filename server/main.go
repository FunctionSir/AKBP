/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2024-09-12 22:07:34
 * @LastEditTime: 2024-09-14 23:19:03
 * @LastEditors: FunctionSir
 * @Description: AKBP Server, main file.
 * @FilePath: /AKBP/server/main.go
 */

package main

import (
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

// You can find global consts and vars in globals.go //

func initialize() {
	for _, x := range os.Args {
		switch x {
		case "--no-color":
			color.NoColor = true
			gin.DisableConsoleColor()
		case "--debug":
			gin.SetMode(gin.DebugMode)
			DebugMode = true
		}
	}
	if !color.NoColor {
		gin.ForceConsoleColor()
	}
	if !DebugMode {
		gin.SetMode(gin.ReleaseMode)
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

func beaconReportHandler(c *gin.Context) {
	msgType := c.GetHeader(KEY_AKBP_MSG_TYPE)
	switch msgType {
	case "ping":
		c.Header(KEY_AKBP_MSG_TYPE, "pong")
		c.String(http.StatusOK, "pong")
	case "beacon-report":
		id := c.GetHeader(KEY_AKBP_BEACON_ID)
		if len(id) == 0 {
			c.Header(KEY_AKBP_MSG_TYPE, "error")
			c.String(http.StatusBadRequest, ERR_WRONG_BEACON_ID_OR_AUTH_INFO)
			return
		}
		key := c.GetHeader(KEY_AKBP_AUTH)
		if len(key) == 0 {
			c.Header(KEY_AKBP_MSG_TYPE, "error")
			c.String(http.StatusUnauthorized, ERR_NO_VALID_AUTH_HEADER)
			return
		}
		if !BeaconAuthOk(id, key) {
			c.Header(KEY_AKBP_MSG_TYPE, "error")
			c.String(http.StatusForbidden, ERR_WRONG_BEACON_ID_OR_AUTH_INFO)
			return
		}
		eid := c.GetHeader(KEY_AKBP_EVENT_ID)
		if len(eid) == 0 {
			c.Header(KEY_AKBP_MSG_TYPE, "error")
			c.String(http.StatusBadRequest, ERR_NO_VALID_EID)
			return
		}
		ts, err := strconv.Atoi(c.GetHeader(KEY_AKBP_TIMESTAMP))
		if err != nil {
			c.Header(KEY_AKBP_MSG_TYPE, "error")
			c.String(http.StatusBadRequest, ERR_BAD_TIMESTAMP)
			return
		}
		AddRecord(id, eid, ts, c.PostForm("msg"))
	default:
		c.Header(KEY_AKBP_MSG_TYPE, "error")
		c.String(http.StatusBadRequest, ERR_UNKNOWN_MSG_TYPE)
	}
}

func main() {
	initialize()
	LoadConf()
	ginEng := gin.Default()

	// FAVICON //
	ginEng.GET("/favicon.ico", func(c *gin.Context) {
		c.Header(KEY_AKBP_MSG_TYPE, "favicon")
		c.String(http.StatusOK, "")
	})

	// Ping //
	ginEng.GET("/ping", func(c *gin.Context) {
		c.Header(KEY_AKBP_MSG_TYPE, "pong")
		c.String(http.StatusOK, "pong")
	})

	// Fortune //
	if runtime.GOOS == "linux" && FileExists(FORTUNE) && !ALWAYS_TEAPOT {
		ginEng.GET("/fortune", func(c *gin.Context) {
			c.Header(KEY_AKBP_MSG_TYPE, "coffee")
			output, _ := exec.Command(FORTUNE).Output()
			c.String(http.StatusOK, strings.Trim(string(output), "\n"))
		})
	} else {
		ginEng.GET("/fortune", func(c *gin.Context) {
			c.Header(KEY_AKBP_MSG_TYPE, "teapot")
			c.String(418, I_AM_A_TEAPOT)
		})
	}

	// Beacon-Report //
	ginEng.POST("/beacon-report", beaconReportHandler)

	ginEng.Run(Addr)
}
