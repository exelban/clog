package main

import (
	"clog"
	"log"
)

func main () {
	w := clog.Install(clog.Cyan)
	filter := &clog.LevelFilter{
		Levels: []string{"ERROR", "INFO", "WARN", "DEBUG"},
		MinLevel: "WARN",
	}
	w.SetFilters(filter)

	log.Print("[ERROR] error text")
	log.Print("[INFO] info text")
	log.Print("[WARN] warn text")
	log.Print("[DEBUG] debug text")

	log.Print("some text")

	w.Custom("[CUSTOM]", clog.HiBlue, clog.Black, clog.Bold)
	log.Print("[CUSTOM] custom text")

	w.SetMinLevel("INFO")
	log.Print("[INFO] min level text")

	w.Uninstall()

	log.Print("[INFO] uninstall text")
}