/*
 * @Author: FunctionSir
 * @Date: 2023-07-17 22:47:42
 * @LastEditTime: 2023-08-01 02:40:41
 * @LastEditors: FunctionSir
 * @Description: Public consts, vars, and functions of AKBP Server.
 * @FilePath: /AKBP/server/public.go
 */
package main

import (
	"bufio"
	"crypto/sha512"
	"encoding/hex"
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
	BeaconsDB          string = DEFAULT_BEACONS_DB    //Beacons DB.
	BeaconsDBLines            = []string{}            //Lines in beacons DB.
	BeaconUUIDs               = []string{}            //UUIDs of beacons.
	BeaconSaltPosOfsts        = []int{}               //Beacon Salt Position Offsets, Specially, -1=to add salt @ the end of the key.
	BeaconSalts               = []string{}            //Salts of beacons.
	BeaconKPSHashes           = []string{}            //Hashes of beacons' [K]ey +([P]lus) Salt.
	RcvrLogFile        string = DEFAULT_RCVR_LOG_FILE //RcvrLog file.
	API_VER_AVL               = [...]string{"APIv1"}  // API version(s) available.
)

func Err_handle(where string, err error) bool {
	if err != nil {
		fmt.Println(time.Now().String() + " [E] Error occurred at " + where + ": " + err.Error() + ".")
		return true
	} else {
		if DEBUG {
			fmt.Println(time.Now().String() + " [I] At " + where + ": an operate was successfully done with err == nil.")
		}
		return false
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

func Find_str(source []string, target string) int {
	for i := 0; i < len(source); i++ {
		if source[i] == target {
			return i
		}
	}
	return -1
}

func Gen_KPS(key string, offset int, salt string) string {
	var kps string = ""
	switch offset {
	case 0:
		kps = salt + key
	case -1:
		kps = key + salt
	default:
		if offset > len(key) {
			kps = key + salt
		} else {
			kps = key[:offset-1] + salt + key[offset:]
		}
	}
	return kps
}

func Gen_KPS_hash(key string, offset int, salt string) string {
	kps := Gen_KPS(key, offset, salt)
	h := sha512.New()
	h.Write([]byte(kps))
	return hex.EncodeToString(h.Sum(nil))
}
