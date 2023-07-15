/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2023-07-14 23:11:35
 * @LastEditTime: 2023-07-16 01:26:50
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

const (
	VER                   string = "0.1-alpha"        // Version
	CODE_NAME             string = "Capsule"          // Code name of this version
	SPLIT_LINE            string = "----------------" // Set the split line you want to use here.
	DEFAULT_PORT          int    = 3690               // Default port (=(ord("A")+ord("K")+ord("B")+ord("S"))*10).
	DEFAULT_RCVR_LOG_FILE string = "akbo-rcvr.log"    // Default [A]nti [K]idnapping [B]eacon [O]rganization [-] [R]e[C]ei[V]e[R]('s/s').[Log] file.
)

var (
	ProgName    string = ""                    //Program name in os.Args[:].
	Port        int    = DEFAULT_PORT          //Server port.
	RcvrLogFile string = DEFAULT_RCVR_LOG_FILE //RcvrLog file.
)

func err_handle(where string, err error) {
	fmt.Println(time.Now().String() + " Error occurred at " + where + ": " + err.Error() + ".")
}

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
		}
	}
}

// Do this first.
func initial() {
	fmt.Println("[A]nti [K]idnapping [B]eacon [P]roject Server")
	fmt.Println("Version: " + VER + ", CodeName: " + CODE_NAME)
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

// HTTP Server
func http_server() {
	//http.HandleFunc("/", nil) //Developing
	http.HandleFunc("/v1/", Apiv1_handler)
	http.ListenAndServe(":"+strconv.Itoa(Port), nil)
}

// A common but uncommon func main.
func main() {
	initial()
	http_server()
}
