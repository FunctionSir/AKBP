/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2024-09-14 21:18:18
 * @LastEditTime: 2024-09-22 01:15:13
 * @LastEditors: FunctionSir
 * @Description: Beacons related.
 * @FilePath: /AKBP/server/msg.go
 */

package main

import (
	"crypto/sha512"
	"fmt"
)

func CalcMsgHash(bid string, eid string, ts int, msg string) string {
	info := fmt.Sprintf("%s\t%s\t%d\t%s", bid, eid, ts, msg)
	return fmt.Sprintf("%X", sha512.Sum512([]byte(info)))
}

func MsgExists(bid string, eid string, ts int, msg string) bool {
	db := DbOpen()
	defer db.Close()
	stmt := DbPrepare(db, "SELECT HASH FROM RECEIVED WHERE HASH=?")
	var tmp string
	err := stmt.QueryRow(CalcMsgHash(bid, eid, ts, msg)).Scan(&tmp)
	return err == nil
}

func AddRecord(bid string, eid string, ts int, msg string, origin string) bool {
	// If msg exists, do nothing.
	if MsgExists(bid, eid, ts, msg) {
		return false
	}
	db := DbOpen()
	defer db.Close()
	stmt := DbPrepare(db, "INSERT INTO RECORDS VALUES(?,?,?,?,?)")
	stmt.Exec(bid, eid, ts, msg, origin)
	stmt = DbPrepare(db, "INSERT INTO RECEIVED VALUES(?)")
	stmt.Exec(CalcMsgHash(bid, eid, ts, msg))
	return true
}
