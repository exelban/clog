package main

import (
	"github.com/exelban/clog"
	"log"
)

func main () {
	w := clog.Install(clog.Cyan)

	w.Custom("[ERROR]", clog.Red)
	w.Prefix("[INFO]", clog.Colors.HiYellow)
	w.Prefix("[WARN]", clog.Colors.Green)
	w.Prefix("[DEBUG]", clog.Colors.Blue)

	log.Print("[ERROR] error text")
	log.Print("[INFO] info text")
	log.Print("[WARN] warn text")
	log.Print("[DEBUG] debug text")

	log.Print("some text")

	w.Uninstall()

	log.Print("[INFO] uninstall text")
}