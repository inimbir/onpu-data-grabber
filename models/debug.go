package models

import (
	"fmt"
)

const Env = "report"

func PrintRed(text1 string, text string) {
	if Env == "report" {
		fmt.Println("\033[31m" + text1 + "\033[39m " + text)
		//fmt.Println("\033[31m" + text1 + "\033[39m\n" + text)
	}
}
