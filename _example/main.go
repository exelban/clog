package main

import (
	"github.com/pkgz/logg"
	"log"
	"os"
)

func main() {
	logg.NewGlobal(os.Stdout)
	//log := logg.New(os.Stdout)
	//logg.SetFormat(logg.Json)
	//logg.NewGlobal(os.Stdout)
	//logg.DebugMode()

	log.Print("[INF] test")
	log.Print("test")

	//log.Print("[INFO] test UNMARSHAL_ERROR")
	//log.Print("[DEBUG] test UNMARSHAL_ERROR")
	//log.Print("[WARN] test")
	//
	//log.Print("test")
	//
	//log.Print("ERROR test")
	//log.Print("INFO test UNMARSHAL_ERROR")
	//log.Print("DEBUG test UNMARSHAL_ERROR")
	//log.Print("WARN test")
}
