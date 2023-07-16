/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2023-07-16 00:26:53
 * @LastEditTime: 2023-07-17 01:06:01
 * @LastEditors: FunctionSir
 * @Description: APIv1 related funcs.
 * @FilePath: /AKBP/server/apiv1.go
 */
package main

import (
	"fmt"
	"net/http"
	"strings"
)

func Apiv1_handler(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		fmt.Fprintln(w, "ERR::UUID_NOT_FOUND")
		return
	}
	if len(uuid) != 36 || strings.Count(uuid, "-") != 4 { //Just a simple check
		fmt.Fprintln(w, "ERR::ILLEGAL_UUID")
		return
	}
}
