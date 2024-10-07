/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2024-09-24 23:00:35
 * @LastEditTime: 2024-10-05 23:07:59
 * @LastEditors: FunctionSir
 * @Description: -
 * @FilePath: /AKBP/midware-dokuwiki/fileio.go
 */
package main

import "os"

// If file exists and is not dir -> true.
// If file not exists or it's a dir -> false.
func FileExists(path string) bool {
	stat, err := os.Stat(path)
	if os.IsNotExist(err) || stat.IsDir() {
		return false
	}
	return true
}

func DirExists(path string) bool {
	stat, err := os.Stat(path)
	if os.IsNotExist(err) || !stat.IsDir() {
		return false
	}
	return true
}
