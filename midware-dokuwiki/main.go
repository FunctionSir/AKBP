/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2023-08-08 18:02:01
 * @LastEditTime: 2023-10-07 06:56:01
 * @LastEditors: FunctionSir
 * @Description: Midware to connect AKBP-Server and Dokuwiki.
 * @FilePath: /AKBP/midware-dokuwiki/main.go
 */
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-ini/ini"
)

// Global consts and vars are in public.go.
// Some useful funcs are also in public.go.

func LF_process(input string) string {
	return strings.ReplaceAll(input, "\n", "\\n")
}

func get_users() []string {
	var r = []string{}
	lines := Read_lines(BaseDir + Userfile)
	for i := 0; i < len(lines); i++ {
		if !strings.HasPrefix(lines[i], "#") && lines[i] != "" {
			r = append(r, strings.Split(lines[i], ":")[0])
		}
	}
	return r
}

func collect_reg_reqs() ([]string, []string) {
	var r0 = []string{}
	var r1 = []string{}
	users := get_users()
	for i := 0; i < len(users); i++ {
		file := Unify_path(BaseDir, true) + Unify_path(UsersDir, true) + Unify_path(users[i], true) + Unify_path(SecDir, true) + Unify_path(ThirdDir, true) + Unify_path(FileNameOfRegReqs, false)
		if File_is_exist(file) {
			r0 = append(r0, file)
			r1 = append(r1, users[i])
		}
	}
	return r0, r1
}

func check_reg_req(f *ini.File) bool {
	if len(f.SectionStrings()) != 1 {
		return false
	}
	if len(f.Section("").KeyStrings()) != 2 {
		return false
	}
	if Find_str(f.Section("").KeyStrings(), "Key") == -1 {
		return false
	}
	return true
}

func reg_new_bcns() {
	regReqFiles, relatedUser := collect_reg_reqs()
	for i := 0; i < len(regReqFiles); i++ {
		regReq, err := ini.Load(regReqFiles[i])
		Err_handle("main.reg_new_bcns", err)
		if check_reg_req(regReq) {
			regKey := Read_lines(BaseDir + RegKeyFile)
			if len(regKey) == 0 {
				fmt.Println(time.Now().String() + " [W] Seems reg key will be empty.")
				regKey = append(regKey, "")
			}
			req := RegAPI + "?regKey=" + regKey[0] + "&bcnKey=" + regReq.Section("").Key("Key").String() + "&bcnExtInfo=" + regReq.Section("").Key("ExtInfo").String()
			resp, err := http.Get(req)
			var body []byte
			if !Err_handle("main.reg_new_bcns", err) {
				body, err = io.ReadAll(resp.Body)
				resp.Body.Close()
				Err_handle("main.reg_new_bcns", err)
			}
			bodyStr := string(body)
			if DEBUG {
				fmt.Println(time.Now().String() + " [D] main.reg_new_bcns.resp.Body = " + LF_process(bodyStr))
			}
			splitedRespBody := strings.Split(bodyStr, "\n")
			bcnListFile := "/dev/null" // For Linux, it may will pervert some bad things. For windows, it also can do that by causing a error.
			wFlag := false
			if len(splitedRespBody) >= 2 {
				if !DEBUG && Find_str(splitedRespBody, "INFO::SUCCESS") != -1 {
					bcnListFile = Unify_path(BaseDir, true) + Unify_path(UsersDir, true) + Unify_path(relatedUser[i], true) + Unify_path(SecDir, true) + Unify_path(ThirdDir, true) + Unify_path(BcnListFile, false)
					wFlag = true
				} else if DEBUG && Find_str(splitedRespBody, "WARN::DEBUG_MODE_ON") != -1 {
					bcnListFile = Unify_path(BaseDir, true) + Unify_path(UsersDir, true) + Unify_path(relatedUser[i], true) + Unify_path(SecDir, true) + Unify_path(ThirdDir, true) + Unify_path("(DEBUG)"+BcnListFile, false)
					wFlag = true
				} else {
					if DEBUG {
						fmt.Println(time.Now().String() + " [D] (!DEBUG && Find_str(splitedRespBody, \"INFO::SUCCESS\") != -1) || (DEBUG && Find_str(splitedRespBody, \"WARN::DEBUG_MODE_ON\") != -1) == false.")
					}
					wFlag = false
				}
				if wFlag {
					var toWrite = []string{}
					var bcnUUID = ""
					if !DEBUG {
						bcnUUID = splitedRespBody[Find_str(splitedRespBody, "INFO::SUCCESS")-1]
						toWrite = []string{"[" + bcnUUID + "]", "ExtInfo = " + regReq.Section("").Key("ExtInfo").String()}

					} else {
						bcnUUID = splitedRespBody[Find_str(splitedRespBody, "WARN::DEBUG_MODE_ON")-1]
						toWrite = []string{"[" + bcnUUID + "]", "ExtInfo = " + regReq.Section("").Key("ExtInfo").String()}
					}
					if File_is_exist(bcnListFile) {
						Append_lines(bcnListFile, toWrite)
					} else {
						Create_file(bcnListFile, toWrite)
					}
					gen_rels_db_entry(bcnUUID, relatedUser[i])
					err := os.Remove(regReqFiles[i]) // Remove reg req if succeeded.
					Err_handle("main.reg_new_bcns", err)
				}
			} else if DEBUG {
				fmt.Println(time.Now().String() + " [D] Got an invalid resp body, len(splitedRespBody) >= 1 is not true.")
			}
		} else {
			if DEBUG {
				fmt.Println(time.Now().String() + " [D] Got an invalid reg req.")
			}
			lines := Read_lines(regReqFiles[i])
			toWrite := []string{"# INVALID REG REQ, PLEASE CHECK IT AGAIN! #"}
			if len(lines) == 0 {
				Append_lines(regReqFiles[i], toWrite)
			} else if len(lines) >= 1 && lines[len(lines)-1] != "# INVALID REG REQ, PLEASE CHECK IT AGAIN! #" {
				Append_lines(regReqFiles[i], toWrite)
			}
		}
	}
}

