package main

import (
	"clog"
	"log"
)

func main () {
	w := clog.Install(clog.Cyan)

	log.Print("[ERROR] error text")
	log.Print("[INFO] info text")
	log.Print("[WARN] warn text")
	log.Print("[DEBUG] debug text")

	log.Print("some text")

	w.Custom("[CUSTOM]", clog.HiBlue, clog.Black, clog.Bold)
	log.Print("[CUSTOM] custom text")

	w.Uninstall()

	log.Print("[INFO] uninstall text")
}