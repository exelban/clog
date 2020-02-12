package logg

import (
	"log"
	"strconv"
	"sync"
	"time"
)

type message struct {
	level     level
	calldepth int
	flags     int
	format    format
	color     bool

	buf []byte
}

var messagePool = sync.Pool{
	New: func() interface{} {
		return &message{
			buf:   make([]byte, 0, 500),
			level: Empty,
		}
	},
}

// fetch a message from sync.Pool.
func newMessage(level level, calldepth int, flags int, format format, color bool) *message {
	m := messagePool.Get().(*message)

	m.level = level
	m.calldepth = calldepth
	m.flags = flags
	m.format = format
	m.color = color

	m.buf = m.buf[:0]

	return m
}

// reset message object and put to sync.Pool.
func (m *message) put() {
	const maxSize = 1 << 16 // 64KiB
	if cap(m.buf) > maxSize {
		return
	}

	messagePool.Put(m)
}

func (m *message) build(b []byte) []byte {
	if len(b) != 0 {
		if m.format == Json {
			m.buildJSON(b)
		} else {
			m.buildPretty(b)
		}
	}

	return append(m.buf, '\n')
}

func (m *message) buildJSON(b []byte) {
	js := newJson()

	if m.flags&(log.Ldate|log.Ltime|log.Lmicroseconds) != 0 {
		js.buf = append(js.addField("time", js.buf), time.Now().Format(time.RFC3339)...)
	}

	if m.flags&(log.Lshortfile|log.Llongfile) != 0 {
		file, line := caller(m.calldepth, m.flags&log.Lshortfile != 0)

		js.buf = append(js.addField("file", js.buf), file...)
		js.buf = append(js.addField("line", js.buf), strconv.Itoa(line)...)
	}

	if m.level != Empty {
		js.buf = append(js.addField("level", js.buf), levels[m.level]...)
	}

	if len(b) != 0 {
		js.buf = append(js.addField("message", js.buf), b...)
	}

	js.close()
	m.buf = js.buf

	js.put()
}

func (m *message) buildPretty(b []byte) {
	if m.color && m.level != Empty {
		m.buf = append(m.buf, colors[m.level]...)
	}

	if m.flags&(log.Ldate|log.Ltime|log.Lmicroseconds) != 0 {
		m.buf = appendTimestamp(time.Now(), m.flags, m.buf)
	}

	if m.flags&(log.Lshortfile|log.Llongfile) != 0 {
		file, line := caller(m.calldepth, m.flags&log.Lshortfile != 0)

		if len(m.buf) != 0 && m.buf[len(m.buf)-1] != ' ' {
			m.buf = append(m.buf, ' ')
		}

		m.buf = append(m.buf, file+":"...)
		m.buf = append(m.buf, strconv.Itoa(line)...)
	}

	if m.level != Empty {
		if m.flags&(log.Ldate|log.Ltime|log.Lmicroseconds) != 0 || m.flags&(log.Lshortfile|log.Llongfile) != 0 {
			m.buf = append(m.buf, ' ')
		}
		m.buf = append(m.buf, levels[m.level]...)
	}

	if len(b) != 0 {
		if len(m.buf) != 0 && m.buf[len(m.buf)-1] != ' ' {
			m.buf = append(m.buf, ' ')
		}

		m.buf = append(m.buf, b...)
	}

	if m.color && m.level != Empty {
		m.buf = append(m.buf, []byte(escapeClose)...)
	}
}
