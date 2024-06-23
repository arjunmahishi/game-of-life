package main

import "fmt"

type colorFunc func(string) string

func makeColor(code string) colorFunc {
	return func(s string) string {
		return fmt.Sprintf("\033[%sm%s\033[0m", code, s)
	}
}

var (
	green   = makeColor("32")
	yellow  = makeColor("33")
	noColor = makeColor("0")
)
