package logg

import (
	"bytes"
	"fmt"
	"github.com/francoispqt/gojay"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"
)

type message struct {
	data  []byte
	time  time.Time
	level string
	color string

	flags int
	file  string
	line  int

	mu sync.Mutex
}

func (m *message) MarshalJSONObject(enc *gojay.Encoder) {
	enc.TimeKey("time", &m.time, "2006-01-02T15:04:05.000Z")
	enc.StringKey("level", m.level)
	enc.StringKey("message", m.getMessage())

	if m.flags&(log.Lshortfile|log.Llongfile) != 0 {
		enc.StringKey("file", m.file)
		enc.IntKey("line", m.line)
	}
}
func (m *message) IsNil() bool {
	return m == nil
}

func (m *message) getMessage() string {
	mess := string(m.data)
	mess = strings.TrimSuffix(mess, "\n")

	if m.level == "" {
		return mess
	}

	level := fmt.Sprintf("[%v]", m.level)
	if !bytes.Contains(m.data, []byte(level)) {
		level = m.level
	}

	mess = strings.TrimPrefix(mess, level)
	if mess[0] == 32 { // trip space
		mess = strings.TrimPrefix(mess, " ")
	}

	return mess
}

func (m *message) fileLine() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.flags&(log.Lshortfile|log.Llongfile) != 0 {
		m.mu.Unlock()
		var ok bool
		var file string
		var line int
		_, file, line, ok = runtime.Caller(4)
		if !ok {
			file = "???"
			line = 0
		}
		m.mu.Lock()

		if m.flags&log.Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}

		m.line = line
		m.file = file
	}
}
