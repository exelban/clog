package logg

import (
	"fmt"
	"testing"
)

func TestColorsManager_CustomColor(t *testing.T) {
	cm := ColorsManager{
		list: make(map[string]string),
	}

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
			cm.list = make(map[string]string)
			cm.CustomColor(tc.prefix, tc.color)

			if cm.list[tc.prefix] == "" {
				t.Errorf("color not find: expected %v, received: %v", tc.prefix, cm.list[tc.prefix])
			}
		})
	}
}

func TestColors_define(t *testing.T) {
	c := colors{}
	cm := ColorsManager{
		list: map[string]string{
			"HIDDEN":    c.Black(),
			"ERROR":     c.Red(),
			"INFO":      c.Yellow(),
			"WARN":      c.Green(),
			"DEBUG":     c.Cyan(),
			"PANIC":     c.Blue(),
			"OWN":       c.Magenta(),
			"TEST":      c.White(),
			"HiBlack":   c.HiBlack(),
			"HiRed":     c.HiRed(),
			"HiBlue":    c.HiBlue(),
			"HiMagenta": c.HiMagenta(),
			"HiWhite":   c.HiWhite(),
			"HiGreen":   c.HiGreen(),
			"HiYellow":  c.HiYellow(),
			"HiCyan":    c.HiCyan(),
		},
	}

	tests := map[string]struct {
		data  []byte
		level string
		color int
	}{
		"HIDDEN": {
			data:  []byte("[HIDDEN] test"),
			level: "HIDDEN",
			color: Black,
		},
		"ERROR": {
			data:  []byte("[ERROR] test"),
			level: "ERROR",
			color: Red,
		},
		"INFO": {
			data:  []byte("[INFO] test"),
			level: "INFO",
			color: Yellow,
		},
		"WARN": {
			data:  []byte("[WARN] test"),
			level: "WARN",
			color: Green,
		},
		"DEBUG": {
			data:  []byte("[DEBUG] test"),
			level: "DEBUG",
			color: Cyan,
		},
		"PANIC": {
			data:  []byte("[PANIC] test"),
			level: "PANIC",
			color: Blue,
		},
		"OWN": {
			data:  []byte("[OWN] test"),
			level: "OWN",
			color: Magenta,
		},
		"TEST": {
			data:  []byte("[TEST] test"),
			level: "TEST",
			color: White,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			color := cm.define(tc.level)
			expectedColor := fmt.Sprintf("3%d", tc.color)

			if color != expectedColor {
				t.Errorf("wrong color: expected %v, received: %v", expectedColor, color)
			}
		})
	}
}
