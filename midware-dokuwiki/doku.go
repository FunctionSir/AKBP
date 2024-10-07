/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2024-09-30 23:07:57
 * @LastEditTime: 2024-10-07 18:03:46
 * @LastEditors: FunctionSir
 * @Description: -
 * @FilePath: /AKBP/midware-dokuwiki/doku.go
 */
package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"gopkg.in/ini.v1"
)

type NotIniPart struct {
	RowId  int
	Bid    string
	Eid    string
	Ts     int
	Msg    string
	Origin string
	IsIni  bool
}

type DecodedPart struct {
	Bid     string
	Eid     string
	EncType string
	Pos     Coord
	Extra   string
}

func coordToStrLatFirstForHuman(coord *Coord, encrypted bool) string {
	markL := ""
	markR := ""
	if encrypted {
		markL = "<decrypt>"
		markR = "</decrypt>"
	}
	return markL + coord.Lat + markR + ", " + markL + coord.Lon + markR + ", " + markL + coord.Asl + markR
}

func genBbox(point *Coord) string {
	tmpLat, err := strconv.ParseFloat(point.Lat, 64)
	if err != nil {
		return ""
	}
	tmpLon, err := strconv.ParseFloat(point.Lon, 64)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%f,%f,%f,%f", tmpLon-GEN_BBOX_LON_OFFSET, tmpLat-GEN_BBOX_LAT_OFFSET, tmpLon+GEN_BBOX_LON_OFFSET, tmpLat+GEN_BBOX_LAT_OFFSET)
}

func NonIniFillTemplate(entry *NotIniPart, template *string) string {
	result := *template
	result = strings.Replace(result, "%BEACON_ID%", entry.Bid, -1)
	result = strings.Replace(result, "%EVENT_ID%", entry.Eid, -1)
	result = strings.Replace(result, "%TIMESTAMP_MS%", strconv.Itoa(entry.Ts), -1)
	result = strings.Replace(result, "%TIME%", time.UnixMilli(int64(entry.Ts)).Format(TimeTemplate), -1)
	result = strings.Replace(result, "%ORIGIN%", entry.Origin, -1)
	result = strings.Replace(result, "%VERSION%", VER, -1)
	result = strings.Replace(result, "%CODENAME%", CODENAME, -1)
	result = strings.Replace(result, "%GOOS%", runtime.GOOS, -1)
	result = strings.Replace(result, "%GOARCH%", runtime.GOARCH, -1)
	result = strings.Replace(result, "%ROW_ID%", strconv.Itoa(entry.RowId), -1)
	result = strings.Replace(result, "%RAW_MSG%", entry.Msg, -1)
	return result
}

