/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2023-07-14 23:11:35
 * @LastEditTime: 2023-07-19 00:48:57
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

var (
	ProgName    string = ""                    //Program name in os.Args[:].
	Port        int    = DEFAULT_PORT          //Server port.
	RcvrLogFile string = DEFAULT_RCVR_LOG_FILE //RcvrLog file.
	API_VER_AVL        = [...]string{"APIv1"}  // API version(s) available.
)

func err_handle(where string, err error) {
	fmt.Println(time.Now().String() + " Error occurred at " + where + ": " + err.Error() + ".")
}

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
				err_handle("main.args_parser", err)
			}
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
	fmt.Println(SPLIT_LINE)
	fmt.Println("Command = " + strings.Join(os.Args[:], " "))
	args_parser()
	fmt.Println(SPLIT_LINE)
	fmt.Println("ProgName = " + ProgName)
	fmt.Println("Port = " + strconv.Itoa(Port))
	fmt.Println("RcvrLogFile = " + RcvrLogFile)
	fmt.Println(SPLIT_LINE)
}

// Feel free to customize it!
func default_handler(w http.ResponseWriter, r *http.Request) {
	_ = r // Don't need var "r" in this ver.
	fmt.Fprintln(w, "This is an [A]nti [K]idnapping [B]eacon [P]roject Server.")
	fmt.Fprintln(w, "Version: "+VER+", Codename: "+CODENAME)
	fmt.Fprintln(w, "!WARNING! !THIS SOFTWARE IS UNDER DEVELOPING AND SHOULD NOT BE USED IN ANY FORMAL SERVERS! !WARNING!") // Remove after the developments.
	fmt.Fprintln(w, "!WARNING! !ANY PACKAGE SENT TO HERE MAY NOT BE HANDLE CORRECTLY! !WARNING!")                           // Remove after the developments.
	fmt.Fprintln(w, "API Version(s) Available: "+strings.Join(API_VER_AVL[:], ", "))
}

// HTTP(S) Server.
func http_server() {
	http.HandleFunc("/", default_handler)
	http.HandleFunc("/v1/", Apiv1_handler)
	http.ListenAndServe(":"+strconv.Itoa(Port), nil)
}

// A common but uncommon func main.
func main() {
	initial()
	http_server()
}
