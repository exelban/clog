package main

import (
	"log"
	"logg"
)

func main() {
	logg.SetFormat(logg.Pretty)

	log.Print("[ERROR] test")
	log.Print("[INFO] test")
	log.Print("WARN test")

	log.Print("test")
}
