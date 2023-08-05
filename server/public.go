/*
 * @Author: FunctionSir
 * @Date: 2023-07-17 22:47:42
 * @LastEditTime: 2023-08-05 22:28:03
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
	"io"
	"os"
	"strconv"
	"strings"
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
			fmt.Println(time.Now().String() + " [D] At " + where + ": an operate was successfully done with err == nil.")
		}
		return false
	}
}

func Remove_CR_and_LF(str string) string {
	str = strings.ReplaceAll(str, "\r", "")
	str = strings.ReplaceAll(str, "\n", "")
	return str
}

func Read_lines(name string) []string {
	var r = []string{}
	f, e := os.Open(name)
	Err_handle("main.Read_lines", e)
	defer func() {
		e := f.Close()
		c := 0
		for Err_handle("main.Read_lines", e) && c <= 8 {
			e = f.Close()
			c++
		}
	}()
	fileScanner := bufio.NewScanner(f)
	for fileScanner.Scan() {
		r = append(r, fileScanner.Text())
	}
	return r
}

func Write_lines(name string, lines []string) []error {
	var errs = []error{}
	s, e := os.Stat(name)
	if Err_handle("main.Write_lines", e) {
		errs = append(errs, e)
	}
	f, e := os.OpenFile(name, os.O_RDWR|os.O_APPEND, s.Mode())
	if Err_handle("main.Write_lines", e) {
		errs = append(errs, e)
	}
	defer func() {
		e := f.Close()
		c := 0
		for Err_handle("main.Write_lines", e) && c <= 8 {
			e = f.Close()
			c++
		}
	}()
	tmp := []byte{0}
	if s.Size() > 0 {
		_, e = f.Seek(-1, io.SeekEnd)
		if Err_handle("main.Write_lines", e) {
			errs = append(errs, e)
		}
		_, _ = f.Read(tmp)
	}
	if s.Size() != 0 && string(tmp) != "\n" {
		f.Seek(0, io.SeekEnd)
		if Err_handle("main.Write_lines", e) {
			errs = append(errs, e)
		}
		_, e := f.WriteString("\n")
		if Err_handle("main.Write_lines", e) {
			errs = append(errs, e)
		}
	} else {
		f.Seek(0, io.SeekEnd)
	}
	for i := 0; i < len(lines); i++ {
		if !strings.HasSuffix(lines[i], "\n") {
			lines[i] = lines[i] + "\n"
		}
		_, e := f.WriteString(lines[i])
		if Err_handle("main.Write_lines", e) {
			errs = append(errs, e)
		}
	}
	return errs
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
