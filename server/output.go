/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2024-09-12 22:47:30
 * @LastEditTime: 2024-09-22 11:52:47
 * @LastEditors: FunctionSir
 * @Description: Print logs, or other output things.
 * @FilePath: /AKBP/server/output.go
 */

package main

import (
	"log"

	"github.com/fatih/color"
)

func LogFatalln(s string) {
	c := color.New(color.FgHiRed, color.Underline)
	log.Fatalln(c.Sprint(s))
}

func LogWarnln(s string) {
	c := color.New(color.FgHiYellow)
	log.Println(c.Sprint(s))
}

func LogInfoln(s string) {
	c := color.New(color.FgHiGreen)
	log.Println(c.Sprint(s))
}

func Hello() {
	c := color.New(color.FgHiBlue)
	c.Println("[A]nti [K]idnapping [B]eacon [P]roject Server")
	c.Println("Version: " + VER + " (" + CODENAME + ")")
}
