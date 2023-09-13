/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2023-08-11 18:47:23
 * @LastEditTime: 2023-09-03 17:22:02
 * @LastEditors: FunctionSir
 * @Description: Public consts, vars, and functions of midware-dokuwiki.
 * @FilePath: /AKBP/midware-dokuwiki/public.go
 */
package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	DEBUG                         bool   = true               // Debug flag, for debugging or developing purposes ONLY!
	VER                           string = "0.1-alpha"        // Version.
	CODENAME                      string = "Capsule"          // Code name of this version.
	SPLIT_LINE                    string = "----------------" // Set the split line you want to use here.
	DEFAULT_USERS_DIR             string = "dw-users"         // Default value of var "UsersDir".
	DEFAULT_SEC_DIR               string = "private"          // Default value of var "SecDir".
	DEFAULT_THIRD_DIR             string = "AKBP"
	DEFAULT_REG_KEY_FILE          string = "regkey.txt"                 // Default reg key file.
	DEFAULT_FILE_NAME_OF_REG_REQS string = "akbp-new-bcn-reg.txt"       // Default value of RegReqFile.
	DEFAULT_USERFILE              string = "users.auth.php"             // Default value of Userfile.
	DEFAULT_RELS_DB               string = "relations.db"               // Default value of RelsDB.
	DEFAULT_REG_API               string = "http://127.0.0.1:3690/reg/" // Default reg API.
)

var (
	UsersDir          string = DEFAULT_USERS_DIR             // Dir of users' spaces
	SecDir            string = DEFAULT_SEC_DIR               // Secondary dir.
	RegKeyFile        string = DEFAULT_REG_KEY_FILE          // Reg key file.
	FileNameOfRegReqs string = DEFAULT_FILE_NAME_OF_REG_REQS // File name of reg reqs.
	Userfile          string = DEFAULT_USERFILE              // Typically, it is users.auth.php.
	RelsDB            string = DEFAULT_RELS_DB               // Relations DB.
	RegAPI            string = DEFAULT_REG_API               // Reg API.
)

func Is_windows() bool {
	if runtime.GOOS == "windows" {
		return true
	} else {
		return false
	}
}

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

func Find_str(source []string, target string) int {
	for i := 0; i < len(source); i++ {
		if source[i] == target {
			return i
		}
	}
	return -1
}

func Unify_path(p string, f bool) string {
	check_EOP := func(path string, flag bool) bool {
		if ((!strings.HasSuffix(path, "/")) && flag) || (strings.HasSuffix(path, "/") && (!flag)) {
			return false
		}
		return true
	}
	if Is_windows() {
		p = strings.ReplaceAll(p, "\\", "/")
	}
	if check_EOP(p, false) && f {
		p = p + "/"
	}
	for check_EOP(p, true) && Is_windows() && (!f) {
		p = p[:len(p)-2]
	}
	return p
}