func gen_rels_db_entry(uuid, username string) {
	f := Unify_path(BaseDir, true) + Unify_path(RelsDB, false)
	if File_is_exist(f) {
		Append_lines(f, []string{uuid + " " + username})
	} else {
		Create_file(f, []string{uuid + " " + username})
	}
}

// func rcvr_log_parser() {
// 	file, err := ini.Load(Unify_path(BaseDir, true) + Unify_path(RcvrLogFile, false))
// 	if !Err_handle("main.rcvr_log_parser", err) {
// 		sectionStrs := file.SectionStrings()
// 		for i := 1; i < len(sectionStrs); i++ {
// 			_ = i
// 		}
// 	}
// }

func args_parser() {
	ProgName = os.Args[0]
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-b", "--base-dir":
			BaseDir = os.Args[i+1]
		case "-ud", "--users-dir":
			UsersDir = os.Args[i+1]
		case "-sd", "--sec-dir":
			SecDir = os.Args[i+1]
		case "-td", "--third-dir":
			ThirdDir = os.Args[i+1]
		case "-blf", "--bcn-list-file":
			BcnListFile = os.Args[i+1]
		case "-rkf", "--reg-key-file":
			RegKeyFile = os.Args[i+1]
		case "-rrfn", "--reg-req-files-name":
			FileNameOfRegReqs = os.Args[i+1]
		case "-uf", "--userfile":
			Userfile = os.Args[i+1]
		case "-rdb", "--rels-db":
			RelsDB = os.Args[i+1]
		case "-rapi", "--reg-api":
			RegAPI = os.Args[i+1]
		case "-ag", "--act-gap":
			var tmp int
			tmp, err := strconv.Atoi(os.Args[i+1])
			if err == nil {
				if tmp >= 0 {
					ActGap = tmp
				}
			} else {
				Err_handle("main.args_parser", err)
			}
		}
	}
}

func initial() {
	fmt.Println("AKBP Midware For Dokuwiki")
	fmt.Println("Version: " + VER + " (" + CODENAME + ")")
	fmt.Println("This is a FOSS under the AGPLv3.")
	args_parser()
}

func main() {
	initial()
	flag := false
	for !flag {
		// rcvr_log_parser()
		reg_new_bcns()
		time.Sleep(time.Duration(ActGap) * time.Second)
	}
}
