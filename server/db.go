/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2024-09-14 21:33:33
 * @LastEditTime: 2024-09-22 01:38:47
 * @LastEditors: FunctionSir
 * @Description: DB related.
 * @FilePath: /AKBP/server/db.go
 */

package main

import (
	"database/sql"
	"strings"
)

// '0'~'9', 'A'~'Z', '_', 'a'~'z' are considered as safe chars.
func chrIsSafe(ch rune) bool {
	return strings.ContainsRune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz", ch)
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

// Init a new DB.
func DbInit() {
	db := DbOpen()
	defer db.Close()
	db.Exec(DB_INIT)
}
