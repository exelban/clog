package logg

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"
)

func TestLogg_newMessage(t *testing.T) {
	tests := map[string]struct {
		data      []byte
		level     level
		flags     int
		calldepth int
	}{
		"empty": {
			data:      []byte{},
			level:     Empty,
			flags:     0,
			calldepth: 0,
		},
		"not empty bytes": {
			data:      []byte{'0', '1', '2'},
			level:     Empty,
			flags:     0,
			calldepth: 0,
		},
		"info level": {
			data:      []byte("test"),
			level:     Info,
			flags:     0,
			calldepth: 0,
		},
		"some flags": {
			data:      []byte("test"),
			level:     Info,
			flags:     log.LstdFlags,
			calldepth: 0,
		},
		"some calldepth": {
			data:      []byte("test"),
			level:     Panic,
			flags:     log.LstdFlags,
			calldepth: 4,
		},
		"UTC time (must be not UTC)": {
			data:      []byte("test"),
			level:     Panic,
			flags:     log.Ldate,
			calldepth: 4,
		},
		"UTC time (must be UTC)": {
			data:      []byte("test"),
			level:     Panic,
			flags:     log.LUTC,
			calldepth: 4,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			m := newMessage(tc.data, tc.level, tc.flags, tc.calldepth)

			if m == nil {
				t.Error("message object is nil")
				t.FailNow()
			}

			if len(m.buf) != 0 {
				t.Error("message buffer must be empty")
			}

			if string(m.data) != string(tc.data) {
				t.Errorf("message data in not the same as expected. Expected: %s, received: %s", string(tc.data), string(m.data))
			}

			if m.calldepth != ContextCallDepth+tc.calldepth {
				t.Errorf("message calldepth must be equal to %d, received: %d", ContextCallDepth+tc.calldepth, m.calldepth)
			}

			if m.lvl != tc.level {
				t.Errorf("message lvl must be equal to %d, received: %d", tc.level, m.lvl)
			}

			now := time.Now()
			if m.t.After(now) {
				t.Errorf("message timestamp is too old. Expected: %v, received: %v", now, m.t)
			}

			zone, _ := m.t.Zone()
			localZone, _ := time.Now().Zone()
			if tc.flags&log.LUTC != 0 && zone != "UTC" {
				t.Errorf("message timestamp must be in UTC zone. Received: %s", zone)
			} else if tc.flags&log.LUTC == 0 && zone != localZone {
				t.Errorf("message timestamp must be in %s zone. Received: %s", localZone, zone)
			}
		})
	}
}

func TestLogg_message_put(t *testing.T) {
	tests := map[string]*message{
		"empty": {
			data: []byte{},

			time:    []byte{},
			level:   []byte{},
			file:    []byte{},
			line:    []byte{},
			message: []byte{},

			t:         time.Now(),
			lvl:       Empty,
			flags:     0,
			calldepth: 0,

			buf: []byte{},
		},
		"not empty data and buf": {
			data: []byte{'0'},

			time:    []byte{},
			level:   []byte{},
			file:    []byte{},
			line:    []byte{},
			message: []byte{},

			t:         time.Now(),
			lvl:       Info,
			flags:     0,
			calldepth: 0,

			buf: []byte{'{'},
		},
		"main arrays": {
			data: []byte{'0'},

			time:    []byte("time"),
			level:   []byte("level"),
			file:    []byte("file"),
			line:    []byte("line"),
			message: []byte("message"),

			t:         time.Now(),
			lvl:       Info,
			flags:     0,
			calldepth: 0,

			buf: []byte{'{'},
		},
		"some flags": {
			data: []byte{'0'},

			time:    []byte("time"),
			level:   []byte("level"),
			file:    []byte("file"),
			line:    []byte("line"),
			message: []byte("message"),

			t:         time.Now(),
			lvl:       Panic,
			flags:     log.LstdFlags,
			calldepth: 4,

			buf: []byte{'{'},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.put()

			m := messagePool.Get().(*message)

			if len(m.data) != 0 {
				t.Errorf("message data must empty, received: %s", string(m.data))
			}

			if len(m.time) != 0 {
				t.Errorf("message time must empty, received: %s", string(m.time))
			}

			if len(m.level) != 0 {
				t.Errorf("message level must empty, received: %s", string(m.level))
			}

			if len(m.line) != 0 {
				t.Errorf("message line must empty, received: %s", string(m.line))
			}

			if len(m.file) != 0 {
				t.Errorf("message file must empty, received: %s", string(m.file))
			}

			if len(m.message) != 0 {
				t.Errorf("message message must empty, received: %s", string(m.message))
			}

			if m.lvl != Empty {
				t.Errorf("message lvl must Empty, received: %d", m.lvl)
			}

			if m.calldepth != 0 {
				t.Errorf("message calldepth must 0, received: %d", m.calldepth)
			}

			if m.flags != 0 {
				t.Errorf("message flags must 0, received: %d", m.flags)
			}
		})
	}

	m := &message{
		buf: make([]byte, 1<<16+1),
		t:   time.Now(),
	}
	m.put()

	mFromPool := messagePool.Get().(*message)
	if m.t == mFromPool.t {
		t.Error("message has a buf size maxSize+1, must not be saved in sync.Pool")
	}
}

