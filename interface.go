package logg

import (
	"fmt"
	"io"
	"log"
)

// PRINT

func (l *Logg) Print(args ...interface{}) {
	l.write(0, Empty, []byte(fmt.Sprint(args...)))
}

func (l *Logg) Printf(format string, args ...interface{}) {
	l.write(0, Empty, []byte(fmt.Sprintf(format, args...)))
}

func (l *Logg) Debug(args ...interface{}) {
	l.write(0, Debug, []byte(fmt.Sprint(args...)))
}

func (l *Logg) Debugf(format string, args ...interface{}) {
	l.write(0, Debug, []byte(fmt.Sprintf(format, args...)))
}

func (l *Logg) Info(args ...interface{}) {
	l.write(0, Info, []byte(fmt.Sprint(args...)))
}

func (l *Logg) Infof(format string, args ...interface{}) {
	l.write(0, Info, []byte(fmt.Sprintf(format, args...)))
}

func (l *Logg) Error(args ...interface{}) {
	l.write(0, Error, []byte(fmt.Sprint(args...)))
}

func (l *Logg) Errorf(format string, args ...interface{}) {
	l.write(0, Error, []byte(fmt.Sprintf(format, args...)))
}

func (l *Logg) Warn(args ...interface{}) {
	l.write(0, Warning, []byte(fmt.Sprint(args...)))
}

func (l *Logg) Warnf(format string, args ...interface{}) {
	l.write(0, Warning, []byte(fmt.Sprintf(format, args...)))
}

func (l *Logg) Panic(args ...interface{}) {
	l.write(0, Panic, []byte(fmt.Sprint(args...)))
}

func (l *Logg) Panicf(format string, args ...interface{}) {
	l.write(0, Panic, []byte(fmt.Sprintf(format, args...)))
}

// SETTINGS

func (l *Logg) DebugMode() {
	l.flags = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile
	l.minLevel = Debug
}

func (l *Logg) SetFormat(format format) {
	l.format = format
}

func (l *Logg) SetFlags(flags int) {
	l.flags = flags
}

func (l *Logg) SetWriter(w io.Writer) {
	l.out = w
}

func (l *Logg) ToggleColor(value bool) {
	l.color = value
}

// Global

func Print(args ...interface{}) { logg.Print(args...) }

func Printf(format string, args ...interface{}) { logg.Printf(format, args...) }

func Deb(args ...interface{}) { logg.Debug(args...) }

func Debugf(format string, args ...interface{}) { logg.Debugf(format, args...) }

func DebugMode() { logg.DebugMode() }

func SetFormat(format format) { logg.SetFormat(format) }

func SetFlags(flags int) { logg.SetFlags(flags) }

func SetWriter(w io.Writer) { logg.SetWriter(w) }

func ToggleColor(value bool) { logg.ToggleColor(value) }
