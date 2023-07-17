/*
 * @Author: FunctionSir
 * @Date: 2023-07-17 22:47:42
 * @LastEditTime: 2023-07-17 22:49:23
 * @LastEditors: FunctionSir
 * @Description: Public functions of AKBP Server.
 * @FilePath: /AKBP/server/pub_funcs.go
 */
package main

import "strconv"

func Is_pure_num_str(str string) bool {
	_, err := strconv.Atoi(str)
	if err == nil {
		return true
	} else {
		return false
	}
}
