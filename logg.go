package logg

import (
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

	out io.Writer
	mu  sync.Mutex
}

var Logger *Logg

func init() {
	logg := &Logg{
		colors: ColorsManager{},
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

		out: os.Stderr,
	}

	log.SetOutput(logg)
	log.SetFlags(0)

	Logger = logg
}

func (l *Logg) Write(b []byte) (int, error) {
	m := &message{
		data: b,
		time: time.Now(),
	}
	m.level = l.levels.define(m.data)

	if !l.levels.check(m.level) {
		return len(b), nil
	}

	if l.format == Pretty {
		timePrefix := time.Now().Format(time.RFC3339)
		b = append([]byte(timePrefix+" "), b...)
		return l.out.Write(b)
	}

	b, err := gojay.MarshalJSONObject(m)
	if err != nil {
		return len(b), err
	}

	b = append(b, []byte("\n")...)

	return l.out.Write(b)
}

func SetFormat(format int) {
	Logger.mu.Lock()
	Logger.format = format
	Logger.mu.Unlock()
}

func SetFlags(flags int) {
	Logger.mu.Lock()
	Logger.flags = flags
	Logger.mu.Unlock()
}
