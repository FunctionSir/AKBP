/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2023-07-16 00:26:53
 * @LastEditTime: 2023-07-26 23:07:29
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
)

func Apiv1_handler(w http.ResponseWriter, r *http.Request) {
	// Get and check uuid.
	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		fmt.Fprintln(w, "ERR::UUID_NOT_FOUND")
		return
	}
	if len(uuid) != 36 || strings.Count(uuid, "-") != 4 { //Just a simple check
		fmt.Fprintln(w, "ERR::ILLEGAL_UUID")
		return
	}
	// Get and Check key.
	key := r.URL.Query().Get("key")
	if key == "" {
		fmt.Fprintln(w, "ERR::KEY_NOT_FOUND")
		return
	}
	// Get and check sid.
	sid := r.URL.Query().Get("sid")
	if (sid != "") && (!Is_float_or_int(sid)) {
		fmt.Fprintln(w, "ERR::ILLEGAL_SID")
		return
	}
	// Get eid.
	eid := r.URL.Query().Get("eid")
	//Get and check ts.
	ts := r.URL.Query().Get("ts")
	if (ts != "") && (!Is_float_or_int(ts)) {
		fmt.Fprintln(w, "ERR::ILLEGAL_TS")
		return
	}
	// Get and check lat.
	lat := r.URL.Query().Get("lat")
	if (lat != "") && (!Is_float_or_int(lat)) {
		fmt.Fprintln(w, "ERR::ILLEGAL_LAT")
		return
	} else {
		if lat != "" {
			var tmp float64 = 0.0
			tmp, err := strconv.ParseFloat(lat, 64)
			if err != nil || math.Abs(tmp) > 90 {
				fmt.Fprintln(w, "ERR::ILLEGAL_LAT")
				return
			}
		}
	}
	// Get and check lon.
	lon := r.URL.Query().Get("lon")
	if (lon != "") && (!Is_float_or_int(lon)) {
		fmt.Fprintln(w, "ERR::ILLEGAL_LON")
		return
	} else {
		if lon != "" {
			var tmp float64 = 0.0
			tmp, err := strconv.ParseFloat(lon, 64)
			if err != nil || math.Abs(tmp) > 180 {
				fmt.Fprintln(w, "ERR::ILLEGAL_LON")
				return
			}
		}
	}
	// Get and check alt.
	alt := r.URL.Query().Get("alt")
	if (alt != "") && (!Is_float_or_int(alt)) {
		fmt.Fprintln(w, "ERR::ILLEGAL_ALT")
		return
	}
	// Get msg.
	msg := r.URL.Query().Get("msg")
	// Get and check img.
	img := r.URL.Query().Get("img")
	if (img != "") && (strings.Count(img, "data:image/") != 1 || strings.Count(img, ";base64,") != 1) {
		fmt.Fprintln(w, "ERR::ILLEGAL_IMG")
		return
	}
	// Get and check file.
	file := r.URL.Query().Get("file")
	if (file != "") && (strings.Count(file, "{") != 1 || strings.Count(file, "}") != 1) {
		fmt.Fprintln(w, "ERR::ILLEGAL_FILE")
		return
	}
	fmt.Println(time.Now().String()+" |", uuid, key, sid, eid, ts, lat, lon, alt, msg, img, file)
	fmt.Fprintln(w, "INFO::SUCCESS")
	if DEBUG {
		fmt.Fprintln(w, "----------------")
		fmt.Fprintln(w, "uuid = "+uuid)
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
		fmt.Fprintln(w, "----------------")
		fmt.Fprintln(w, "[A]nti [K]idnapping [B]eacon [P]roject Server")
		fmt.Fprintln(w, "Version: "+VER+", Codename: "+CODENAME)
		fmt.Fprintln(w, "Version of this API: APIv1")
	}
}
