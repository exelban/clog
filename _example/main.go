package main

import (
	"github.com/exelban/logg"
	"log"
)

func main() {
	//logg.SetFormat(logg.Pretty)
	logg.SetDebug()

	log.Print("[ERROR] test")
	log.Print("[INFO] test UNMARSHAL_ERROR")
	log.Print("[DEBUG] test UNMARSHAL_ERROR")
	log.Print("[WARN] test")

	log.Print("test")
}
