/*
 * @Author: FunctionSir
 * @Date: 2023-07-17 22:47:42
 * @LastEditTime: 2023-07-26 23:18:10
 * @LastEditors: FunctionSir
 * @Description: Public consts, vars, and functions of AKBP Server.
 * @FilePath: /AKBP/server/public.go
 */
package main

import (
	"strconv"
)

const (
	DEBUG                 bool   = true               //Debug flag, for debugging or developing purposes ONLY!
	VER                   string = "0.1-alpha"        // Version.
	CODENAME              string = "Capsule"          // Code name of this version.
	SPLIT_LINE            string = "----------------" // Set the split line you want to use here.
	DEFAULT_PORT          int    = 3690               // Default port (=(ord("A")+ord("K")+ord("B")+ord("S"))*10).
	DEFAULT_BEACONS_DB    string = "beacons.db"       // Default beacons DB.
	DEFAULT_RCVR_LOG_FILE string = "akbo-rcvr.log"    // Default [A]nti [K]idnapping [B]eacon [O]rganization [-] [R]e[C]ei[V]e[R]('s/s').[Log] file.
)

var (
	ProgName    string = ""                    //Program name in os.Args[:].
	Port        int    = DEFAULT_PORT          //Server port.
	BeaconsDB   string = ""                    //Beacons DB.
	RcvrLogFile string = DEFAULT_RCVR_LOG_FILE //RcvrLog file.
	API_VER_AVL        = [...]string{"APIv1"}  // API version(s) available.
)

func Is_float_or_int(str string) bool {
	_, err := strconv.ParseFloat(str, 64)
	if err == nil {
		return true
	} else {
		return false
	}
}
