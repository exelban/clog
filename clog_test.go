package clog

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"testing"
)

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
		text string
		color int
	}{
		{
			text: "black text",
			color: 30,
		},
		{
			text: "red text",
			color: Red,
		},
		{
			text: "green text",
			color: Green,
		},
		{
			text: "yellow text",
			color: Yellow,
		},
		{
			text: "blue text",
			color: Blue,
		},
		{
			text: "magenta text",
			color: Magenta,
		},
		{
			text: "cyan text",
			color: Cyan,
		},
		{
			text: "white text",
			color: White,
		},
	}

	for _, c := range testList {
		color := generate(c.color)

		writer.set(color)
		writer.out.Write([]byte(c.text))
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
		text string
		color func(clog Colors) string
	}{
		{
			prefix: "[HIDDEN]",
			text: "[HIDDEN] black text",
			color: Colors.Black,
		},
		{
			prefix: "[ERROR]",
			text: "[ERROR] error text",
			color: Colors.Red,
		},
		{
			prefix: "[INFO]",
			text: "[INFO] info text",
			color: Colors.Yellow,
		},
		{
			prefix: "[WARN]",
			text: "[WARN] warn text",
			color: Colors.Green,
		},
		{
			prefix: "[DEBUG]",
			text: "[DEBUG] debug text",
			color: Colors.Cyan,
		},
		{
			prefix: "[PANIC]",
			text: "[PANIC] panic text",
			color: Colors.Blue,
		},
		{
			prefix: "[OWN]",
			text: "[OWN] own text",
			color: Colors.Magenta,
		},
		{
			prefix: "[TEST]",
			text: "[TEST] white text",
			color: Colors.White,
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
		text string
		color func(clog Colors) string
	}{
		{
			prefix: "[HIDDEN]",
			text: "[HIDDEN] black text",
			color: Colors.HiBlack,
		},
		{
			prefix: "[ERROR]",
			text: "[ERROR] error text",
			color: Colors.HiRed,
		},
		{
			prefix: "[INFO]",
			text: "[INFO] info text",
			color: Colors.HiYellow,
		},
		{
			prefix: "[WARN]",
			text: "[WARN] warn text",
			color: Colors.HiGreen,
		},
		{
			prefix: "[DEBUG]",
			text: "[DEBUG] debug text",
			color: Colors.HiCyan,
		},
		{
			prefix: "[PANIC]",
			text: "[PANIC] panic text",
			color: Colors.HiBlue,
		},
		{
			prefix: "[OWN]",
			text: "[OWN] own text",
			color: Colors.HiMagenta,
		},
		{
			prefix: "[TEST]",
			text: "[TEST] white text",
			color: Colors.HiWhite,
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
		text string
		color int
	}{
		{
			prefix: "[HIDDEN]",
			text: "[HIDDEN] black text",
			color: Black,
		},
		{
			prefix: "[ERROR]",
			text: "[ERROR] error text",
			color: Red,
		},
		{
			prefix: "[INFO]",
			text: "[INFO] info text",
			color: Yellow,
		},
		{
			prefix: "[WARN]",
			text: "[WARN] warn text",
			color: Green,
		},
		{
			prefix: "[DEBUG]",
			text: "[DEBUG] debug text",
			color: Cyan,
		},
		{
			prefix: "[PANIC]",
			text: "[PANIC] panic text",
			color: Blue,
		},
		{
			prefix: "[OWN]",
			text: "[OWN] own text",
			color: Magenta,
		},
		{
			prefix: "[TEST]",
			text: "[TEST] white text",
			color: White,
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

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic on wrong parameters in Custom()")
		}
	}()

	writer.Custom("[TEST]")

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic on wrong parameters in Custom()")
		}
	}()
	writer.Custom("[TEST]", "1")
}

func TestWriter_Style(t *testing.T) {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	writer := Install()
	buf := new(bytes.Buffer)
	writer.out = buf

	testList := []struct {
		prefix string
		text string
		style int
	}{
		{
			prefix: "[BOLD]",
			text: "[BOLD] bold text",
			style: Bold,
		},
		{
			prefix: "[FAINT]",
			text: "[FAINT] faint text",
			style: Faint,
		},
		{
			prefix: "[ITALIC]",
			text: "[ITALIC] italic text",
			style: Italic,
		},
		{
			prefix: "[UNDERLINE]",
			text: "[UNDERLINE] underline text",
			style: Underline,
		},
		{
			prefix: "[BLINKSLOW]",
			text: "[BLINKSLOW] blink slow text",
			style: BlinkSlow,
		},
		{
			prefix: "[BLINKRAPID]",
			text: "[BLINKRAPID] blink rapid text",
			style: BlinkRapid,
		},
		{
			prefix: "[REVERSEVIDEO]",
			text: "[REVERSEVIDEO] reverse video text",
			style: ReverseVideo,
		},
		{
			prefix: "[CANCEALED]",
			text: "[CANCEALED] concealed text",
			style: Concealed,
		},
		{
			prefix: "[CROSSEDOUT]",
			text: "[CROSSEDOUT] crossed out text",
			style: CrossedOut,
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
	writer.Uninstall()

	log.Print("visibility test")
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

func readFromBuffer(buf *bytes.Buffer) string {
	readBuf, _ := ioutil.ReadAll(buf)
	line := strings.Replace(string(readBuf), "\n", "", 1)
	return string(line)
}