func TestLogg_message_build(t *testing.T) {
	tests := map[string]struct {
		data      []byte
		lvl       level
		flags     int
		calldepth int

		color   []byte
		time    []byte
		message []byte

		expectedLevel level
	}{
		"empty": {
			data:      []byte{},
			lvl:       Empty,
			flags:     0,
			calldepth: 0,

			expectedLevel: Empty,
		},
		"string log": {
			data:      []byte("test"),
			lvl:       Empty,
			flags:     0,
			calldepth: 0,

			message:       []byte("test"),
			expectedLevel: Empty,
		},
		"level log": {
			data:      []byte("INF test"),
			lvl:       Empty,
			flags:     0,
			calldepth: 0,

			message:       []byte("test"),
			expectedLevel: Info,
		},
		"level log with predefined level": {
			data:      []byte("INF test"),
			lvl:       Error,
			flags:     0,
			calldepth: 0,

			message:       []byte("INF test"),
			expectedLevel: Error,
		},
		"caller": {
			data:      []byte("INF test"),
			lvl:       Empty,
			flags:     log.Lshortfile,
			calldepth: -2,

			message:       []byte("test"),
			expectedLevel: Info,
		},
		"timestamp": {
			data:  []byte("INF test"),
			lvl:   Empty,
			flags: log.LstdFlags,

			message:       []byte("test"),
			expectedLevel: Info,
			time:          []byte(time.Now().Format("2006-01-02 15:04:05")),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			m := newMessage(tc.data, tc.lvl, tc.flags, tc.calldepth)

			m.build()
			file, line := caller(1, tc.flags&log.Lshortfile != 0)

			if !bytes.Equal(m.time, tc.time) {
				t.Errorf("time is not equeal. Expected: %s, received: %s", string(tc.time), string(m.time))
			}

			if m.lvl != tc.expectedLevel {
				t.Errorf("level is not equeal. Expected: %d, received: %d", tc.expectedLevel, m.lvl)
			}

			if tc.flags&(log.Lshortfile|log.Llongfile) != 0 {
				if string(m.file) != file {
					t.Errorf("file is not equeal. Expected: %s, received: %s", file, string(m.file))
				}
				if string(m.line) != fmt.Sprint(line-1) {
					t.Errorf("line is not equeal. Expected: %s, received: %s", fmt.Sprint(line-1), string(m.line))
				}
			}

			if !bytes.Equal(m.message, tc.message) {
				t.Errorf("message is not equeal. Expected: %s, received: %s", string(tc.message), string(m.message))
			}

			m.put()
		})
	}
}

func TestLogg_message_exec(t *testing.T) {
	tests := map[string]struct {
		data      []byte
		flags     int
		calldepth int

		pretty []byte
		json   []byte
	}{
		"empty": {
			data:      []byte{},
			flags:     0,
			calldepth: 0,

			pretty: []byte("\n"),
			json:   []byte(`{}`),
		},
		"string": {
			data:      []byte("test"),
			flags:     0,
			calldepth: 0,

			pretty: []byte("test\n"),
			json:   []byte(`{"message": "test"}`),
		},
		"level": {
			data:      []byte("INF test"),
			flags:     0,
			calldepth: 0,

			pretty: []byte("INF test\n"),
			json:   []byte(`{"level": "INF", "message": "test"}`),
		},
		"level with brackets": {
			data:      []byte("[INFO] test"),
			flags:     0,
			calldepth: 0,

			pretty: []byte("INF test\n"),
			json:   []byte(`{"level": "INF", "message": "test"}`),
		},
		"time": {
			data:      []byte("INF test"),
			flags:     log.LstdFlags,
			calldepth: 0,

			pretty: []byte(fmt.Sprintf("%s INF test\n", time.Now().Format("2006-01-02 15:04:05"))),
			json:   []byte(fmt.Sprintf(`{"time": "%s", "level": "INF", "message": "test"}`, time.Now().Format("2006-01-02 15:04:05"))),
		},
		"caller": {
			data:      []byte("INF test"),
			flags:     log.Lshortfile,
			calldepth: -2,

			pretty: []byte("%1:%2 INF test\n"),
			json:   []byte(`{"level": "INF", "message": "test", "file": "%1", "line": "%2"}`),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			m := newMessage(tc.data, Empty, tc.flags, tc.calldepth)

			m.build()
			file, line := caller(1, tc.flags&log.Lshortfile != 0)
			if tc.flags&(log.Lshortfile|log.Llongfile) != 0 {
				tc.pretty = []byte(strings.Replace(string(tc.pretty), "%1", file, 1))
				tc.pretty = []byte(strings.Replace(string(tc.pretty), "%2", fmt.Sprint(line-1), 1))

				tc.json = []byte(strings.Replace(string(tc.json), "%1", file, 1))
				tc.json = []byte(strings.Replace(string(tc.json), "%2", fmt.Sprint(line-1), 1))
			}

			m.exec(Pretty, false)
			if !bytes.Equal(m.buf, tc.pretty) {
				t.Errorf("pretty: output is not the same. Expected: %s, received: %s", string(tc.pretty), string(m.buf))
			}

			m.buf = m.buf[:0]
			m.exec(Json, false)
			tc.json = append(tc.json, 10)
			if !bytes.Equal(m.buf, tc.json) {
				t.Errorf("json: output is not the same. Expected: %s, received: %s", string(tc.json), string(m.buf))
			}

			m.put()
		})
	}
}

func TestLogg_message_exec_color(t *testing.T) {
	data := []byte("INF test")

	input := make([]byte, len(data))
	copy(input, data)
	m := newMessage(input, Empty, 0, 0)
	m.build()
	m.exec(Pretty, true)

	colorize(&data, colors[Info])
	data = append(data, '\n')

	if !bytes.Equal(data, m.buf) {
		t.Errorf("Wrong color. Expected: %s, received: %s", string(data), string(m.buf))
	}

	m.buf = m.buf[:0]
	m.exec(Json, true)
	expected := []byte(`{"level": "INF", "message": "test"}`)
	expected = append(expected, '\n')
	if !bytes.Equal(expected, m.buf) {
		t.Errorf("Json cannot be colorized. Expected: %s, received: %s", string(expected), string(m.buf))
	}

	m.put()
}
