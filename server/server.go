/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2023-07-14 23:11:35
 * @LastEditTime: 2023-08-05 22:14:38
 * @LastEditors: FunctionSir
 * @Description: Server of AKBP for beacons to link.
 * @FilePath: /AKBP/server/server.go
 */
package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
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
}

func beacons_db_reader() {
	fmt.Println(time.Now().String() + " [I] Reading the BeaconsDB...")
	BeaconsDBLines = Read_lines(BeaconsDB)
	validLines := len(BeaconsDBLines)
	for i := 0; i < len(BeaconsDBLines); i++ {
		tmp := strings.Split(BeaconsDBLines[i], " ")
		if len(tmp) == 4 {
			bspoTmp, err := strconv.Atoi(tmp[1])
			if !Err_handle("main.beacons_db_reader", err) {
				BeaconUUIDs = append(BeaconUUIDs, tmp[0])
				BeaconSaltPosOfsts = append(BeaconSaltPosOfsts, bspoTmp)
				BeaconSalts = append(BeaconSalts, tmp[2])
				BeaconKPSHashes = append(BeaconKPSHashes, tmp[3])
			} else {
				fmt.Println(time.Now().String() + " [W] Ignored line #" + strconv.Itoa(i+1) + " of " + BeaconsDB + ": wrong format.")
				validLines = validLines - 1
			}
		} else {
			fmt.Println(time.Now().String() + " [W] Ignored line #" + strconv.Itoa(i+1) + " of " + BeaconsDB + ": wrong format.")
			validLines = validLines - 1
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

func log_viewer_handler(w http.ResponseWriter, r *http.Request) {
	lines := Read_lines(RcvrLogFile)
	for i := 0; i < len(lines); i++ {
		fmt.Fprintln(w, lines[i])
	}
}

// HTTP(S) Server.
func http_server() {
	http.HandleFunc("/", default_handler)
	http.HandleFunc("/LogViewer/", log_viewer_handler)
	http.HandleFunc("/v1/", Apiv1_handler)
	http.ListenAndServe(":"+strconv.Itoa(Port), nil)
}

// A common but uncommon func main.
func main() {
	initial()
	beacons_db_reader()
	http_server()
}
