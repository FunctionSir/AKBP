/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2024-09-16 21:23:59
 * @LastEditTime: 2024-09-22 00:48:36
 * @LastEditors: FunctionSir
 * @Description: -
 * @FilePath: /AKBP/server/auth.go
 */
package main

import (
	"crypto/sha512"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func CalcHash(key string, salt string) string {
	return fmt.Sprintf("%X", sha512.Sum512([]byte(key+salt)))
}

func AuthOK(table string, id string, key string) bool {
	ChkStrWithExit(&table)
	db := DbOpen()
	defer db.Close()
	stmt := DbPrepare(db, "SELECT SALT,HASH FROM \""+table+"\" WHERE ID=?")
	var salt string
	var hash string
	stmt.QueryRow(id).Scan(&salt, &hash)
	if CalcHash(key, salt) == hash {
		return true
	} else {
		return false
	}
}
