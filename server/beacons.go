/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2024-09-14 21:18:18
 * @LastEditTime: 2024-09-16 22:42:53
 * @LastEditors: FunctionSir
 * @Description: Beacons related.
 * @FilePath: /AKBP/server/beacons.go
 */

package main

import (
	_ "github.com/mattn/go-sqlite3"
)

func AddRecord(bid string, eid string, ts int, msg string, origin string) {
	db := DbOpen()
	defer db.Close()
	stmt := DbPrepare(db, "INSERT INTO RECORDS VALUES(?,?,?,?,?)")
	stmt.Exec(bid, eid, ts, msg, origin)
}
