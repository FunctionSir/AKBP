/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2024-09-14 21:33:33
 * @LastEditTime: 2024-09-14 22:14:15
 * @LastEditors: FunctionSir
 * @Description: DB related.
 * @FilePath: /AKBP/server/db.go
 */

package main

import (
	"database/sql"
	"strings"
)

// Open DB using DB file specified in global var.
func DbOpen() *sql.DB {
	db, err := sql.Open("sqlite3", Db)
	if err != nil {
		LogFatalln("Error occurred when opening the database: " + strings.Trim(err.Error(), "\n"))
	}
	return db
}

func DbPrepare(db *sql.DB, query string) *sql.Stmt {
	stmt, err := db.Prepare(query)
	if err != nil {
		LogFatalln("Error occurred when preparing the SQL statement: " + strings.Trim(err.Error(), "\n"))
	}
	return stmt
}
