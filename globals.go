package logg

import (
	"log"
	"os"
)

// Global Logg configuration
var logg = New(os.Stderr)

// Default parameters
const (
	DefaultFormat       = Pretty
	DefaultColorOutput  = true
	DefaultFlags        = log.LstdFlags // log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile
	DefaultMinimumLevel = Debug

	ContextCallDepth = 4

	escape      = "\x1b"
	escapeClose = "\x1b[0m"

	digits01 = "0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789"
	digits10 = "0000000000111111111122222222223333333333444444444455555555556666666666777777777788888888889999999999"
)

// Base types
type (
	format int
	level  int
)

var (
	levels = []string{"DBG", "INF", "ERR", "WRN", "PNC", "DEBUG", "INFO", "ERROR", "WARN", "PANIC"}
	colors = [][]byte{generate(HiCyan), generate(HiYellow), generate(Red), generate(HiGreen), generate(Red)}
)

// Output formats
const (
	Pretty format = iota
	Json
)

// Log levels.
const (
	Debug level = iota
	Info
	Error
	Warning
	Panic

	Empty level = -1
)

// Base colors for console
const (
	Black int = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// High intensity colors
const (
	HiBlack int = iota + 60
	HiRed
	HiGreen
	HiYellow
	HiBlue
	HiMagenta
	HiCyan
	HiWhite
)
