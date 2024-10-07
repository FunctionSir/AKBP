/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2024-09-14 21:33:33
 * @LastEditTime: 2024-10-07 17:50:21
 * @LastEditors: FunctionSir
 * @Description: DB related.
 * @FilePath: /AKBP/server/db.go
 */

package main

import (
	"database/sql"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// '0'~'9', 'A'~'Z', '_', 'a'~'z' are considered as safe chars.
func chrIsSafe(ch rune, extra string) bool {
	return strings.ContainsRune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz"+extra, ch)
}

// A kind of protection before SQL actions.
func ChkStrWithExit(str *string) {
	for _, x := range *str {
		if !chrIsSafe(x, "") {
			LogFatalln("Unsafe char found in string \"" + *str + "\" by strict str checker for DB.")
		}
	}
}

// Chk str, but do not exit.
func ChkStrNoExit(str *string, extra string) bool {
	for _, x := range *str {
		if !chrIsSafe(x, extra) {
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
	stmt := DbPrepare(db, "SELECT BID,EID,TS,MSG,BANNED FROM RECORDS;")
	rows, err := stmt.Query()
	if err != nil {
		LogWarnln("An error occurred when performing a query.")
		return nil
	}
	return rows
}

// Init a new DB.
func DbInit() {
	db := DbOpen()
	defer db.Close()
	db.Exec(DB_INIT)
}
