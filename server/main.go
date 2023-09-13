/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2023-07-14 23:11:35
 * @LastEditTime: 2023-09-03 14:38:53
 * @LastEditors: FunctionSir
 * @Description: Server of AKBP for beacons to link.
 * @FilePath: /AKBP/server/main.go
 */
package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Global consts and vars are in public.go.
// Some useful funcs are also in public.go.

// The os.Args parser.
func args_parser() {
	var i int
	ProgName = os.Args[0]
	for i = 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-p", "--port":
			var tmp int
			tmp, err := strconv.Atoi(os.Args[i+1])
			if err == nil {
				Port = tmp
			} else {
				Err_handle("main.args_parser", err)
			}
		case "-b", "--beacons-db":
			BeaconsDB = os.Args[i+1]
		case "-r", "--rcvr-log":
			RcvrLogFile = os.Args[i+1]
		}
	}
}

// Do this first.
func initial() {
	fmt.Println("[A]nti [K]idnapping [B]eacon [P]roject Server")
	fmt.Println("Version: " + VER + ", Codename: " + CODENAME)
	fmt.Println("This is a FOSS under the AGPLv3.")
	if DEBUG {
		fmt.Println(SPLIT_LINE)
		fmt.Println("!!! DEBUG flag was set to TRUE !!!")
	}
	fmt.Println(SPLIT_LINE)
	fmt.Println("Command = " + strings.Join(os.Args[:], " "))
	args_parser()
	fmt.Println(SPLIT_LINE)
	fmt.Println("ProgName = " + ProgName)
	fmt.Println("Port = " + strconv.Itoa(Port))
	fmt.Println("BeaconsDB = " + BeaconsDB)
	fmt.Println("RcvrLogFile = " + RcvrLogFile)
	fmt.Println(SPLIT_LINE)
	fmt.Println(time.Now().String() + " [I] Reading reg key file...")
	PrevRegKey = Read_lines(Reg_key_file)[0]
	CurrentRegKey = PrevRegKey
}

func beacons_db_reader() {
	fmt.Println(time.Now().String() + " [I] Reading the BeaconsDB...")
	BeaconsDBLines = Read_lines(BeaconsDB)
	validLines := len(BeaconsDBLines)
	for i := 0; i < len(BeaconsDBLines); i++ {
		tmp := strings.Split(BeaconsDBLines[i], " ")
		if len(tmp) == 4 || len(tmp) == 5 {
			bspoTmp, err := strconv.Atoi(tmp[1])
			if !Err_handle("main.beacons_db_reader", err) {
				BeaconUUIDs = append(BeaconUUIDs, tmp[0])
				BeaconSaltPosOfsts = append(BeaconSaltPosOfsts, bspoTmp)
				BeaconSalts = append(BeaconSalts, tmp[2])
				BeaconKPSHashes = append(BeaconKPSHashes, tmp[3])
			} else {
				fmt.Println(time.Now().String() + " [W] Ignored line #" + strconv.Itoa(i+1) + " of " + BeaconsDB + ": wrong format.")
				validLines--
			}
		} else {
			fmt.Println(time.Now().String() + " [W] Ignored line #" + strconv.Itoa(i+1) + " of " + BeaconsDB + ": wrong format.")
			validLines--
		}
	}
	fmt.Println(time.Now().String() + " [I] Done! Noticed " + strconv.Itoa(validLines) + " valid lines in " + BeaconsDB + ".")
}

// Feel free to customize it!
func default_handler(w http.ResponseWriter, r *http.Request) {
	_ = r // Don't need var "r" in this ver.
	fmt.Fprintln(w, "This is an [A]nti [K]idnapping [B]eacon [P]roject Server.")
	fmt.Fprintln(w, "Version: "+VER+", Codename: "+CODENAME)
	if DEBUG {
		fmt.Fprintln(w, "!WARNING! !THIS SOFTWARE IS UNDER DEVELOPING AND SHOULD NOT BE USED IN ANY FORMAL SERVERS! !WARNING!")
		fmt.Fprintln(w, "!WARNING! !ANY PACKAGE SENT TO HERE MAY NOT BE HANDLE CORRECTLY! !WARNING!")
	}
	fmt.Fprintln(w, "API Version(s) Available: "+strings.Join(API_VER_AVL[:], ", "))
}

