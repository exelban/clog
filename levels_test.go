package logg

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestLevelFilter_Check(t *testing.T) {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	writer := Install()
	buf := new(bytes.Buffer)
	writer.out = buf

	levelFilter := &LevelFilter{
		Levels:   []string{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: "DEBUG",
	}
	writer.SetFilters(levelFilter)

	log.Print("[DEBUG] foo")
	log.Print("[INFO] baz")
	log.Print("[WARN] buzz")
	log.Print("[ERROR] bar")

	line := readFromBuffer(buf)
	if !strings.Contains(line, "DEBUG") ||
		!strings.Contains(line, "INFO") ||
		!strings.Contains(line, "WARN") ||
		!strings.Contains(line, "ERROR") {
		t.Errorf("Must be all levels printed, received: %v", line)
	}

	buf = new(bytes.Buffer)
	writer.out = buf
	writer.SetMinLevel("WARN")

	log.Print("[DEBUG] foo")
	log.Print("[INFO] baz")
	log.Print("[WARN] buzz")
	log.Print("[ERROR] bar")

	line = readFromBuffer(buf)
	if strings.Contains(line, "DEBUG") ||
		strings.Contains(line, "INFO") ||
		!strings.Contains(line, "WARN") ||
		!strings.Contains(line, "ERROR") {
		t.Errorf("Must be only WARN and ERROR levels printed, received: %v", line)
	}
}
