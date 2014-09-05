package main

import (
	"tpl"

	"github.com/hoisie/web"
)

func init_web() {
	web.Get("/index", tpl.Index)
}

func hello(val string) string {
	return "hello " + val
}

func main() {
	init_web()
	web.Get("/(.*)", hello)
	web.Run("0.0.0.0:9999")
}
