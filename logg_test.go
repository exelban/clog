package logg

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"testing"
)

var messages = [][]byte{
	[]byte("[TRACE] foo"),
	[]byte("[DEBUG] foo"),
	[]byte("[INFO] foo"),
	[]byte("[WARN] foo"),
	[]byte("[ERROR] foo"),
}

func TestInstall(t *testing.T) {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	writer := Install()
	buf := new(bytes.Buffer)
	writer.out = buf

	message := "test install"
	log.Print(message)

	line := readFromBuffer(buf)
	colored := fmt.Sprintf("%s", message)

	//fmt.Println("Has: ", []byte(line))
	//fmt.Println("Want:", []byte(colored))

	if line != colored {
		t.Errorf("Expecting %s, got '%s'\n", colored, line)
	}
}

func TestInstall2(t *testing.T) {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	writer := Install(Black, Green)
	buf := new(bytes.Buffer)
	writer.out = buf
	color := generate(Black, Green)

	message := "test install"
	log.Print(message)

	line := readFromBuffer(buf)
	colored := fmt.Sprintf("\x1b[%sm%s\x1b[0m", color, message)

	if line != colored {
		t.Errorf("Expecting %s, got '%s'\n", colored, line)
	}
}

func TestColor(t *testing.T) {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	writer := Install()
	buf := new(bytes.Buffer)
	writer.out = buf

	testList := []struct {
		text  string
		color int
	}{
		{
			text:  "black text",
			color: 30,
		},
		{
			text:  "red text",
			color: Red,
		},
		{
			text:  "green text",
			color: Green,
		},
		{
			text:  "yellow text",
			color: Yellow,
		},
		{
			text:  "blue text",
			color: Blue,
		},
		{
			text:  "magenta text",
			color: Magenta,
		},
		{
			text:  "cyan text",
			color: Cyan,
		},
		{
			text:  "white text",
			color: White,
		},
	}

	for _, c := range testList {
		color := generate(c.color)

		writer.set(color)
		n, err := writer.out.Write([]byte(c.text))
		if err != nil {
			t.Errorf("Not expected error '%s'\n", err)
		}
		if n != len([]byte(c.text)) {
			t.Errorf("Writed number must be the same length as text\n")
		}
		writer.unset()

		line := readFromBuffer(buf)
		colored := fmt.Sprintf("\x1b[%sm%s\x1b[0m", color, c.text)

		if line != colored {
			t.Errorf("Expecting '%s', got '%s'\n", colored, line)
		}
	}
}

func TestWriter_Prefix(t *testing.T) {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	writer := Install()
	buf := new(bytes.Buffer)
	writer.out = buf

	testList := []struct {
		prefix string
		text   string
		color  func(logg Colors) string
	}{
		{
			prefix: "[HIDDEN]",
			text:   "[HIDDEN] black text",
			color:  Colors.Black,
		},
		{
			prefix: "[ERROR]",
			text:   "[ERROR] error text",
			color:  Colors.Red,
		},
		{
			prefix: "[INFO]",
			text:   "[INFO] info text",
			color:  Colors.Yellow,
		},
		{
			prefix: "[WARN]",
			text:   "[WARN] warn text",
			color:  Colors.Green,
		},
		{
			prefix: "[DEBUG]",
			text:   "[DEBUG] debug text",
			color:  Colors.Cyan,
		},
		{
			prefix: "[PANIC]",
			text:   "[PANIC] panic text",
			color:  Colors.Blue,
		},
		{
			prefix: "[OWN]",
			text:   "[OWN] own text",
			color:  Colors.Magenta,
		},
		{
			prefix: "[TEST]",
			text:   "[TEST] white text",
			color:  Colors.White,
		},
	}

	for _, c := range testList {
		writer.Prefix(c.prefix, c.color)
		color := c.color(&colors{})

		log.Print(c.text)

		line := readFromBuffer(buf)
		colored := fmt.Sprintf("\x1b[%sm%s\x1b[0m", color, c.text)

		fmt.Println(line)
		if line != colored {
			t.Errorf("Expecting %s, got '%s'\n", colored, line)
		}
	}
}

