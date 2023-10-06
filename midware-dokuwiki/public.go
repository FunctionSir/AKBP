/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2023-08-11 18:47:23
 * @LastEditTime: 2023-10-07 07:22:28
 * @LastEditors: FunctionSir
 * @Description: Public consts, vars, and functions of midware-dokuwiki.
 * @FilePath: /AKBP/midware-dokuwiki/public.go
 */
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	DEBUG                         bool   = false                          // Debug flag, for debugging or developing purposes ONLY!
	DEBUG_VERBOSE                 bool   = false                          // Be verbose.
	VER                           string = "0.1-alpha"                    // Version.
	CODENAME                      string = "Capsule"                      // Code name of this version.
	SPLIT_LINE                    string = "----------------"             // Set the split line you want to use here.
	DEFAULT_BASE_DIR              string = "/home/funcsir/Projects/TEST/" // Default value of var BaseDir.
	DEFAULT_USERS_DIR             string = "dw-users"                     // Default value of var "UsersDir".
	DEFAULT_SEC_DIR               string = "private"                      // Default value of var "SecDir".
	DEFAULT_THIRD_DIR             string = "AKBP"                         // Default value of var "ThirdDir"
	DEFAULT_BCN_LIST_FILE         string = "MyBeacons"                    // Default value of var "BcnListFile"
	DEFAULT_REG_KEY_FILE          string = "regkey.txt"                   // Default reg key file.
	DEFAULT_FILE_NAME_OF_REG_REQS string = "akbp-new-bcn-reg.txt"         // Default value of RegReqFile.
	DEFAULT_USERFILE              string = "users.auth.php"               // Default value of Userfile.
	DEFAULT_RELS_DB               string = "relations.db"                 // Default value of RelsDB.
	DEFAULT_RCVR_LOG_FILE         string = "akbo-rcvr.log"                // Default value of RcvrLogFile.
	DEFAULT_REG_API               string = "http://127.0.0.1:3690/reg/"   // Default reg API.
	DEFAULT_ACT_GAP               int    = 1                              // Default action gap (seconds).
)

var (
	ProgName          string = ""                            // Program name in os.Args[:].
	BaseDir           string = DEFAULT_BASE_DIR              // Base dir if you want to use this without full path of required dirs.
	UsersDir          string = DEFAULT_USERS_DIR             // Dir of users' spaces
	SecDir            string = DEFAULT_SEC_DIR               // Secondary dir.
	ThirdDir          string = DEFAULT_THIRD_DIR             // Third dir.
	BcnListFile       string = DEFAULT_BCN_LIST_FILE         // Beacons list in dokuwiki.
	RegKeyFile        string = DEFAULT_REG_KEY_FILE          // Reg key file.
	FileNameOfRegReqs string = DEFAULT_FILE_NAME_OF_REG_REQS // File name of reg reqs.
	Userfile          string = DEFAULT_USERFILE              // Typically, it is users.auth.php.
	RelsDB            string = DEFAULT_RELS_DB               // Relations DB.
	RcvrLogFile       string = DEFAULT_RCVR_LOG_FILE         // [A]nti [K]idnapping [B]eacon [O]rganization [-] [R]e[C]ei[V]e[R]('s/s').[Log] file.
	RegAPI            string = DEFAULT_REG_API               // Reg API.
	ActGap            int    = DEFAULT_ACT_GAP               // Gap between actions. Should >= 0.
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
		if DEBUG && DEBUG_VERBOSE {
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
	if p == "" {
		return p
	}
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

func File_is_exist(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	} else if os.IsExist(err) {
		return true
	} else {
		return false
	}
}

func Append_lines(name string, lines []string) []error {
	var errs = []error{}
	s, e := os.Stat(name)
	if Err_handle("main.Append_lines", e) {
		errs = append(errs, e)
	}
	f, e := os.OpenFile(name, os.O_RDWR|os.O_APPEND, s.Mode())
	if Err_handle("main.Append_lines", e) {
		errs = append(errs, e)
	}
	defer func() {
		e := f.Close()
		c := 0
		for Err_handle("main.Append_lines", e) && c <= 8 {
			e = f.Close()
			c++
		}
	}()
	tmp := []byte{0}
	if s.Size() > 0 {
		_, e = f.Seek(-1, io.SeekEnd)
		if Err_handle("main.Append_lines", e) {
			errs = append(errs, e)
		}
		_, _ = f.Read(tmp)
	}
	if s.Size() != 0 && string(tmp) != "\n" {
		f.Seek(0, io.SeekEnd)
		if Err_handle("main.Append_lines", e) {
			errs = append(errs, e)
		}
		_, e := f.WriteString("\n")
		if Err_handle("main.Append_lines", e) {
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
		if Err_handle("main.Append_lines", e) {
			errs = append(errs, e)
		}
	}
	return errs
}

func Create_file(name string, lines []string) error {
	var s string = ""
	f, e := os.Create(name)
	Err_handle("main.Create_file", e)
	for i := 0; i < len(lines); i++ {
		if strings.HasSuffix(lines[i], "\n") {
			s = s + lines[i]
		} else {
			s = s + lines[i] + "\n"
		}
	}
	_, e = f.WriteString(s)
	Err_handle("main.Create_file", e)
	return e
}
