package main

import (
	"github.com/exelban/logg"
	"log"
)

func main() {
	logg.SetFormat(logg.Json)
	logg.SetDebug()

	log.Print("[ERROR] test")
	log.Print("[INFO] test")
	log.Print("[DEBUG] test")
	log.Print("[WARN] test")

	log.Print("test")
}
