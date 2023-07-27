/*
 * @Author: FunctionSir
 * @Date: 2023-07-17 22:47:42
 * @LastEditTime: 2023-07-28 00:37:10
 * @LastEditors: FunctionSir
 * @Description: Public consts, vars, and functions of AKBP Server.
 * @FilePath: /AKBP/server/public.go
 */
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
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
	ProgName           string = ""                    //Program name in os.Args[:].
	Port               int    = DEFAULT_PORT          //Server port.
	BeaconsDB          string = ""                    //Beacons DB.
	BeaconsDBLines            = []string{}            //Lines in beacons DB.
	BeaconUUIDs               = []string{}            //UUIDs of beacons.
	BeaconSaltPosOfsts        = []int{}               //Beacon Salt Position Offsets, Specially, -1=to add salt @ the end of the key.
	BeaconSalts               = []string{}            //Salts of beacons.
	BeaconKPSHashes           = []string{}            //Hashes of beacons' [K]ey +([P]lus) Salt.
	RcvrLogFile        string = DEFAULT_RCVR_LOG_FILE //RcvrLog file.
	API_VER_AVL               = [...]string{"APIv1"}  // API version(s) available.
)

func Err_handle(where string, err error) {
	if err != nil {
		fmt.Println(time.Now().String() + " [E] Error occurred at " + where + ": " + err.Error() + ".")
	} else {
		if DEBUG {
			fmt.Println(time.Now().String() + " [I] At " + where + ": an operate was successfully done with err == nil.")
		}
	}
}

func Read_lines(name string) []string {
	var r = []string{}
	f, err := os.Open(name)
	Err_handle("main.Read_lines", err)
	fileScanner := bufio.NewScanner(f)
	for fileScanner.Scan() {
		r = append(r, fileScanner.Text())
	}
	return r
}

func Is_float_or_int(str string) bool {
	_, err := strconv.ParseFloat(str, 64)
	if err == nil {
		return true
	} else {
		return false
	}
}
