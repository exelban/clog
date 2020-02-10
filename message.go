package logg

import (
	"encoding/json"
	"log"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type message struct {
	data []byte // original log data

	color   []byte // log color
	time    []byte // message time
	level   []byte // message level
	caller  []byte // message caller
	message []byte // log message

	allowed  bool   // minimum log level
	flags    int    // time flags
	file     string // caller file path
	line     int    // caller line number
	brackets bool

	t   time.Time // log time
	buf []byte
}

const digits01 = "0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789"
const digits10 = "0000000000111111111122222222223333333333444444444455555555556666666666777777777788888888889999999999"

var messagePool = sync.Pool{
	New: func() interface{} {
		return &message{
			buf: make([]byte, 0, 500),
		}
	},
}

// reset message object before putting to sync.Pool.
func (m *message) reset() {
	m.color = []byte{}
	m.time = []byte{}
	m.message = []byte{}
	m.buf = m.buf[:0]
	m.brackets = false
}

// marshal message to json.
func (m *message) MarshalJSON() ([]byte, error) {
	jsonObject := &struct {
		Time    *time.Time `json:"time,omitempty"`
		Level   string     `json:"level,omitempty"`
		Message string     `json:"message"`
		File    string     `json:"file,omitempty"`
		Line    int        `json:"line,omitempty"`
	}{
		Level:   string(m.level),
		Message: string(m.message),
	}

	if m.flags&(log.Lshortfile|log.Llongfile) != 0 {
		jsonObject.File = m.file
		jsonObject.Line = m.line
	}

	return json.Marshal(jsonObject)
}

// receive log message from log and put to separate slice.
func (m *message) defineMessage() {
	message := m.data
	if message[len(message)-1] == '\n' {
		message = m.data[:len(m.data)-1]
	}

	for message[0] == 32 { // trip space
		message = append(message[:0], message[1:]...)
	}

	m.message = message
}

// define a time for a message.
func (m *message) defineTime() {
	if m.t.IsZero() {
		m.time = append(m.time, "0000-00-00 00:00:00"...)
		return
	}

	v := m.t
	if m.flags^log.LUTC != 0 {
		v = v.In(time.UTC)
	}

	v = v.Add(time.Nanosecond * 500) // To round under microsecond
	year := v.Year()
	year100 := year / 100
	year1 := year % 100
	month := v.Month()
	day := v.Day()
	hour := v.Hour()
	minute := v.Minute()
	second := v.Second()
	micro := v.Nanosecond() / 1000

	if m.flags&log.Ldate != 0 && m.flags&log.Ltime != 0 {
		m.time = append(m.time, []byte{
			digits10[year100], digits01[year100],
			digits10[year1], digits01[year1],
			'-',
			digits10[month], digits01[month],
			'-',
			digits10[day], digits01[day],
			' ',
			digits10[hour], digits01[hour],
			':',
			digits10[minute], digits01[minute],
			':',
			digits10[second], digits01[second],
		}...)
	} else if m.flags&log.Ldate != 0 && m.flags&log.Ltime == 0 {
		m.time = append(m.time, []byte{
			digits10[year100], digits01[year100],
			digits10[year1], digits01[year1],
			'-',
			digits10[month], digits01[month],
			'-',
			digits10[day], digits01[day],
			' ',
		}...)
	} else if m.flags&log.Ldate == 0 && m.flags&log.Ltime != 0 {
		m.time = append(m.time, []byte{
			digits10[hour], digits01[hour],
			':',
			digits10[minute], digits01[minute],
			':',
			digits10[second], digits01[second],
		}...)
	}

	if m.flags&log.Lmicroseconds != 0 {
		micro10000 := micro / 10000
		micro100 := micro / 100 % 100
		micro1 := micro % 100
		m.time = append(m.time, []byte{
			'.',
			digits10[micro10000], digits01[micro10000],
			digits10[micro100], digits01[micro100],
			digits10[micro1], digits01[micro1],
		}...)
	}
}

// define a caller for a message.
func (m *message) defineCaller() {
	if m.flags&(log.Lshortfile|log.Llongfile) == 0 {
		return
	}

	_, file, line, ok := runtime.Caller(5)
	if !ok {
		file = "???"
		line = 0
	}

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

	m.caller = append(m.caller, file+":"+strconv.Itoa(line)...)
}
