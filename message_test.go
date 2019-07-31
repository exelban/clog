package logg

import (
	"encoding/json"
	"log"
	"testing"
	"time"
)

type jsonLog struct {
	Time    time.Time `json:"time"`
	Level   string    `json:"level"`
	Message string    `json:"message"`
	File    string    `json:"file"`
	Line    int       `json:"line"`
}

func TestMessage_MarshalJSONObject(t *testing.T) {
	tests := map[string]struct {
		m               *message
		expectedMessage string
		expectedLevel   string
		expectedFile    string
		expectedLine    int
	}{
		"test1": {
			m: &message{
				time:  time.Now(),
				level: "",
				data:  []byte(" test 1"),
				file:  "",
				line:  0,
			},
			expectedMessage: " test 1",
			expectedLevel:   "",
			expectedFile:    "",
			expectedLine:    0,
		},
		"test2": {
			m: &message{
				time:  time.Now(),
				level: "",
				data:  []byte("[INFO] test 2"),
				file:  "",
				line:  0,
			},
			expectedMessage: "[INFO] test 2",
			expectedLevel:   "",
			expectedFile:    "",
			expectedLine:    0,
		},
		"test3": {
			m: &message{
				time:  time.Now(),
				level: "INFO",
				data:  []byte("[INFO] test 3"),
				file:  "",
				line:  0,
			},
			expectedMessage: "test 3",
			expectedLevel:   "INFO",
			expectedFile:    "",
			expectedLine:    0,
		},
		"test4": {
			m: &message{
				time:  time.Now(),
				level: "INFO",
				data:  []byte("INFO test 4"),
				file:  "",
				line:  0,
			},
			expectedMessage: "test 4",
			expectedLevel:   "INFO",
			expectedFile:    "",
			expectedLine:    0,
		},
		"test5": {
			m: &message{
				time:  time.Now(),
				level: "WoW",
				data:  []byte("WoW test 5"),
				file:  "",
				line:  0,
			},
			expectedMessage: "test 5",
			expectedLevel:   "WoW",
			expectedFile:    "",
			expectedLine:    0,
		},
		"test6": {
			m: &message{
				time:  time.Now(),
				level: "WoW",
				data:  []byte("test 6 WoW"),
				file:  "",
				line:  0,
			},
			expectedMessage: "test 6 WoW",
			expectedLevel:   "WoW",
			expectedFile:    "",
			expectedLine:    0,
		},
		"test7": {
			m: &message{
				time:  time.Now(),
				level: "info",
				data:  []byte("[info] test 7"),
				file:  "",
				line:  0,
			},
			expectedMessage: "test 7",
			expectedLevel:   "info",
			expectedFile:    "",
			expectedLine:    0,
		},
		"test8": {
			m: &message{
				time:  time.Now(),
				level: "info",
				data:  []byte("info test 8"),
				file:  "",
				line:  0,
			},
			expectedMessage: "test 8",
			expectedLevel:   "info",
			expectedFile:    "",
			expectedLine:    0,
		},
		"test9": {
			m: &message{
				time:  time.Now(),
				level: "DEBUG",
				data:  []byte("debug test 9"),
				file:  "main.go",
				line:  0,
			},
			expectedMessage: "debug test 9",
			expectedLevel:   "DEBUG",
			expectedFile:    "",
			expectedLine:    0,
		},
		"test10": {
			m: &message{
				time:  time.Now(),
				level: "DEBUG",
				data:  []byte("debug test 10"),
				file:  "main.go",
				line:  0,
				flags: log.Lshortfile,
			},
			expectedMessage: "debug test 10",
			expectedLevel:   "DEBUG",
			expectedFile:    "main.go",
			expectedLine:    0,
		},
		"test11": {
			m: &message{
				time:  time.Now(),
				level: "DEBUG",
				data:  []byte("debug test 11"),
				file:  "main.go",
				line:  125,
				flags: log.Lshortfile,
			},
			expectedMessage: "debug test 11",
			expectedLevel:   "DEBUG",
			expectedFile:    "main.go",
			expectedLine:    125,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			b, err := tc.m.MarshalJSON()
			if err != nil {
				t.Error(err)
			}

			var response jsonLog
			err = json.Unmarshal(b, &response)
			if err != nil {
				t.Error(err)
			}

			if response.Level != tc.expectedLevel {
				t.Errorf("wrong level: expected %v, received: %v", tc.expectedLevel, response.Level)
			}
			if response.Message != tc.expectedMessage {
				t.Errorf("wrong message: expected %v, received: %v", tc.expectedMessage, response.Message)
			}
			if response.File != tc.expectedFile {
				t.Errorf("wrong file: expected %v, received: %v", tc.expectedFile, response.File)
			}
			if response.Line != tc.expectedLine {
				t.Errorf("wrong line: expected %v, received: %v", tc.expectedLine, response.Line)
			}
		})
	}
}

func TestMessage_fileLine(t *testing.T) {
	tests := map[string]struct {
		m            *message
		expectedFile string
		expectedLine int
		calldepth    int
	}{
		"test1": {
			m:            &message{},
			expectedFile: "",
			expectedLine: 0,
		},
		"test2": {
			m: &message{
				file:  "",
				line:  0,
				flags: log.Lshortfile,
			},
			expectedFile: "message_test.go",
			expectedLine: 0,
			calldepth:    1,
		},
		"test3": {
			m: &message{
				file:  "",
				line:  0,
				flags: log.Lshortfile,
			},
			expectedFile: "message_test.go",
			expectedLine: 255,
			calldepth:    1,
		},
		"test4": {
			m: &message{
				file:  "",
				line:  0,
				flags: log.Lshortfile,
			},
			expectedFile: "???",
			expectedLine: 0,
			calldepth:    4,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.m.fileLine(tc.calldepth)

			if tc.m.file != tc.expectedFile {
				t.Errorf("wrong file: expected %v, received: %v", tc.expectedFile, tc.m.file)
			}
			if tc.expectedLine != 0 && tc.m.line == 0 {
				t.Errorf("wrong line: expected %v, received: %v", tc.expectedLine, tc.m.line)
			}
		})
	}
}