func TestWriter_Prefix2(t *testing.T) {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	writer := Install()
	buf := new(bytes.Buffer)
	writer.out = buf

	testList := []struct {
		prefix string
		text   string
		color  func(logg Colors) string
	}{
		{
			prefix: "[HIDDEN]",
			text:   "[HIDDEN] black text",
			color:  Colors.HiBlack,
		},
		{
			prefix: "[ERROR]",
			text:   "[ERROR] error text",
			color:  Colors.HiRed,
		},
		{
			prefix: "[INFO]",
			text:   "[INFO] info text",
			color:  Colors.HiYellow,
		},
		{
			prefix: "[WARN]",
			text:   "[WARN] warn text",
			color:  Colors.HiGreen,
		},
		{
			prefix: "[DEBUG]",
			text:   "[DEBUG] debug text",
			color:  Colors.HiCyan,
		},
		{
			prefix: "[PANIC]",
			text:   "[PANIC] panic text",
			color:  Colors.HiBlue,
		},
		{
			prefix: "[OWN]",
			text:   "[OWN] own text",
			color:  Colors.HiMagenta,
		},
		{
			prefix: "[TEST]",
			text:   "[TEST] white text",
			color:  Colors.HiWhite,
		},
	}

	for _, c := range testList {
		writer.Prefix(c.prefix, c.color)
		color := c.color(&colors{})

		log.Print(c.text)

		line := readFromBuffer(buf)
		colored := fmt.Sprintf("\x1b[%sm%s\x1b[0m", color, c.text)

		fmt.Println(line)
		if line != colored {
			t.Errorf("Expecting %s, got '%s'\n", colored, line)
		}
	}
}

func TestWriter_Custom(t *testing.T) {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	writer := Install()
	buf := new(bytes.Buffer)
	writer.out = buf

	testList := []struct {
		prefix string
		text   string
		color  int
	}{
		{
			prefix: "[HIDDEN]",
			text:   "[HIDDEN] black text",
			color:  Black,
		},
		{
			prefix: "[ERROR]",
			text:   "[ERROR] error text",
			color:  Red,
		},
		{
			prefix: "[INFO]",
			text:   "[INFO] info text",
			color:  Yellow,
		},
		{
			prefix: "[WARN]",
			text:   "[WARN] warn text",
			color:  Green,
		},
		{
			prefix: "[DEBUG]",
			text:   "[DEBUG] debug text",
			color:  Cyan,
		},
		{
			prefix: "[PANIC]",
			text:   "[PANIC] panic text",
			color:  Blue,
		},
		{
			prefix: "[OWN]",
			text:   "[OWN] own text",
			color:  Magenta,
		},
		{
			prefix: "[TEST]",
			text:   "[TEST] white text",
			color:  White,
		},
	}

	for _, c := range testList {
		writer.Custom(c.prefix, c.color)
		color := generate(c.color)

		log.Print(c.text)

		line := readFromBuffer(buf)
		colored := fmt.Sprintf("\x1b[%sm%s\x1b[0m", color, c.text)

		fmt.Println(line)
		if line != colored {
			t.Errorf("Expecting %s, got '%s'\n", colored, line)
		}
	}
}

func TestWriter_Custom2(t *testing.T) {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	writer := Install()
	buf := new(bytes.Buffer)
	writer.out = buf

	writer.Custom("[TEST]", Red, Blue)
	color := generate(Red, Blue)

	message := "[TEST] test custom with two parameters"
	log.Print(message)

	line := readFromBuffer(buf)
	colored := fmt.Sprintf("\x1b[%sm%s\x1b[0m", color, message)

	writer.Custom("[TEST]", Red, Blue, Green)
	color = generate(Red, Blue, Green)

	message = "[TEST] test custom with all parameters"
	log.Print(message)

	line = readFromBuffer(buf)
	colored = fmt.Sprintf("\x1b[%sm%s\x1b[0m", color, message)

	if line != colored {
		t.Errorf("Expecting %s, got '%s'\n", colored, line)
	}

	prefix := "[TEST]"
	defer func() {
		r := recover()

		fmt.Println("WTF 1")
		if r != fmt.Sprintf("logg: missed configuration for %s", prefix) {
			t.Error("Must throw missed configuration")
		}

		if r == nil {
			t.Errorf("The code did not panic on wrong parameters in Custom()")
		}
	}()
	writer.Custom(prefix)
}

func TestWriter_Custom3(t *testing.T) {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	writer := Install()
	buf := new(bytes.Buffer)
	writer.out = buf

	prefix := "[TEST]"
	defer func() {
		r := recover()

		if !strings.Contains(fmt.Sprintf("%v", r), fmt.Sprintf("logg: wrong configuration for %s", prefix)) {
			t.Error("Must throw wrong configuration")
		}

		if r == nil {
			t.Errorf("The code did not panic on wrong parameters in Custom()")
		}
	}()

	writer.Custom(prefix, "1")
}

