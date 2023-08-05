/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2023-07-16 00:26:53
 * @LastEditTime: 2023-08-05 22:31:15
 * @LastEditors: FunctionSir
 * @Description: APIv1 related funcs.
 * @FilePath: /AKBP/server/apiv1.go
 */
package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	HANDLER_NAME string = "Apiv1_handler"
)

func gen_log_entry(bUUID, sid, eid, ts, lat, lon, alt, msg, img, file string) (string, []error) {
	var errs = []error{}
	var lines = []string{}
	logEntryUUID := uuid.New().String()
	lines = append(lines, "["+logEntryUUID+"]")
	if !DEBUG {
		lines = append(lines, "Handler = Apiv1_handler")
	} else {
		lines = append(lines, "Handler = (DEBUG)Apiv1_handler")
	}
	lines = append(lines, "Puuid = "+Remove_CR_and_LF(bUUID))
	lines = append(lines, "Psid = "+Remove_CR_and_LF(sid))
	lines = append(lines, "Peid = "+Remove_CR_and_LF(eid))
	lines = append(lines, "Pts = "+Remove_CR_and_LF(ts))
	lines = append(lines, "Plat = "+Remove_CR_and_LF(lat))
	lines = append(lines, "Plon = "+Remove_CR_and_LF(lon))
	lines = append(lines, "Palt = "+Remove_CR_and_LF(alt))
	lines = append(lines, "Pmsg = "+Remove_CR_and_LF(msg))
	lines = append(lines, "Pimg = "+Remove_CR_and_LF(img))
	lines = append(lines, "Pfile = "+Remove_CR_and_LF(file))
	errs = append(errs, Write_lines(RcvrLogFile, lines)...)
	return logEntryUUID, errs
}