func IniFillTemplate(notIni *NotIniPart, decoded *DecodedPart, template *string) string {
	result := *template
	result = NonIniFillTemplate(notIni, &result)
	result = strings.Replace(result, "%ENC_TYPE%", decoded.EncType, -1)
	encMsgLabelL := ""
	encMsgLabelR := ""
	guardL := "<file>\n"
	guardR := "\n</file>"
	// This is why you need "Encrypted Passwords Plugin" installed on dokuwiki.
	if decoded.EncType == "aes-256-cbc-pbkdf2" {
		encMsgLabelL = "<decrypt>"
		encMsgLabelR = "</decrypt>"
		guardL = ""
		guardR = ""
	}
	result = strings.Replace(result, "%LON%", encMsgLabelL+decoded.Pos.Lon+encMsgLabelR, -1)
	result = strings.Replace(result, "%LAT%", encMsgLabelL+decoded.Pos.Lat+encMsgLabelR, -1)
	result = strings.Replace(result, "%ASL%", encMsgLabelL+decoded.Pos.Asl+encMsgLabelR, -1)
	result = strings.Replace(result, "%EXTRA%", guardL+encMsgLabelL+decoded.Extra+encMsgLabelR+guardR, -1)
	// If position is encrypted or illegal, replace OSM related keys with "" (Empty string).
	if decoded.EncType != "plaintext" || decoded.Pos.Lat == "" || decoded.Pos.Lat == "nil" || decoded.Pos.Lon == "" || decoded.Pos.Lon == "nil" {
		osmKeys := []string{"%OSM_TO%", "%OSM_QUERY%", "%OSM_MAP_WITH_MARKER%", "%OSM_EMBED_WITH_MARKER%"}
		for _, x := range osmKeys {
			result = strings.Replace(result, x, "", -1)
		}
		return result
	}
	result = strings.Replace(result, "%OSM_TO%", fmt.Sprintf(
		"https://www.openstreetmap.org/directions?to=%s%%2C%s#map=%s/%s/%s",
		decoded.Pos.Lat, decoded.Pos.Lon, OSM_ZOOM, decoded.Pos.Lat, decoded.Pos.Lon), -1)
	result = strings.Replace(result, "%OSM_QUERY%", fmt.Sprintf(
		"https://www.openstreetmap.org/query?lat=%s&lon=%s#map=%s/%s/%s",
		decoded.Pos.Lat, decoded.Pos.Lon, OSM_ZOOM, decoded.Pos.Lat, decoded.Pos.Lon), -1)
	result = strings.Replace(result, "%OSM_MAP_WITH_MARKER%", fmt.Sprintf(
		"https://www.openstreetmap.org/?mlat=%s&mlon=%s#map=%s/%s/%s",
		decoded.Pos.Lat, decoded.Pos.Lon, OSM_ZOOM, decoded.Pos.Lat, decoded.Pos.Lon), -1)
	result = strings.Replace(result, "%OSM_EMBED_WITH_MARKER%", fmt.Sprintf(
		"https://www.openstreetmap.org/export/embed.html?bbox=%s&marker=%s%%2C%s",
		url.QueryEscape(genBbox(&decoded.Pos)), decoded.Pos.Lat, decoded.Pos.Lon), -1)
	return result
}

func GenAllSummary(file string, notIniPart []NotIniPart) {
	content := ""
	content += "====== Summary Page Of All Messages The Server Stored ======\n"
	content += fmt.Sprintf("** Total Messages Received: %d. ** \\\\\n", len(notIniPart))
	content += "** The latest one is the top one in this table. ** \\\\\n"
	content += "|  INI  |  Row ID  |  Beacon ID  |  Event ID  |  Reported Time  |  Received From  |  Details  |\n"
	for i := len(notIniPart) - 1; i >= 0; i-- {
		x := &notIniPart[i]
		isIniStr := "❌"
		if x.IsIni {
			isIniStr = "✔️"
		}
		content += fmt.Sprintf("|  %s  |  %s  |  %s  |  %s  |  %s  |  %s  |  [[%s|GO ➡️]]  |\n",
			isIniStr, strconv.Itoa(x.RowId), x.Bid, x.Eid,
			time.UnixMilli(int64(x.Ts)).Format(TimeTemplate), x.Origin,
			strings.ToLower(strconv.Itoa(x.RowId)+"-"+x.Bid+"-"+x.Eid+"-"+strconv.Itoa(x.Ts)))
	}
	content += "** The gap between two updates is " + strconv.Itoa(UpdGap) + " second(s). ** \\\\\n"
	content += "\n----\n// Generated by AKBP Midware for DokuWiki //"
	os.WriteFile(file, []byte(content), 0644)
}

