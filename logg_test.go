package logg

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestLogg(t *testing.T) {
	buf := new(bytes.Buffer)

	tests := map[string]struct {
		data   string
		prefix string
		style  string

		format format
		flags  int
		color  bool
		error  bool

		expectedData string
		emptyOutput  bool
	}{
		"empty": {
			data:         "",
			expectedData: "",
			format:       Pretty,
			color:        true,
		},
		"simple text": {
			data:         "Hello World",
			expectedData: "Hello World",
			format:       Pretty,
			color:        true,
		},
		"INFO": {
			data:         "[INFO] Hello World",
			expectedData: "[INFO] Hello World",
			format:       Pretty,
			color:        true,
		},
		"DEBUG": {
			data:        "[DEBUG] min level test",
			emptyOutput: true,
			format:      Pretty,
			color:       true,
		},
		"ERROR": {
			data:         "[ERROR] min level test",
			expectedData: "[ERROR] min level test",
			format:       Pretty,
			color:        true,
		},
		"ERROR_2": {
			data:         "ERROR min level test",
			expectedData: "ERROR min level test",
			format:       Pretty,
			color:        true,
		},
		"empty_json": {
			data:         "",
			expectedData: `{"message":""}`,
			format:       Json,
			color:        false,
		},
		"simple_json": {
			data:         "Hello World",
			expectedData: `{"message":"Hello World"}`,
			format:       Json,
			color:        false,
		},
		"INFO_json": {
			data:         "[INFO] Hello World",
			expectedData: `{"level":"INFO","message":"Hello World"}`,
			format:       Json,
			color:        false,
		},
		"DEBUG_json": {
			data:        "[DEBUG] Hello World",
			emptyOutput: true,
			format:      Json,
			color:       false,
		},
	}

	SetOutput(buf)
	SetFlags(0)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			SetFormat(tc.format)
			Logger.color = tc.color

			log.Print(tc.data)

			b := []byte(tc.data)
			expectedOutput := tc.expectedData
			if tc.color {
				expectedOutput = fmt.Sprintf("\x1b[%sm%s", Logger.colors.define(&b), tc.expectedData)
			}
			output := readFromBuffer(buf)

			//fmt.Println(Logger.buf, []byte(output), []byte(expectedOutput))

			if !tc.emptyOutput && output != expectedOutput {
				t.Errorf("expected: %s, writed: %s", expectedOutput, output)
			} else if tc.emptyOutput && len([]byte(output)) != 0 {
				t.Errorf("expected empty output, writed: %s", output)
			}
		})
	}

	SetOutput(os.Stderr)
	SetFlags(log.Ldate | log.Ltime)
	SetFormat(Pretty)
	Logger.color = true
}

func TestLogg_Write(t *testing.T) {
	tests := map[string]struct {
		data   []byte
		time   bool
		prefix bool
	}{
		"empty": {
			data:   []byte(""),
			time:   true,
			prefix: false,
		},
		"1": {
			data:   []byte("1"),
			time:   true,
			prefix: false,
		},
		"zero": {
			data:   []byte("zero"),
			time:   true,
			prefix: false,
		},
		"info": {
			data:   []byte("[INFO] Hello world"),
			time:   true,
			prefix: true,
		},
		"error": {
			data:   []byte("[ERROR] Hello world"),
			time:   true,
			prefix: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			n, err := Logger.Write(tc.data)
			if err != nil {
				t.Error(err)
			}

			expectedBytes := len(tc.data)
			if tc.time {
				expectedBytes += 24
			}
			if tc.prefix {
				expectedBytes += 3
			}

			if expectedBytes != n {
				t.Errorf("expected %d bytes, writed: %d", expectedBytes, n)
			}
		})
	}
}

func TestCustomColor(t *testing.T) {
	tests := map[string]struct {
		prefix string
		color  int
	}{
		"HIDDEN": {
			prefix: "[HIDDEN]",
			color:  Black,
		},
		"ERROR": {
			prefix: "[ERROR]",
			color:  Red,
		},
		"INFO": {
			prefix: "[INFO]",
			color:  Yellow,
		},
		"WARN": {
			prefix: "[WARN]",
			color:  Green,
		},
		"DEBUG": {
			prefix: "[DEBUG]",
			color:  Cyan,
		},
		"PANIC": {
			prefix: "[PANIC]",
			color:  Blue,
		},
		"OWN": {
			prefix: "[OWN]",
			color:  Magenta,
		},
		"TEST": {
			prefix: "[TEST]",
			color:  White,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			CustomColor(tc.prefix, tc.color)

			if Logger.colors.list[tc.prefix] == "" {
				t.Errorf("color not find: expected %v, received: %v", tc.prefix, Logger.colors.list[tc.prefix])
			}
		})
	}
}

func TestSetDebug(t *testing.T) {
	if Logger.flags != log.Ldate|log.Ltime {
		t.Errorf("default flags must be %v, not: %v", log.Ldate|log.Ltime, Logger.flags)
	}

	SetDebug()

	if Logger.flags != log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile {
		t.Errorf("flags must be %v, not: %v", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile, Logger.flags)
	}
}

func TestSetFormat(t *testing.T) {
	if Logger.format != Pretty {
		t.Errorf("default format must be %v, not: %v", Pretty, Logger.format)
	}

	SetFormat(Json)

	if Logger.format != Json {
		t.Errorf("format must be %v, not: %v", Json, Logger.format)
	}
}

func TestSetFlags(t *testing.T) {
	SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	if Logger.flags != log.Ldate|log.Ltime|log.Lmicroseconds {
		t.Errorf("flags must be %v, not: %v", log.Ldate|log.Ltime|log.Lmicroseconds, Logger.flags)
	}
}

func TestSetLevels(t *testing.T) {
	defaultLevelsList := []string{"DEBUG", "INFO", "WARN", "ERROR"}

	if !reflect.DeepEqual(defaultLevelsList, Logger.levels.List) {
		t.Errorf("default level must be: %v, not: %v", defaultLevelsList, Logger.levels.List)
	}

	SetLevels([]string{})
	if len(Logger.levels.List) != 0 {
		t.Errorf("level must empty, not: %v", Logger.levels.List)
	}
}

func TestSetMinLevel(t *testing.T) {
	defaultMinLevel := "INFO"

	if defaultMinLevel != Logger.levels.Min {
		t.Errorf("default min level must be: %v, not: %v", defaultMinLevel, Logger.levels.List)
	}

	SetMinLevel("DEBUG")
	if Logger.levels.Min != "DEBUG" {
		t.Errorf("level must %s, not: %v", defaultMinLevel, Logger.levels.Min)
	}
}

func TestSetOutput(t *testing.T) {
	buf := new(bytes.Buffer)

	if Logger.out != os.Stderr {
		t.Error("default output must be os.Stderr")
	}

	SetOutput(buf)
	if Logger.out != buf {
		t.Error("output must be buf")
	}
}

func readFromBuffer(buf *bytes.Buffer) string {
	readBuf, _ := ioutil.ReadAll(buf)
	line := strings.Replace(string(readBuf), "\n", "", 1)
	return string(line)
}