func Apiv1_handler(w http.ResponseWriter, r *http.Request) {
	// Get and check bUUID.
	bUUID := r.URL.Query().Get("uuid")
	if bUUID == "" {
		fmt.Fprintln(w, "ERR::UUID_NOT_FOUND")
		if !DEBUG {
			return
		}
	}
	if len(bUUID) != 36 || strings.Count(bUUID, "-") != 4 { //Just a simple check
		fmt.Fprintln(w, "ERR::ILLEGAL_UUID")
		if !DEBUG {
			return
		}
	}
	beaconIndex := Find_str(BeaconUUIDs, bUUID)
	if beaconIndex == -1 {
		fmt.Fprintln(w, "ERR::UUID_NOT_REGISTERED")
		if !DEBUG {
			return
		} else {
			beaconIndex = 0
			if len(BeaconUUIDs) == 0 {
				BeaconUUIDs = append(BeaconUUIDs, "")
				BeaconSalts = append(BeaconSalts, "")
				BeaconSaltPosOfsts = append(BeaconSaltPosOfsts, 0)
				BeaconKPSHashes = append(BeaconKPSHashes, "")
			}
		}
	}
	// Get and Check key.
	key := r.URL.Query().Get("key")
	if key == "" {
		fmt.Fprintln(w, "ERR::KEY_NOT_FOUND")
		if !DEBUG {
			return
		}
	}
	if Gen_KPS_hash(key, BeaconSaltPosOfsts[beaconIndex], BeaconSalts[beaconIndex]) != BeaconKPSHashes[beaconIndex] {
		fmt.Fprintln(w, "ERR::WRONG_KEY")
		if !DEBUG {
			return
		}
	}
	// Get and check sid.
	sid := r.URL.Query().Get("sid")
	if (sid != "") && (!Is_float_or_int(sid)) {
		fmt.Fprintln(w, "ERR::ILLEGAL_SID")
		if !DEBUG {
			return
		}
	}
	// Get eid.
	eid := r.URL.Query().Get("eid")
	//Get and check ts.
	ts := r.URL.Query().Get("ts")
	if (ts != "") && (!Is_float_or_int(ts)) {
		fmt.Fprintln(w, "ERR::ILLEGAL_TS")
		if !DEBUG {
			return
		}
	}
	// Get and check lat.
	lat := r.URL.Query().Get("lat")
	if (lat != "") && (!Is_float_or_int(lat)) {
		fmt.Fprintln(w, "ERR::ILLEGAL_LAT")
		if !DEBUG {
			return
		}
	} else {
		if lat != "" {
			var tmp float64 = 0.0
			tmp, err := strconv.ParseFloat(lat, 64)
			if err != nil || math.Abs(tmp) > 90 {
				fmt.Fprintln(w, "ERR::ILLEGAL_LAT")
				if !DEBUG {
					return
				}
			}
		}
	}
	// Get and check lon.
	lon := r.URL.Query().Get("lon")
	if (lon != "") && (!Is_float_or_int(lon)) {
		fmt.Fprintln(w, "ERR::ILLEGAL_LON")
		if !DEBUG {
			return
		}
	} else {
		if lon != "" {
			var tmp float64 = 0.0
			tmp, err := strconv.ParseFloat(lon, 64)
			if err != nil || math.Abs(tmp) > 180 {
				fmt.Fprintln(w, "ERR::ILLEGAL_LON")
				if !DEBUG {
					return
				}
			}
		}
	}
	// Get and check alt.
	alt := r.URL.Query().Get("alt")
	if (alt != "") && (!Is_float_or_int(alt)) {
		fmt.Fprintln(w, "ERR::ILLEGAL_ALT")
		if !DEBUG {
			return
		}
	}
	// Get msg.
	msg := r.URL.Query().Get("msg")
	// Get and check img.
	img := r.URL.Query().Get("img")
	if (img != "") && (strings.Count(img, "data:image/") != 1 || strings.Count(img, ";base64,") != 1) {
		fmt.Fprintln(w, "ERR::ILLEGAL_IMG")
		if !DEBUG {
			return
		}
	}
	// Get and check file.
	file := r.URL.Query().Get("file")
	if (file != "") && (strings.Count(file, "{") != 1 || strings.Count(file, "}") != 1) {
		fmt.Fprintln(w, "ERR::ILLEGAL_FILE")
		if !DEBUG {
			return
		}
	}
	if DEBUG {
		fmt.Println(time.Now().String() + " [D] GET request received: uuid=" + bUUID + ", key=" + key + ", sid=" + sid + ", eid=" + eid + ", ts=" + ts + ", lat=" + lat + ", lon=" + lon + ", alt=" + alt + ", msg=" + msg + ", img=" + img + ", file=" + file + ".")
	}
	logEntryUUID, errs := gen_log_entry(bUUID, sid, eid, ts, lat, lon, alt, msg, img, file)
	if len(errs) != 0 {
		fmt.Fprintln(w, "ERR::FAILED_TO_GEN_LOG_ENTRY")
		if !DEBUG {
			return
		}
	}
	if !DEBUG {
		fmt.Fprintln(w, "INFO::SUCCESS")
	} else {
		fmt.Fprintln(w, "WARN::DEBUG_MODE_ON")
	}
	if DEBUG {
		fmt.Fprintln(w, SPLIT_LINE)
		fmt.Fprintln(w, "uuid = "+bUUID)
		fmt.Fprintln(w, "key = "+key)
		fmt.Fprintln(w, "sid = "+sid)
		fmt.Fprintln(w, "eid = "+eid)
		fmt.Fprintln(w, "ts = "+ts)
		fmt.Fprintln(w, "lat = "+lat)
		fmt.Fprintln(w, "lon = "+lon)
		fmt.Fprintln(w, "alt = "+alt)
		fmt.Fprintln(w, "msg = "+msg)
		fmt.Fprintln(w, "img = "+img)
		fmt.Fprintln(w, "file = "+file)
		fmt.Fprintln(w, SPLIT_LINE)
		fmt.Fprintln(w, "saltPosOfst = "+strconv.Itoa(BeaconSaltPosOfsts[beaconIndex]))
		fmt.Fprintln(w, "salt = "+BeaconSalts[beaconIndex])
		fmt.Fprintln(w, "keyPlusSalt = "+Gen_KPS(key, BeaconSaltPosOfsts[beaconIndex], BeaconSalts[beaconIndex]))
		fmt.Fprintln(w, "keyPlusSaltHash = "+Gen_KPS_hash(key, BeaconSaltPosOfsts[beaconIndex], BeaconSalts[beaconIndex]))
		fmt.Fprintln(w, SPLIT_LINE)
		fmt.Fprintln(w, "logEntryUUID = "+logEntryUUID)
		if !DEBUG {
			fmt.Fprintln(w, "handler = Apiv1_handler")
		} else {
			fmt.Fprintln(w, "handler = (DEBUG)Apiv1_handler")
		}
		fmt.Fprintln(w, SPLIT_LINE)
		fmt.Fprintln(w, "[A]nti [K]idnapping [B]eacon [P]roject Server")
		fmt.Fprintln(w, "Version: "+VER+", Codename: "+CODENAME)
		fmt.Fprintln(w, "Version of this API: APIv1")
	}
}
