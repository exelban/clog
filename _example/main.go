package main

import (
	"log"
	"logg"
)

func main() {
	w := logg.Install()
	filter := &logg.LevelFilter{
		Levels:   []string{"ERROR", "INFO", "WARN", "DEBUG"},
		MinLevel: "WARN",
	}
	w.SetFilters(filter)

	log.Print("[ERROR] error text")
	log.Print("[INFO] info text")
	log.Print("[WARN] warn text")
	log.Print("[DEBUG] debug text")

	log.Print("some text")

	w.Custom("[CUSTOM]", logg.HiBlue, logg.Black, logg.Bold)
	log.Print("[CUSTOM] custom text")

	w.SetMinLevel("INFO")
	log.Print("[INFO] min level text")

	w.Uninstall()

	log.Print("[INFO] uninstall text")
}
