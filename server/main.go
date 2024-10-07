/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2024-09-12 22:07:34
 * @LastEditTime: 2024-10-07 18:07:40
 * @LastEditors: FunctionSir
 * @Description: AKBP Server, main file.
 * @FilePath: /AKBP/server/main.go
 */

package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"

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
			DebugMode = true
		}
	}
	if !color.NoColor {
		gin.ForceConsoleColor()
	}
	if !DebugMode {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
		LogWarnln("You are in debug mode currently.")
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

func doPost(url string, headers *http.Header, kv *map[string]string) {
	tmpList := []string{}
	for k, v := range *kv {
		tmpList = append(tmpList, k+"="+v)
	}
	bodyStr := strings.Join(tmpList, "&")
	req, _ := http.NewRequest("POST", url, strings.NewReader(bodyStr))
	req.Header.Set(KEY_CONTENT_TYPE, "application/x-www-form-urlencoded")
	for k, v := range *headers {
		req.Header.Add(k, strings.Join(v, ","))
	}
	// Just let it alone.
	go http.DefaultClient.Do(req)
}

func reportOk(c *gin.Context) {
	c.Header(KEY_AKBP_MSG_TYPE, "server-report")
	c.String(http.StatusOK, HTTP_200_OK+" "+fmt.Sprint(time.Now().UnixMilli()))
}

func pong(c *gin.Context) {
	c.Header(KEY_AKBP_MSG_TYPE, "pong")
	c.String(http.StatusOK, "pong")
}

func srvToSrvExchange() {
	for {
		LogInfoln("Started server to server records exchange.")
		// If no valid server ID or no known server, do not exchange.
		if ServerId == "" || KnownServers == nil {
			return
		}
		rows := QueryRecs()
		for rows.Next() {
			var bid string
			var eid string
			var ts int
			var msg string
			var banned string
			rows.Scan(&bid, &eid, &ts, &msg, &banned)
			notForwardTo := strings.Split(banned, ",")
			// Make a new header.
			newHeader := make(http.Header)
			newHeader.Set(KEY_AKBP_TTL, strconv.Itoa(ADD_TTL))
			// Set type and server ID info.
			newHeader.Set(KEY_AKBP_MSG_TYPE, "forwarded")
			newHeader.Set(KEY_AKBP_SERVER_ID, ServerId)
			newHeader.Set(KEY_AKBP_DO_NOT_FORWARD_TO, banned)
			// BID, EID, TS
			newHeader.Set(KEY_AKBP_BEACON_ID, bid)
			newHeader.Set(KEY_AKBP_EVENT_ID, eid)
			newHeader.Set(KEY_AKBP_TIMESTAMP, strconv.Itoa(ts))
			// Encoded msg here.
			encodedMsg := url.QueryEscape(msg)
			newKv := map[string]string{"msg": encodedMsg}
			// Range known servers.
			for k, v := range KnownServers {
				// If the server is contained in notForwardTo.
				if slices.Contains(notForwardTo, k) {
					continue
				}
				// Set new auth info.
				newHeader.Set(KEY_AKBP_AUTH, KeysForAuth[k])
				// Post it!
				doPost(v, &newHeader, &newKv)
			}
		}
		LogInfoln("Successfully done once server to server records exchange.")
		// Sit back and relax.
		time.Sleep(time.Duration(ExchangeGap) * time.Second)
	}
}

func srvToSrvForward(c *gin.Context) {
	// If no valid server ID or no known server, do not forward.
	if ServerId == "" || KnownServers == nil {
		return
	}
	// Prepare a new header.
	newHeader := c.Request.Header.Clone()
	if newHeader == nil {
		newHeader = make(http.Header)
	}
	// Handle TTL.
	curTtlStr := newHeader.Get(KEY_AKBP_TTL)
	if curTtlStr == "" {
		newHeader.Add(KEY_AKBP_TTL, strconv.Itoa(ADD_TTL))
	} else {
		curTtlInt, err := strconv.Atoi(curTtlStr)
		if err != nil {
			newHeader.Set(KEY_AKBP_TTL, strconv.Itoa(ADD_TTL))
		} else {
			if curTtlInt <= 0 {
				return
			}
			newHeader.Set(KEY_AKBP_TTL, strconv.Itoa(curTtlInt-1))
		}
	}
	// Remove additional header.
	newHeader.Del(KEY_CONTENT_TYPE)
	// Add this server to chain.
	newHeader.Add(KEY_AKBP_DO_NOT_FORWARD_TO, ServerId)
	// Set type and server ID info.
	newHeader.Set(KEY_AKBP_MSG_TYPE, "forwarded")
	newHeader.Set(KEY_AKBP_SERVER_ID, ServerId)
	// For safety.
	encodedMsg := url.QueryEscape(c.PostForm("msg"))
	// New k-v pair.
	newKv := map[string]string{"msg": encodedMsg}
	// Get not forward to. Use new header will let the server itself be included.
	notForwardTo := strings.Split(newHeader.Get(KEY_AKBP_DO_NOT_FORWARD_TO), ",")
	// Range known servers.
	for k, v := range KnownServers {
		// If the server is contained in notForwardTo.
		if slices.Contains(notForwardTo, k) {
			continue
		}
		// Set new auth info.
		newHeader.Set(KEY_AKBP_AUTH, KeysForAuth[k])
		// Post it!
		doPost(v, &newHeader, &newKv)
	}
}

func beaconReportHandler(c *gin.Context) {
	// Get beacon ID.
	bid := c.GetHeader(KEY_AKBP_BEACON_ID)
	if len(bid) == 0 || !ChkStrNoExit(&bid, "@.") || len(bid) > 32 {
		c.Header(KEY_AKBP_MSG_TYPE, "error")
		c.String(http.StatusBadRequest, ERR_WRONG_BEACON_ID_OR_AUTH_INFO)
		return
	}
	// Get key.
	key := c.GetHeader(KEY_AKBP_AUTH)
	if len(key) == 0 {
		c.Header(KEY_AKBP_MSG_TYPE, "error")
		c.String(http.StatusUnauthorized, ERR_NO_VALID_AUTH_HEADER)
		return
	}
	// Auth here.
	if !AuthOK(TABLE_BEACONS, bid, key) {
		c.Header(KEY_AKBP_MSG_TYPE, "error")
		c.String(http.StatusForbidden, ERR_WRONG_BEACON_ID_OR_AUTH_INFO)
		return
	}
	// To process the message.
	msgType := c.GetHeader(KEY_AKBP_MSG_TYPE)
	switch msgType {
	case "ping": // This ping can use as a kind of auth test.
		pong(c)
	case "beacon-report":
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
		// Add new record.
		tmp := strings.Split(c.GetHeader(KEY_AKBP_DO_NOT_FORWARD_TO), ",")
		notForwardTo := ""
		for _, x := range tmp {
			notForwardTo += strings.TrimSpace(x) + ","
		}
		notForwardTo += ServerId
		if AddRecord(bid, eid, ts, c.PostForm("msg"), STR_BEACON+"-"+bid, notForwardTo) {
			// Forward.
			srvToSrvForward(c)
		}
		// Report.
		reportOk(c)
	default:
		c.Header(KEY_AKBP_MSG_TYPE, "error")
		c.String(http.StatusBadRequest, ERR_UNKNOWN_MSG_TYPE)
	}
}

func fromServerHandler(c *gin.Context) {
	sid := c.GetHeader(KEY_AKBP_SERVER_ID)
	key := c.GetHeader(KEY_AKBP_AUTH)
	if !ChkStrNoExit(&sid, "") || len(sid) > 16 || !AuthOK(TABLE_SERVERS, sid, key) {
		return
	}
	msgType := c.GetHeader(KEY_AKBP_MSG_TYPE)
	switch msgType {
	case "ping":
		pong(c)
	case "forwarded":
		bid := c.GetHeader(KEY_AKBP_BEACON_ID)
		eid := c.GetHeader(KEY_AKBP_EVENT_ID)
		tsStr := c.GetHeader(KEY_AKBP_TIMESTAMP)
		tsInt, err := strconv.Atoi(tsStr)
		msg := c.PostForm("msg")
		if bid == "" || eid == "" || tsStr == "" || err != nil || (!ChkStrNoExit(&bid, "@.")) || len(bid) > 32 {
			return
		}
		tmp := strings.Split(c.GetHeader(KEY_AKBP_DO_NOT_FORWARD_TO), ",")
		notForwardTo := ""
		for _, x := range tmp {
			notForwardTo += strings.TrimSpace(x) + ","
		}
		notForwardTo += ServerId
		if AddRecord(bid, eid, tsInt, msg, STR_SERVER+"-"+sid, notForwardTo) {
			srvToSrvForward(c)
		}
		reportOk(c)
	default:
		c.Header(KEY_AKBP_MSG_TYPE, "error")
		c.String(http.StatusBadRequest, ERR_UNKNOWN_MSG_TYPE)
	}
}

func main() {
	initialize()
	LoadConf()
	if !FileExists(Db) {
		LogWarnln("Specified DB not found, will create and init a new one.")
		DbInit()
	}
	// If no valid server ID.
	if ServerId == "" {
		LogWarnln("No valid server ID, server to server forward will be disabled.")
	}

	// If no known servers.
	if KnownServers == nil {
		LogWarnln("No known server, server to server forward will be disabled.")
	}

	ginEng := gin.Default()

	// FAVICON //
	ginEng.GET("/favicon.ico", func(c *gin.Context) {
		c.Header(KEY_AKBP_MSG_TYPE, "favicon")
		c.String(http.StatusOK, "")
	})

	// Ping //
	ginEng.GET("/ping", func(c *gin.Context) {
		pong(c)
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

	// From-Server //
	ginEng.POST("/from-server", fromServerHandler)

	// Start Server To Server Exchange //
	go srvToSrvExchange()

	LogInfoln("Everything is OK, ready to start the HTTP(S) service.")
	ginEng.Run(Addr)
}
