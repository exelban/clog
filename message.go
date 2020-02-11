package logg

import (
	"log"
	"strconv"
	"sync"
	"time"
)

type message struct {
	data []byte // original log data

	time    []byte // message time
	level   []byte // message level
	file    []byte // caller file path
	line    []byte // caller line number
	message []byte // log message

	t         time.Time
	lvl       level
	flags     int
	calldepth int

	buf []byte
}

var messagePool = sync.Pool{
	New: func() interface{} {
		return &message{
			buf: make([]byte, 0, 500),
			lvl: Empty,
		}
	},
}

// fetch a message from sync.Pool.
func newMessage(b []byte, level level, flags int, calldepth int) *message {
	m := messagePool.Get().(*message)
	m.t = time.Now()
	if flags&log.LUTC != 0 {
		m.t = m.t.In(time.UTC)
	}

	m.data = b
	m.flags = flags
	m.calldepth = ContextCallDepth + calldepth
	m.lvl = level

	return m
}

// reset message object and put to sync.Pool.
func (m *message) put() {
	m.data = m.data[:0]

	m.time = []byte{}
	m.level = []byte{}
	m.line = []byte{}
	m.file = []byte{}
	m.message = []byte{}

	m.lvl = Empty
	m.calldepth = 0
	m.flags = 0

	m.buf = m.buf[:0]

	const maxSize = 1 << 16 // 64KiB
	if cap(m.buf) > maxSize {
		return
	}
	messagePool.Put(m)
}

// Parse a provided byte array and try to fetch time, caller and message from log.
func (m *message) build() {
	if m.lvl == Empty {
		m.lvl = defineLevel(&m.data)
	}

	if m.flags&(log.Lshortfile|log.Llongfile) != 0 {
		file, line := caller(m.calldepth, m.flags&log.Lshortfile != 0)

		m.line = append(m.line, strconv.Itoa(line)...)
		m.file = append(m.file, file...)
	}

	if m.flags != 0 {
		m.time = timestamp(m.t, m.flags)
	}

	m.message = m.data
}

// Build a log from message.
func (m *message) exec(format format, color bool) {
	if format == Json {
		js := newJSON()

		if len(m.time) != 0 {
			js.addField("time", m.time)
		}
		if m.lvl != Empty {
			js.addField("level", []byte(levels[m.lvl]))
		}
		if len(m.message) != 0 {
			js.addField("message", m.message)
		}
		if len(m.file) != 0 {
			js.addField("file", m.file)
		}
		if len(m.line) != 0 {
			js.addField("line", m.line)
		}

		js.close()
		m.buf = append(m.buf, js.buf...)

		js.put()
	} else {
		// add time
		if len(m.time) != 0 {
			m.buf = append(m.buf, m.time...)
			m.buf = append(m.buf, ' ')
		}

		// add caller
		if len(m.file) != 0 {
			m.buf = append(m.buf, m.file...)
			if len(m.line) != 0 {
				m.buf = append(m.buf, ':')
				m.buf = append(m.buf, m.line...)
			}
			m.buf = append(m.buf, ' ')
		}

		//add level
		if m.lvl != Empty {
			m.buf = append(m.buf, levels[m.lvl]...)
		}

		// add message
		if len(m.message) > 0 && len(m.buf) > 0 && m.message[0] != ' ' && m.buf[len(m.buf)-1] != ' ' {
			m.buf = append(m.buf, ' ')
		}
		m.buf = append(m.buf, m.message...)

		// colorize output
		if color && m.lvl != Empty {
			colorize(&m.buf, colors[int(m.lvl)])
		}
	}

	if len(m.buf) == 0 || m.buf[len(m.buf)-1] != '\n' {
		m.buf = append(m.buf, '\n')
	}
}