func GenByEvent(file string, event string, pages []string, indexNotIni []int, indexDecoded []int, notIniPart []NotIniPart, decodedPart []DecodedPart) {
	content := ""
	content += "~~NOTOC~~\n"
	content += "====== Summary Page Of Event " + event + " ======\n"
	content += "===== Basic Info =====\n"
	content += fmt.Sprintf("Reported Time: From: %s To: %s \\\\\n",
		time.UnixMilli(int64(notIniPart[indexNotIni[0]].Ts)).Format(TimeTemplate),
		time.UnixMilli(int64(notIniPart[indexNotIni[len(indexDecoded)-1]].Ts)).Format(TimeTemplate))
	content += fmt.Sprintf("Total related pages (received messages): %d \\\\\n", len(pages))
	content += "===== KML Download =====\n"
	content += fmt.Sprintf("** -> Download KML File Of Route We Decoded And Assembled: {{%s:%s.kml| 💾 Download}} <- ** \\\\\n", KmlsNs, strings.ToLower(event))
	content += "P.S. KML file might be unavailable if there are no points we can decode in some cases. \\\\\n"
	content += "===== Related Pages =====\n"
	content += "** The latest one is the top one in the table. ** \\\\\n"
	content += "|  INI  |  Row ID  |  Encrypted  |  Decoded Position  |  Reported Time  |  Received From  |  Details  |\n"
	for i := len(pages) - 1; i >= 0; i-- {
		tmpNotIni := &notIniPart[indexNotIni[i]]
		isIniStr := "❌"
		isEncrypted := "❓"
		decodePos := "UNABLE TO DECODE"
		if tmpNotIni.IsIni && indexDecoded[i] >= 0 {
			isIniStr = "✔️"
			tmpDecoded := &decodedPart[indexDecoded[i]]
			if tmpDecoded.EncType != "plaintext" {
				isEncrypted = "✔️"
			} else {
				isEncrypted = "❌"
			}
			tmpEencrypted := false
			if tmpDecoded.EncType == "aes-256-cbc-pbkdf2" {
				tmpEencrypted = true
			}
			decodePos = coordToStrLatFirstForHuman(&tmpDecoded.Pos, tmpEencrypted)
		}
		rowId := tmpNotIni.RowId
		reportedTime := time.UnixMilli(int64(tmpNotIni.Ts)).Format(TimeTemplate)
		content += fmt.Sprintf("|  %s  |  %d  |  %s  |  %s  |  %s  |  %s  |  [[%s:%s|GO ➡️]]  |\n",
			isIniStr, rowId, isEncrypted, decodePos, reportedTime, tmpNotIni.Origin, AllEntriesNs, pages[i])
	}
	content += "** The gap between two updates is " + strconv.Itoa(UpdGap) + " second(s). ** \\\\\n"
	content += "\n----\n// Generated by AKBP Midware for DokuWiki //"
	os.WriteFile(file, []byte(content), 0644)
}

func GenByEventSummary(file string, events []string, ts []int, msgCnt []int) {
	content := ""
	content += "====== Summary Page Of All Events The Server Stored ======\n"
	content += fmt.Sprintf("** Total Events Calculated: %d. ** \\\\\n", len(events))
	content += "** The latest one is the top one in this table. ** \\\\\n"
	content += "|  #  |  Event  |  Last Message  |  Total Received  |  Details  |\n"
	for i := range events {
		content += fmt.Sprintf("|  %d  |  %s  |  %s  |  %d  |  [[%s:%s|GO ➡️]]|\n",
			i+1, events[i], time.UnixMilli(int64(ts[i])).Format(TimeTemplate), msgCnt[i], EventsNs, strings.ToLower(events[i]))
	}
	content += "** The gap between two updates is " + strconv.Itoa(UpdGap) + " second(s). ** \\\\\n"
	content += "\n----\n// Generated by AKBP Midware for DokuWiki //"
	os.WriteFile(file, []byte(content), 0644)
}

func ProcessRegRequests(dir string) (int, int) {
	regCnt := 0
	processCnt := 0
	files := make([]string, 0)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if info.Name() == "akbp_beacon_registration.txt" {
			files = append(files, path)
		}
		return nil
	})
	for _, x := range files {
		processCnt++
		conf, err := ini.Load(x)
		if err != nil {
			continue
		}
		sec := conf.Section("DEFAULT")
		if !sec.HasKey("BID") || !sec.HasKey("KEY") {
			continue
		}
		bid := strings.TrimSpace(sec.Key("BID").String())
		if !ChkStrNoExit(&bid) || len(bid) > 32 {
			continue
		}
		bid += "@" + Domain
		key := strings.TrimSpace(sec.Key("KEY").String())
		note := strings.TrimSpace(sec.Key("NOTE").String())
		if IsBidExists(bid) {
			continue
		}
		if RegBeacon(bid, key, note) {
			os.Remove(x)
			regCnt++
		}
	}
	return processCnt, regCnt
}
