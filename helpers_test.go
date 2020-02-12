package logg

import (
	"bytes"
	"log"
	"runtime"
	"strings"
	"testing"
	"time"
)

func Test_caller(t *testing.T) {
	file, line := caller(60, false)
	if file != "???" {
		t.Errorf("caller path must be undefined. Received: %s", file)
	}
	if line != 0 {
		t.Errorf("caller line must be 0. Received: %d", line)
	}

	file, line = caller(1, false)
	_, f, l, _ := runtime.Caller(0)
	if file != f {
		t.Errorf("wrong caller path. Expected: %s, received: %s", f, file)
	}
	if line != l-1 {
		t.Errorf("wrong caller line. Expected: %d, received: %d", l-1, line)
	}

	file, line = caller(1, true)
	_, f, l, _ = runtime.Caller(0)
	if !strings.Contains(f, file) {
		t.Errorf("wrong caller path. Expected: %s, received: %s", f, file)
	}
	if line != l-1 {
		t.Errorf("wrong caller line. Expected: %d, received: %d", l-1, line)
	}
}

func Test_timestamp(t *testing.T) {
	now := time.Now()
	tests := map[string]struct {
		t     time.Time
		flags int
		buf   []byte
	}{
		"empty": {
			t:     time.Time{},
			flags: 0,
			buf:   []byte("0000-00-00 00:00:00"),
		},
		"default format": {
			t:     now,
			flags: log.LstdFlags,
			buf:   []byte(now.Format("2006-01-02 15:04:05")),
		},
		"date only": {
			t:     now,
			flags: log.Ldate,
			buf:   []byte(now.Format("2006-01-02")),
		},
		"time only": {
			t:     now,
			flags: log.Ltime,
			buf:   []byte(now.Format("15:04:05")),
		},
		"milliseconds only": {
			t:     now,
			flags: log.Lmicroseconds,
			buf:   []byte(now.Format(".999999")),
		},
		"time with milliseconds": {
			t:     now,
			flags: log.Ltime | log.Lmicroseconds,
			buf:   []byte(now.Format("15:04:05.999999")),
		},
		"date with time with milliseconds": {
			t:     now,
			flags: log.Ldate | log.Ltime | log.Lmicroseconds,
			buf:   []byte(now.Format("2006-01-02 15:04:05.999999")),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			buf := []byte{}
			buf = appendTimestamp(tc.t, tc.flags, buf)

			if !bytes.Equal(buf, tc.buf) {
				t.Errorf("wrong timestamp. Expected: %s, received: %s", string(tc.buf), string(buf))
			}
		})
	}
}

func Benchmark_appendTimestamp(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	now := time.Now()
	buf := []byte{}

	for n := 0; n < b.N; n++ {
		appendTimestamp(now, log.LstdFlags, buf)
	}
}

func Test_defineLevel(t *testing.T) {
	tests := map[string]struct {
		data   []byte
		result []byte
		level  level
	}{
		"empty": {
			data:   []byte{},
			result: []byte{},
			level:  Empty,
		},
		"empty with string": {
			data:   []byte("test"),
			result: []byte("test"),
			level:  Empty,
		},
		"short": {
			data:   []byte("INF test"),
			result: []byte("test"),
			level:  Info,
		},
		"short with brackets": {
			data:   []byte("[ERR] test"),
			result: []byte("test"),
			level:  Error,
		},
		"double short": {
			data:   []byte("WRN ERR test"),
			result: []byte("ERR test"),
			level:  Warning,
		},
		"long": {
			data:   []byte("DEBUG test"),
			result: []byte("test"),
			level:  Debug,
		},
		"long with brackets": {
			data:   []byte("[PANIC] test"),
			result: []byte("test"),
			level:  Panic,
		},
		"double long": {
			data:   []byte("PANIC DEBUG test"),
			result: []byte("DEBUG test"),
			level:  Panic,
		},

		"debug short": {
			data:   []byte("DBG test"),
			result: []byte("test"),
			level:  Debug,
		},
		"debug long": {
			data:   []byte("DEBUG test"),
			result: []byte("test"),
			level:  Debug,
		},
		"info short": {
			data:   []byte("INF test"),
			result: []byte("test"),
			level:  Info,
		},
		"info long": {
			data:   []byte("INFO test"),
			result: []byte("test"),
			level:  Info,
		},
		"error short": {
			data:   []byte("ERR test"),
			result: []byte("test"),
			level:  Error,
		},
		"error long": {
			data:   []byte("ERROR test"),
			result: []byte("test"),
			level:  Error,
		},
		"warning short": {
			data:   []byte("WRN test"),
			result: []byte("test"),
			level:  Warning,
		},
		"warning long": {
			data:   []byte("WARN test"),
			result: []byte("test"),
			level:  Warning,
		},
		"panic short": {
			data:   []byte("PNC test"),
			result: []byte("test"),
			level:  Panic,
		},
		"panic long": {
			data:   []byte("PANIC test"),
			result: []byte("test"),
			level:  Panic,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			lvl := defineLevel(&tc.data)
			tc.data = removeLevel(tc.data, lvl)

			if lvl != tc.level {
				t.Errorf("wrong level. Expected: %d, received: %d", tc.level, lvl)
			}

			if !bytes.Equal(tc.data, tc.result) {
				t.Errorf("level not removed from data. Expected: %s, received: %s", string(tc.result), string(tc.data))
			}
		})
	}
}

func Benchmark_defineLevel(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		t := []byte("INF test")
		_ = defineLevel(&t)
	}
}

func Benchmark_removeLevel(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		t := []byte("INFO test")
		t = removeLevel(t, Info)
	}
}
