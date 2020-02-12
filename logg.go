// Better log experience in golang.
/*
Usage
	package main

	import (
		"github.com/pkgz/logg"
		"log"
		"os"
	)

	func main () {
		logg.NewGlobal(os.Stdout)

		log.Print("DEBUG some text")
		log.Print("INF some text")
		log.Print("[ERROR] some text")
		log.Print("[WARN] some text")
		log.Print("some text")
	}
*/

package logg

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// A Logg represents an active logging object that generates lines of
// output to an io.Writer. Each logging operation makes a single call to
// the Logg's Write method.
type Logg struct {
	format format // output format (string/json)
	flags  int    // time format flags
	color  bool   // colorize output

	minLevel level
	out      io.Writer
}

// Create new a new logg.
func New(w io.Writer) *Logg {
	if w == nil {
		w = ioutil.Discard
	}

	return &Logg{
		out: w,

		format:   DefaultFormat,
		flags:    DefaultFlags,
		color:    DefaultColorOutput,
		minLevel: DefaultMinimumLevel,
	}
}

// Allows to create a new global logger.
func NewGlobal(w io.Writer) {
	logg = New(w)

	log.SetOutput(logg)
	log.SetFlags(0)
}

// Write writes len(p) bytes from p to the underlying data stream.
// It returns the number of bytes written from p (0 <= n <= len(p))
// and any error encountered that caused the write to stop early.
// Write must return a non-nil error if it returns n < len(p).
func (l *Logg) Write(b []byte) (n int, err error) {
	n = len(b)
	if n > 0 && b[n-1] == '\n' {
		b = b[0 : n-1]
	}

	l.write(3, Empty, b)
	return
}

func (l *Logg) write(calldepth int, level level, b []byte) {
	if b == nil {
		return
	}

	if level == Empty {
		level = defineLevel(&b)
		b = removeLevel(b, level)
	}

	if level != Empty && level < l.minLevel {
		return
	}

	m := newMessage(level, ContextCallDepth+calldepth, l.flags, l.format, l.color)

	out := l.out
	if out == os.Stdout && (level > Error) {
		out = os.Stderr
	}

	if err := write(l.out, m.build(b)); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "logg: could not write message: %v\n", err)
	}
	m.put()
}

// Writer returns the output destination for the standard logger.
func (l *Logg) Writer() io.Writer {
	return l.out
}
