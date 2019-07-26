package logg

import (
	"fmt"
	"github.com/francoispqt/gojay"
	"io"
	"log"
	"os"
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
	buf []byte
}

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

	m.fileLine()
	m.level = l.levels.define(b)
	if !l.levels.check(m.level) {
		return len(b), nil
	}

	l.buf = l.buf[:0]
	if l.format == Pretty {
		l.formatHeader(m)
		l.buf = append(l.buf, m.data...)
	} else if l.format == Json {
		b, err := gojay.MarshalJSONObject(m)
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
		m.color = l.colors.define(&m.data)
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
func SetColor() {
	Logger.mu.Lock()
	Logger.color = true
	Logger.mu.Unlock()
}

// SetLevels - set the levels of logs.
func SetLevels(list []string) {
	Logger.mu.Lock()
	Logger.levels.List = list
	Logger.mu.Unlock()
}

// SetMinLevel - set the minimum levels of logs.
func SetMinLevel(minLevel string) {
	Logger.mu.Lock()
	Logger.levels.Min = minLevel
	Logger.mu.Unlock()
}

// CustomColor - allow to set custom colors for prefix.
func CustomColor(prefix string, v ...interface{}) {
	Logger.colors.CustomColor(prefix, v...)
}