func TestWriter_Style(t *testing.T) {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	writer := Install()
	buf := new(bytes.Buffer)
	writer.out = buf

	testList := []struct {
		prefix string
		text   string
		style  int
	}{
		{
			prefix: "[BOLD]",
			text:   "[BOLD] bold text",
			style:  Bold,
		},
		{
			prefix: "[FAINT]",
			text:   "[FAINT] faint text",
			style:  Faint,
		},
		{
			prefix: "[ITALIC]",
			text:   "[ITALIC] italic text",
			style:  Italic,
		},
		{
			prefix: "[UNDERLINE]",
			text:   "[UNDERLINE] underline text",
			style:  Underline,
		},
		{
			prefix: "[BLINKSLOW]",
			text:   "[BLINKSLOW] blink slow text",
			style:  BlinkSlow,
		},
		{
			prefix: "[BLINKRAPID]",
			text:   "[BLINKRAPID] blink rapid text",
			style:  BlinkRapid,
		},
		{
			prefix: "[REVERSEVIDEO]",
			text:   "[REVERSEVIDEO] reverse video text",
			style:  ReverseVideo,
		},
		{
			prefix: "[CANCEALED]",
			text:   "[CANCEALED] concealed text",
			style:  Concealed,
		},
		{
			prefix: "[CROSSEDOUT]",
			text:   "[CROSSEDOUT] crossed out text",
			style:  CrossedOut,
		},
	}

	for _, c := range testList {
		writer.Custom(c.prefix, HiCyan, Black, c.style)
		color := generate(HiCyan, Black, c.style)

		log.Print(c.text)

		line := readFromBuffer(buf)
		colored := fmt.Sprintf("\x1b[%sm%s\x1b[0m", color, c.text)

		fmt.Println(line)
		if line != colored {
			t.Errorf("Expecting %s, got '%s'\n", colored, line)
		}
	}
}

func TestWriter_Uninstall(t *testing.T) {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	writer := Install()
	buf := new(bytes.Buffer)
	writer.out = buf
	writer.Uninstall()

	log.Print("visibility test")

	if buf.Len() != 0 {
		t.Error("Buffer must be empty")
	}
}

func TestWriter_SetOutput(t *testing.T) {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	writer := Install(Black, Green)
	buf := new(bytes.Buffer)
	writer.SetOutput(buf)
	color := generate(Black, Green)

	message := "test install"
	log.Print(message)

	line := readFromBuffer(buf)
	colored := fmt.Sprintf("\x1b[%sm%s\x1b[0m", color, message)

	if line != colored {
		t.Errorf("Expecting %s, got '%s'\n", colored, line)
	}
}

func TestWriter_SetFilter(t *testing.T) {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	writer := Install()
	buf := new(bytes.Buffer)
	writer.out = buf

	levelFilter := &LevelFilter{
		Levels:   []string{"TESTTTTT1", "TESTTTTT2", "TESTTTTT3", "TESTTTTT4"},
		MinLevel: "TESTTTTT2",
	}
	writer.SetFilters(levelFilter)

	if writer.filters != levelFilter {
		t.Errorf("Level fiters are not initialized properly")
	}

	log.Print("[TESTTTTT1] foo")
	log.Print("[TESTTTTT2] baz")
	log.Print("[TESTTTTT3] buzz")
	log.Print("[TESTTTTT4] bar")

	line := readFromBuffer(buf)
	if strings.Contains(line, "TESTTTTT1") ||
		!strings.Contains(line, "TESTTTTT2") ||
		!strings.Contains(line, "TESTTTTT3") ||
		!strings.Contains(line, "TESTTTTT4") {
		t.Errorf("Must be only 3 levels printed, received: %v", line)
	}
}

func TestWriter_SetMinLevel(t *testing.T) {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	writer := Install()
	buf := new(bytes.Buffer)
	writer.out = buf

	levelFilter := &LevelFilter{
		Levels:   []string{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: "INFO",
	}
	writer.SetFilters(levelFilter)

	buf = new(bytes.Buffer)
	writer.out = buf
	writer.SetMinLevel("ERROR")

	if writer.filters.MinLevel != "ERROR" {
		t.Errorf("Minimum level not set properly")
	}

	log.Print("[DEBUG] foo")
	log.Print("[INFO] baz")
	log.Print("[WARN] buzz")
	log.Print("[ERROR] bar")

	line := readFromBuffer(buf)
	if strings.Contains(line, "DEBUG") ||
		strings.Contains(line, "INFO") ||
		strings.Contains(line, "WARN") ||
		!strings.Contains(line, "ERROR") {
		t.Errorf("Must be only ERROR levels printed, received: %v", line)
	}
}

func BenchmarkDiscard(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = ioutil.Discard.Write(messages[i%len(messages)])
	}
}

func BenchmarkLoggWrite(b *testing.B) {
	b.ReportAllocs()

	writer := Install()
	writer.out = ioutil.Discard

	for i := 0; i < b.N; i++ {
		_, _ = writer.Write(messages[i%len(messages)])
	}
}

func BenchmarkLogg(b *testing.B) {
	b.ReportAllocs()

	writer := Install()
	writer.out = ioutil.Discard

	for i := 0; i < b.N; i++ {
		log.Print(messages[i%len(messages)])
	}
}

func BenchmarkLog(b *testing.B) {
	b.ReportAllocs()

	log.SetOutput(ioutil.Discard)

	for i := 0; i < b.N; i++ {
		log.Print(messages[i%len(messages)])
	}
}

func readFromBuffer(buf *bytes.Buffer) string {
	readBuf, _ := ioutil.ReadAll(buf)
	line := strings.Replace(string(readBuf), "\n", "", 1)
	return string(line)
}
