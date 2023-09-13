/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2023-08-08 18:02:01
 * @LastEditTime: 2023-09-13 22:17:20
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
	"strings"
	"time"

	"github.com/go-ini/ini"
)

func get_users() []string {
	var r = []string{}
	lines := Read_lines(Userfile)
	for i := 0; i < len(lines); i++ {
		r = append(r, strings.Split(lines[i], ":")[0])
	}
	return r
}

func collect_reg_reqs() []string {
	var r = []string{}
	users := get_users()
	for i := 0; i < len(users); i++ {
		file := Unify_path(UsersDir, true) + Unify_path(users[i], true) + Unify_path(SecDir, true) + Unify_path(FileNameOfRegReqs, false)
		_, err := os.Stat(file)
		if os.IsExist(err) {
			r = append(r, file)
		}
	}
	return r
}

func check_reg_req(f *ini.File) bool {
	if len(f.SectionStrings()) != 0 {
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
	regReqFiles := collect_reg_reqs()
	for i := 0; i < len(regReqFiles); i++ {
		regReq, err := ini.Load(regReqFiles[i])
		Err_handle("main.reg_new_bcns", err)
		if check_reg_req(regReq) {
			regKey := Read_lines(RegKeyFile)
			req := RegAPI + "?regKey=" + regKey[0] + "&bcnKey=" + regReq.Section("").Key("Key").String() + "bcnExtInfo=" + regReq.Section("").Key("ExtInfo").String()
			resp, err := http.Get(req)
			Err_handle("main.reg_new_bcns", err)
			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			Err_handle("main.reg_new_bcns", err)
			bodyStr := string(body)
			if DEBUG {
				fmt.Println(time.Now().String() + " [D] main.reg_new_bcns.resp.Body = " + bodyStr)
			}
			splitedRespBody := strings.Split(bodyStr, "\n")
			if !DEBUG && splitedRespBody[1] == "INFO::SUCCESS" {
				// In dev..
				_ = 0
			} else if DEBUG && splitedRespBody[1] == "WARN::DEBUG_MODE_ON" {
				// In dev...
				_ = 0
			} else {
				// In dev...
				_ = 0
			}
		} else {
			//
			_ = 0
		}
	}
}

func args_parser() {

}

func initial() {
	args_parser()
}

func main() {
	initial()
	// Developing...
	reg_new_bcns()
}
