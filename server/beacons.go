/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2024-09-14 21:18:18
 * @LastEditTime: 2024-09-14 23:03:05
 * @LastEditors: FunctionSir
 * @Description: Beacons related.
 * @FilePath: /AKBP/server/beacons.go
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

func BeaconAuthOk(id string, key string) bool {
	db := DbOpen()
	defer db.Close()
	stmt := DbPrepare(db, "SELECT SALT,HASH FROM BEACONS WHERE ID=?")
	var salt string
	var hash string
	stmt.QueryRow(id).Scan(&salt, &hash)
	if CalcHash(key, salt) == hash {
		return true
	} else {
		return false
	}
}

func AddRecord(bid string, eid string, ts int, msg string) {
	db := DbOpen()
	defer db.Close()
	stmt := DbPrepare(db, "INSERT INTO RECORDS VALUES(?,?,?,?)")
	stmt.Exec(bid, eid, ts, msg)
}
