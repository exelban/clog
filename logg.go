package logg

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

// Base format types
type format int

const (
	Pretty format = iota
	Json
)

// A Logg represents an active logging object that generates lines of
// output to an io.Writer. Each logging operation makes a single call to
// the Logg's Write method.
type Logg struct {
	format    format // output format (string/json)
	flags     int    // time format flags
	color     bool   // colorize output
	callDepth int    // callDepth is the count of the number of frames to skip when computing the file name and line number

	levels *Levels
	colors *Colors

	out io.Writer
	mu  sync.Mutex
}

// Logger - is a global object which keep Logg configuration.
var Logger *Logg

func init() {
	Logger = New(os.Stdout)
}

// Create new logg instance.
func New(w io.Writer) *Logg {
	if w == nil {
		w = ioutil.Discard
	}

	logg := &Logg{
		colors: &Colors{
			list: map[string][]byte{
				"ERROR": generate(Red),
				"INFO":  generate(HiYellow),
				"WARN":  generate(HiGreen),
				"DEBUG": generate(HiCyan),
			},
		},
		levels: &Levels{
			List: []string{
				"DEBUG",
				"INFO",
				"WARN",
				"ERROR",
			},
			Min: "INFO",
		},

		format: Pretty,
		flags:  log.Ldate | log.Ltime, // log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile, // log.Ldate | log.Ltime
		color:  true,

		out: w,
	}
	logg.levels.init()
	logg.colors.init()

	log.SetOutput(logg)
	log.SetFlags(0)

	return logg
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

	l.write(b)

	return
}

func (l *Logg) write(b []byte) {
	m := messagePool.Get().(*message)
	m.data = b
	m.t = time.Now()
	m.flags = l.flags

	l.levels.define(m) // define level
	if !m.allowed {
		return
	}

	if l.color && l.format == Pretty {
		l.colors.define(m) // define color
	}

	m.defineMessage() // define message
	m.defineTime()    // define time
	m.defineCaller()  // define caller

	m.buf = m.buf[:0]
	if l.format == Pretty {
		// add color
		if l.color {
			m.buf = append(m.buf, m.color...)
		}

		// add time
		if len(m.time) != 0 {
			m.buf = append(m.buf, m.time...)
			m.buf = append(m.buf, ' ')
		}

		// add caller
		if len(m.caller) != 0 {
			m.buf = append(m.buf, m.caller...)
			m.buf = append(m.buf, ' ')
		}

		//add level
		if m.brackets {
			m.buf = append(m.buf, '[')
		}
		m.buf = append(m.buf, m.level...)
		if m.brackets {
			m.buf = append(m.buf, ']')
		}

		// add space between level and message
		if len(m.level) == 4 {
			m.buf = append(m.buf, "  "...)
		} else {
			m.buf = append(m.buf, ' ')
		}

		// add message
		m.buf = append(m.buf, m.message...)

		// close color
		if l.color {
			m.buf = append(m.buf, []byte(escapeClose)...)
		}
	} else if l.format == Json {
		b, err := m.MarshalJSON()
		if err != nil {
			return
		}
		m.buf = append(m.buf, b...)
	}

	if len(m.buf) == 0 || m.buf[len(m.buf)-1] != '\n' {
		m.buf = append(m.buf, '\n')
	}

	_, err := l.out.Write(m.buf)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "logg: could not write message: %v\n", err)
	}

	m.reset()

	const maxSize = 1 << 16 // 64KiB
	if cap(m.buf) > maxSize {
		return
	}
	messagePool.Put(m)
}