func gen_reg_key() {
	var newRegKey string = ""
	for i := 0; i < Reg_key_STR; i++ {
		newRegKey = newRegKey + uuid.New().String()
	}
	Create_file(Reg_key_file, []string{newRegKey})
	PrevRegKey = CurrentRegKey
	CurrentRegKey = Read_lines(Reg_key_file)[0]
	time.Sleep(time.Duration(Reg_key_gen_gap) * time.Second)
}

func reg_handler(w http.ResponseWriter, r *http.Request) {
	regKey := r.URL.Query().Get("regKey")
	if regKey == "" {
		fmt.Fprintln(w, "ERR::REG_KEY_NOT_FOUND")
		if !DEBUG {
			return
		}
	}
	if (regKey != CurrentRegKey && regKey != PrevRegKey) && (regKey != "") {
		fmt.Fprintln(w, "ERR::WRONG_REG_KEY")
		if !DEBUG {
			return
		}
	}
	bcnKey := r.URL.Query().Get("bcnKey")
	if bcnKey == "" {
		fmt.Println(w, "ERR::BCN_KEY_NOT_FOUND")
		if !DEBUG {
			return
		}
	}
	if len(bcnKey) < MinBcnKeyLen {
		fmt.Fprintln(w, "ERR::BCN_KEY_TOO_SHORT")
		if !DEBUG {
			return
		}
	}
	bcnExtInfoB64 := base64.StdEncoding.EncodeToString([]byte(r.URL.Query().Get("bcnExtInfo")))
	bcnUUID := uuid.New().String()
	tmp, e := rand.Int(rand.Reader, big.NewInt(int64(len(bcnKey)+SaltPosOfstOvfl+1)))
	if Err_handle("main.reg_handler", e) && !DEBUG {
		return
	}
	bcnSaltPosOfst := int(tmp.Int64() - 1)
	var bcnSalt string = ""
	for i := 0; i < SaltSTR; i++ {
		bcnSalt = bcnSalt + uuid.New().String()
	}
	lines := []string{bcnUUID + " " + strconv.Itoa(bcnSaltPosOfst) + " " + bcnSalt + " " + Gen_KPS_hash(bcnKey, bcnSaltPosOfst, bcnSalt) + " " + bcnExtInfoB64}
	errs := Append_lines(BeaconsDB, lines)
	if len(errs) > 0 {
		for i := 0; i < len(errs); i++ {
			Err_handle("main.reg_handler", errs[i])
		}
		if !DEBUG {
			return
		}
	}
	fmt.Fprintln(w, bcnUUID)
	BeaconUUIDs = append(BeaconUUIDs, bcnUUID)
	BeaconSaltPosOfsts = append(BeaconSaltPosOfsts, bcnSaltPosOfst)
	BeaconSalts = append(BeaconSalts, bcnSalt)
	BeaconKPSHashes = append(BeaconKPSHashes, Gen_KPS_hash(bcnKey, bcnSaltPosOfst, bcnSalt))
	if !DEBUG {
		fmt.Fprintln(w, "INFO::SUCCESS")
	} else {
		fmt.Fprintln(w, "WARN::DEBUG_MODE_ON")
	}
}

func log_viewer_handler(w http.ResponseWriter, r *http.Request) {
	lines := Read_lines(RcvrLogFile)
	for i := 0; i < len(lines); i++ {
		fmt.Fprintln(w, lines[i])
	}
}

// HTTP(S) Server.
func http_server() {
	http.HandleFunc("/", default_handler)
	http.HandleFunc("/logViewer/", log_viewer_handler)
	http.HandleFunc("/reg/", reg_handler)
	http.HandleFunc("/v1/", Apiv1_handler)
	http.ListenAndServe(":"+strconv.Itoa(Port), nil)
}

// A common but uncommon func main.
func main() {
	initial()
	beacons_db_reader()
	go gen_reg_key()
	http_server()
}
