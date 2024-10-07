/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2024-09-22 14:46:17
 * @LastEditTime: 2024-10-07 16:59:58
 * @LastEditors: FunctionSir
 * @Description: The bridge to dokuwiki.
 * @FilePath: /AKBP/midware-dokuwiki/main.go
 */

package main

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
	"gopkg.in/ini.v1"
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
	oldUmask := syscall.Umask(022)
	defer syscall.Umask(oldUmask)
	if runtime.GOOS == "linux" {
		TmpDir = path.Join(LINUX_DEV_SHM, DEFAULT_TMP_DIR)
	} else {
		TmpDir = DEFAULT_TMP_DIR
	}
	LoadConf()
	if !FileExists(Db) {
		LogFatalln("Database \"" + Db + "\" not exists. Please check your config.")
	}
	if len(TmpDir) < len(LINUX_DEV_SHM) || TmpDir[0:len(LINUX_DEV_SHM)] != LINUX_DEV_SHM {
		LogWarnln("A large amount of write operations might will be preformed. Please make sure you are using a ram disk, or a HDD, or a SSD but without any valuable data in it for tmp dir.")
		LogWarnln("Please make sure that dir " + TmpDir + " is not important, the program will clean the dir.")
	}
	err := os.RemoveAll(TmpDir)
	if err != nil {
		LogFatalln("Failed to remove tmp dir, unable to clean it.")
	}
	err = os.MkdirAll(TmpDir, 0755)
	if err != nil {
		LogFatalln("Failed to make tmp dir.")
	}
	neededDirs := []string{"/dw-ns", "/dw-m", "/dw-ns/all", "/dw-ns/events", "/dw-m/kmls"}
	for _, x := range neededDirs {
		err = os.Mkdir(path.Join(TmpDir, x), 0755)
		if err != nil {
			LogFatalln("Needed dir \"" + path.Join(TmpDir, x) + "\" could not be created.")
		}
	}
	content, err := os.ReadFile(NonIniTemplate)
	if err != nil {
		LogFatalln("Unable to load NoIniTemplate \"" + NonIniTemplate + "\". Error: " + err.Error() + ".")
	}
	NonIniTemplate = string(content)
	content, err = os.ReadFile(IniTypeTemplate)
	if err != nil {
		LogFatalln("Unable to load IniTypeTemplate \"" + IniTypeTemplate + "\". Error: " + err.Error() + ".")
	}
	IniTypeTemplate = string(content)
	LogInfoln("The software is ready to serve.")
	totUpd := 0
	eventsDir := path.Join(TmpDir, "/dw-ns/events")
	allDir := path.Join(TmpDir, "/dw-ns/all")
	kmlsDir := path.Join(TmpDir, "/dw-m/kmls")
	var processed int
	// The main loop.
	for {
		startTime := time.Now().UnixMilli()
		// Handle reg requests.
		totRegRequests, reged := ProcessRegRequests(UsersDir)
		processed = 0
		var rowId int
		var bid string
		var eid string
		var ts int
		var msg string
		var origin string
		toGenKml := make(map[string][]Coord)
		byEvent := make(map[string][]string)
		byEventNotIniIndex := make(map[string][]int)
		byEventDecodedPartIndex := make(map[string][]int)
		notIniParts := make([]NotIniPart, 0)
		decodedParts := make([]DecodedPart, 0)
		addedToByEventSummary := make(set[string])
		rows := QueryRecs()
		// Process every rows.
		for rows.Next() {
			processed++
			rows.Scan(&rowId, &bid, &eid, &ts, &msg, &origin)
			// Gen file in "all" dir.
			filePath := path.Join(allDir, strings.ToLower(strconv.Itoa(rowId)+"-"+bid+"-"+eid+"-"+strconv.Itoa(ts)+".txt"))
			ini, err := ini.Load([]byte(msg))
			notIni := NotIniPart{rowId, bid, eid, ts, msg, origin, true}
			// Append to by event list.
			byEvent[bid+"-"+eid] = append(byEvent[bid+"-"+eid], strings.ToLower(strconv.Itoa(rowId)+"-"+bid+"-"+eid+"-"+strconv.Itoa(ts)))
			// If is not ini format.
			if err != nil {
				notIni.IsIni = false
				os.WriteFile(filePath, []byte(NonIniFillTemplate(&notIni, &NonIniTemplate)), 0644)
				notIniParts = append(notIniParts, notIni)
				byEventNotIniIndex[bid+"-"+eid] = append(byEventNotIniIndex[bid+"-"+eid], len(notIniParts)-1)
				byEventDecodedPartIndex[bid+"-"+eid] = append(byEventDecodedPartIndex[bid+"-"+eid], -1)
				continue
			}
			sec := ini.Section("DEFAULT")
			decoded := DecodedPart{}
			decoded.Bid = notIni.Bid
			decoded.Eid = notIni.Eid
			if sec.HasKey("Encryption") {
				decoded.EncType = strings.TrimSpace(sec.Key("Encryption").String())
			} else {
				decoded.EncType = "plaintext"
			}
			if sec.HasKey("Position") {
				splited := strings.Split(sec.Key("Position").String(), ",")
				if len(splited) <= 1 {
					decoded.Pos = Coord{"nil", "nil", "nil"}
				} else {
					asl := "nil"
					if len(splited) >= 3 {
						asl = splited[2]
					}
					decoded.Pos = Coord{strings.TrimSpace(splited[0]), strings.TrimSpace(splited[1]), strings.TrimSpace(asl)}
				}
			} else {
				decoded.Pos = Coord{"nil", "nil", "nil"}
			}
			coordKeys := []string{"Lat", "Lon", "Asl"}
			for _, x := range coordKeys {
				if sec.HasKey(x) {
					decoded.Pos.Lon = strings.TrimSpace(sec.Key(x).String())
				}
			}
			if sec.HasKey("Extra") {
				decoded.Extra = strings.TrimSpace(sec.Key("Extra").String())
			} else {
				decoded.Extra = "nil"
			}
			notIniParts = append(notIniParts, notIni)
			decodedParts = append(decodedParts, decoded)
			byEventNotIniIndex[bid+"-"+eid] = append(byEventNotIniIndex[bid+"-"+eid], len(notIniParts)-1)
			byEventDecodedPartIndex[bid+"-"+eid] = append(byEventDecodedPartIndex[bid+"-"+eid], len(decodedParts)-1)
			os.WriteFile(filePath, []byte(IniFillTemplate(&notIni, &decoded, &IniTypeTemplate)), 0644)
		}
		// Gen Summary for "all".
		GenAllSummary(path.Join(allDir, "summary.txt"), notIniParts)
		// Gen KMLs.
		for _, x := range decodedParts {
			if x.EncType == "plaintext" {
				tmp := x.Bid + "-" + x.Eid
				toGenKml[tmp] = append(toGenKml[tmp], x.Pos)
			}
		}
		for i, x := range toGenKml {
			os.WriteFile(path.Join(kmlsDir, strings.ToLower(i+".kml")), []byte(GenRoute(i, x, "Route of "+i, "Belongs to event "+i+".")), 0644)
		}
		// Gen "by event"
		for i, x := range byEvent {
			GenByEvent(path.Join(eventsDir, strings.ToLower(i+".txt")), i, x, byEventNotIniIndex[i], byEventDecodedPartIndex[i], notIniParts, decodedParts)
		}
		// Gen summary of "by evnent".
		toGenSummaryOfByEventTime := make([]int, 0)
		toGenSummaryOfByEventEvent := make([]string, 0)
		toGenSummaryOfByEventMsgCnt := make([]int, 0)
		for i := len(notIniParts) - 1; i >= 0; i-- {
			tmp := notIniParts[i].Bid + "-" + notIniParts[i].Eid
			if !addedToByEventSummary.Have(tmp) {
				toGenSummaryOfByEventEvent = append(toGenSummaryOfByEventEvent, tmp)
				toGenSummaryOfByEventTime = append(toGenSummaryOfByEventTime, notIniParts[i].Ts)
				toGenSummaryOfByEventMsgCnt = append(toGenSummaryOfByEventMsgCnt, len(byEventNotIniIndex[tmp]))
				addedToByEventSummary.Ins(tmp)
			}
		}
		GenByEventSummary(path.Join(eventsDir, "summary.txt"), toGenSummaryOfByEventEvent, toGenSummaryOfByEventTime, toGenSummaryOfByEventMsgCnt)
		// Calc used time.
		usedTime := time.Now().UnixMilli() - startTime
		LogInfoln(fmt.Sprintf("Update #%d done. Processed %d records, %d registration request, registered %d beacons. Used %d ms.", totUpd+1, processed, totRegRequests, reged, usedTime))
		if usedTime > 5000 {
			LogWarnln("The program is running too slow. Check your DB, or the location of your tmp dir.")
		}
		// Add 1 to update counter.
		totUpd = totUpd + 1
		// Sit back and relax.
		time.Sleep(time.Duration(UpdGap) * time.Second)
	}
}
