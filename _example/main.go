package main

import (
	"github.com/pkgz/logg"
	"os"
)

func main() {
	log := logg.New(os.Stdout)
	//log.SetFormat(logg.Json)
	log.SetFlags(logg.LstdFlags)
	//log.DebugMode()

	log.Print("DEBUG test logging, but use a somewhat realistic message length.")
	log.Print("INF test logging, but use a somewhat realistic message length.")
	log.Print("ERROR test logging, but use a somewhat realistic message length.")
	log.Print("WRN test logging, but use a somewhat realistic message length.")
}
