package logg

import (
	"fmt"
	"github.com/francoispqt/gojay"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

// Base format types
const (
	Pretty int = iota
	Json
)

type Logg struct {
	colors ColorsManager
	levels LevelsManager

	format int
	flags  int
	color  bool

	out io.Writer
	mu  sync.Mutex
}

var Logger *Logg

func init() {
	logg := &Logg{
		colors: ColorsManager{
			list: make(map[string]string),
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
		color:  false,

		out: os.Stderr,
	}
	logg.flags = log.Flags()

	log.SetOutput(logg)
	log.SetFlags(0)

	if logg.color {
		logg.colors.CustomColor("[ERROR]", Red)
		logg.colors.CustomColor("[INFO]", HiYellow)
		logg.colors.CustomColor("[WARN]", HiGreen)
		logg.colors.CustomColor("[DEBUG]", HiCyan)
	}

	Logger = logg
}

func (l *Logg) Write(b []byte) (int, error) {
	m := &message{
		data:  b,
		time:  time.Now(),
		flags: l.flags,
	}
	m.level = l.levels.define(m.data)

	if !l.levels.check(m.level) {
		return len(b), nil
	}
	if l.flags&(log.Lshortfile|log.Llongfile) != 0 {
		var ok bool
		_, file, line, ok := runtime.Caller(3)
		if !ok {
			file = "???"
			line = 0
		}
		m.file = file
		m.line = line

		if l.flags&log.Lshortfile != 0 {
			short := m.file
			for i := len(m.file) - 1; i > 0; i-- {
				if m.file[i] == '/' {
					short = m.file[i+1:]
					break
				}
			}
			m.file = short
		}
	}

	if l.format == Pretty {
		if l.color {
			var color = l.colors.define(m)
			l.set(color)
		}

		b = append(l.formatHeader(m), b...)
		n, err := l.out.Write(b)
		if l.color {
			l.unset()
		}
		return n, err
	}

	b, err := gojay.MarshalJSONObject(m)
	if err != nil {
		return len(b), err
	}

	b = append(b, []byte("\n")...)

	return l.out.Write(b)
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
func (l *Logg) formatHeader(m *message) []byte {
	t := m.time
	buf := &[]byte{}

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

	return *buf
}

// SetFormat sets the output format (Pretty or Json) for the logger.
func SetFormat(format int) {
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

// SetDebug sets the output flags prepared to debug for the logger.
func SetDebug() {
	Logger.mu.Lock()
	Logger.flags = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile
	Logger.mu.Unlock()
}

// SetColor turn on colors for the logger.
func SetColor(state bool) {
	Logger.mu.Lock()
	Logger.color = state
	Logger.mu.Unlock()
}

// set - set prefix to data with color and style
func (l *Logg) set(c string) {
	str := fmt.Sprintf("%s[%sm", escape, c)
	_, _ = fmt.Fprintf(l.out, str)
}

// unset - unset prefix from data with color and style
func (l *Logg) unset() {
	_, _ = fmt.Fprintf(l.out, "%s[%dm", escape, Reset)
}
