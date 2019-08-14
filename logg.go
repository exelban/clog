// Better log experience in golang.
/*
Usage

	package main
	import (
		_ "github.com/exelban/logg"
		"log"
	)
	func main () {
		log.Print("[ERROR] error text")
	}
*/
package logg

import (
	"fmt"
	"io"
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
	colors ColorsManager
	levels LevelsManager

	format format // pretty or Json
	flags  int    // time flags
	color  bool   // if colors are enabled or not

	out io.Writer
	mu  sync.Mutex
	buf []byte
}

// Logger - is a global object which keep Logg configuration.
var Logger *Logg

func init() {
	c := colors{}
	logg := &Logg{
		colors: ColorsManager{
			list: map[string]string{
				"ERROR": c.Red(),
				"INFO":  c.HiYellow(),
				"WARN":  c.HiGreen(),
				"DEBUG": c.HiCyan(),
			},
		},
		levels: LevelsManager{
			List: []string{
				"DEBUG",
				"INFO",
				"WARN",
				"ERROR",
			},
			Min: "INFO",
		},

		format: Pretty,
		flags:  log.Ldate | log.Ltime,
		color:  true,

		out: os.Stderr,
	}
	logg.flags = log.Flags()

	log.SetOutput(logg)
	log.SetFlags(0)

	Logger = logg
}

// Write writes len(p) bytes from p to the underlying data stream.
// It returns the number of bytes written from p (0 <= n <= len(p))
// and any error encountered that caused the write to stop early.
// Write must return a non-nil error if it returns n < len(p).
func (l *Logg) Write(b []byte) (int, error) {
	m := &message{
		data:  b,
		time:  time.Now(),
		level: "",
		flags: l.flags,
		file:  "",
		line:  0,
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	m.fileLine(4)
	m.level = l.levels.define(b)
	if !l.levels.check(m.level) {
		return len(b), nil
	}

	l.buf = l.buf[:0]
	if l.format == Pretty {
		l.formatHeader(m)
		l.buf = append(l.buf, m.data...)
	} else if l.format == Json {
		b, err := m.MarshalJSON()
		if err != nil {
			return len(b), err
		}
		l.buf = append(l.buf, b...)
	}

	if len(l.buf) == 0 || l.buf[len(l.buf)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}

	return l.out.Write(l.buf)
}

// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

// formatHeader generates time for message in way as log package generate.
func (l *Logg) formatHeader(m *message) {
	t := m.time
	buf := &l.buf

	if l.color {
		m.color = l.colors.define(m.level)
		l.buf = append(l.buf, []byte(fmt.Sprintf("%s[%sm", escape, m.color))...)
	}

	if l.flags&(log.Ldate|log.Ltime|log.Lmicroseconds) != 0 {
		if l.flags&log.LUTC != 0 {
			t = t.UTC()
		}
		if l.flags&log.Ldate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if l.flags&(log.Ltime|log.Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if l.flags&log.Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
		}
	}
	if l.flags&(log.Lshortfile|log.Llongfile) != 0 {
		*buf = append(*buf, m.file...)
		*buf = append(*buf, ':')
		itoa(buf, m.line, -1)
		*buf = append(*buf, ": "...)
	}
}

// SetOutput - sets the output destination for the standard logger.
func SetOutput(writer io.Writer) {
	Logger.mu.Lock()
	Logger.out = writer
	Logger.mu.Unlock()
}

// SetFormat sets the output format (Pretty or Json) for the logger.
func SetFormat(format format) {
	Logger.mu.Lock()
	Logger.format = format
	Logger.mu.Unlock()
}

// SetFlags sets the output flags for the logger.
func SetFlags(flags int) {
	Logger.mu.Lock()
	Logger.flags = flags
	Logger.mu.Unlock()
}

// SetDebug sets the output flags prepared to debug for the logger. And setting minimum level to DEBUG.
func SetDebug() {
	Logger.mu.Lock()
	Logger.levels.Min = "DEBUG"
	Logger.flags = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile
	Logger.mu.Unlock()
}

// SetLevels - sets the levels of logs.
func SetLevels(list []string) {
	Logger.levels.mu.Lock()
	Logger.levels.List = list
	Logger.levels.mu.Unlock()
}

// SetMinLevel - set the minimum levels of logs.
func SetMinLevel(minLevel string) {
	Logger.levels.mu.Lock()
	Logger.levels.Min = minLevel
	Logger.levels.mu.Unlock()
}

// CustomColor - allow to set custom colors for prefix.
func CustomColor(prefix string, v ...interface{}) {
	Logger.colors.CustomColor(prefix, v...)
}
