/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2023-07-16 00:26:53
 * @LastEditTime: 2023-07-18 02:34:12
 * @LastEditors: FunctionSir
 * @Description: APIv1 related funcs.
 * @FilePath: /AKBP/server/apiv1.go
 */
package main

import (
	"fmt"
	"net/http"
	"strings"
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
	if (sid != "") && (!Is_pure_num_str(sid)) {
		fmt.Fprintln(w, "ERR::ILLEGAL_SID")
		return
	}
	// Get eid.
	eid := r.URL.Query().Get("eid")
	//Get and check ts.
	ts := r.URL.Query().Get("ts")
	if (ts != "") && (!Is_pure_num_str(ts)) {
		fmt.Fprintln(w, "ERR::ILLEGAL_TS")
		return
	}
	// Get and check lat.
	lat := r.URL.Query().Get("lat")
	if (lat != "") && (!Is_pure_num_str(lat)) {
		fmt.Fprintln(w, "ERR::ILLEGAL_LAT")
		return
	}
	// Get and check lon.
	lon := r.URL.Query().Get("lon")
	if (lon != "") && (!Is_pure_num_str(lon)) {
		fmt.Fprintln(w, "ERR::ILLEGAL_LON")
		return
	}
	// Get and check alt.
	alt := r.URL.Query().Get("alt")
	if (alt != "") && (!Is_pure_num_str(alt)) {
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
	_, _, _, _, _, _, _, _, _ = sid, eid, ts, lat, lon, alt, msg, img, file
}
