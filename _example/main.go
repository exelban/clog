package main

import (
	"log"
	"logg"
)

func main() {
	logg.SetFormat(logg.Pretty)
	logg.SetDebug()

	log.Print("[ERROR] test")
	log.Print("[INFO] test")
	log.Print("[DEBUG] test")
	log.Print("[WARN] test")

	log.Print("test")
}
