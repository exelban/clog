package logg

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func (m *message) MarshalJSON() ([]byte, error) {
	jsonObject := &struct {
		Time    *time.Time `json:"time,omitempty"`
		Level   string     `json:"level,omitempty"`
		Message string     `json:"message"`
		File    string     `json:"file,omitempty"`
		Line    int        `json:"line,omitempty"`
	}{
		Level:   m.level,
		Message: m.getMessage(),
	}

	if m.flags != 0 {
		t := m.time.UTC()
		jsonObject.Time = &t
	}
	if m.flags&(log.Lshortfile|log.Llongfile) != 0 {
		jsonObject.File = m.file
		jsonObject.Line = m.line
	}

	return json.Marshal(jsonObject)
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

func (m *message) fileLine(calldepth int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.flags&(log.Lshortfile|log.Llongfile) != 0 {
		m.mu.Unlock()
		var ok bool
		var file string
		var line int
		_, file, line, ok = runtime.Caller(calldepth)
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
