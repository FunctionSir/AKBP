/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2024-09-14 21:33:33
 * @LastEditTime: 2024-10-06 19:42:51
 * @LastEditors: FunctionSir
 * @Description: DB related.
 * @FilePath: /AKBP/midware-dokuwiki/db.go
 */

package main

import (
	"crypto/sha512"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// '0'~'9', 'A'~'Z', '_', 'a'~'z' are considered as safe chars.
func chrIsSafe(ch rune) bool {
	return strings.ContainsRune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz", ch)
}

func CalcHash(key string, salt string) string {
	return fmt.Sprintf("%X", sha512.Sum512([]byte(key+salt)))
}

// A kind of protection before SQL actions.
func ChkStrWithExit(str *string) {
	for _, x := range *str {
		if !chrIsSafe(x) {
			LogFatalln("Unsafe char found in string \"" + *str + "\" by strict str checker for DB.")
		}
	}
}

// Chk str, but do not exit.
func ChkStrNoExit(str *string) bool {
	for _, x := range *str {
		if !chrIsSafe(x) {
			return false
		}
	}
	return true
}

// Open DB using DB file specified in global var.
func DbOpen() *sql.DB {
	db, err := sql.Open("sqlite3", Db)
	if err != nil {
		LogFatalln("Error occurred when opening the database: " + strings.Trim(err.Error(), "\n"))
	}
	return db
}

// Prepare a query.
func DbPrepare(db *sql.DB, query string) *sql.Stmt {
	stmt, err := db.Prepare(query)
	if err != nil {
		LogFatalln("Error occurred when preparing the SQL statement: " + strings.Trim(err.Error(), "\n"))
	}
	return stmt
}

func QueryRecs() *sql.Rows {
	db := DbOpen()
	defer db.Close()
	stmt := DbPrepare(db, "SELECT ROWID,BID,EID,TS,MSG,ORIGIN FROM RECORDS ORDER BY TS;")
	rows, err := stmt.Query()
	if err != nil {
		LogWarnln("An error occurred when performing a query.")
		return nil
	}
	return rows
}

func IsBidExists(bid string) bool {
	db := DbOpen()
	defer db.Close()
	stmt := DbPrepare(db, "SELECT ROWID FROM BEACONS WHERE ID=?")
	var tmp int
	err := stmt.QueryRow(bid).Scan(&tmp)
	return err == nil
}

func IsSaltExists(salt string) bool {
	db := DbOpen()
	defer db.Close()
	stmt := DbPrepare(db, "SELECT * FROM BEACONS WHERE SALT=?")
	err := stmt.QueryRow(salt)
	return err == nil
}

func RegBeacon(bid string, key string, note string) bool {
	db := DbOpen()
	defer db.Close()
	stmt := DbPrepare(db, "INSERT INTO BEACONS VALUES(?,?,?,?)")
	salt := GenSalt()
	hash := CalcHash(key, salt)
	_, err := stmt.Exec(bid, salt, hash, note)
	return err == nil
}

func GenSalt() string {
	salt := uuid.New().String()
	for IsSaltExists(salt) {
		salt = uuid.New().String()
	}
	return salt
}
