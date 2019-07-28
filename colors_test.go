package logg

import (
	"fmt"
	"strings"
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

func TestColorsManager_CustomColor_panics(t *testing.T) {
	cm := ColorsManager{
		list: make(map[string]string),
	}

	prefix := "[TEST]"
	defer func() {
		r := recover()

		if r != fmt.Sprintf("logg: missed configuration for %s", prefix) {
			t.Error("Must throw missed configuration")
		}

		if r == nil {
			t.Errorf("The code did not panic on wrong parameters in Custom()")
		}
	}()
	cm.CustomColor(prefix)
}

func TestColorsManager_CustomColor_panics2(t *testing.T) {
	cm := ColorsManager{
		list: make(map[string]string),
	}

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

	cm.CustomColor(prefix, "1")
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
		},
	}

	tests := map[string]struct {
		data  []byte
		color int
	}{
		"HIDDEN": {
			data:  []byte("[HIDDEN] test"),
			color: Black,
		},
		"ERROR": {
			data:  []byte("[ERROR] test"),
			color: Red,
		},
		"INFO": {
			data:  []byte("[INFO] test"),
			color: Yellow,
		},
		"WARN": {
			data:  []byte("[WARN] test"),
			color: Green,
		},
		"DEBUG": {
			data:  []byte("[DEBUG] test"),
			color: Cyan,
		},
		"PANIC": {
			data:  []byte("[PANIC] test"),
			color: Blue,
		},
		"OWN": {
			data:  []byte("[OWN] test"),
			color: Magenta,
		},
		"TEST": {
			data:  []byte("[TEST] test"),
			color: White,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			color := cm.define(&tc.data)
			expectedColor := fmt.Sprintf("3%d;", tc.color)

			if color != expectedColor {
				t.Errorf("wrong color: expected %v, received: %v", expectedColor, color)
			}
		})
	}
}

func TestColors_background(t *testing.T) {
	tests := map[string]int{
		"HiBlack":   HiBlack,
		"HiRed":     HiRed,
		"HiGreen":   HiGreen,
		"HiYellow":  HiYellow,
		"HiBlue":    HiBlue,
		"HiMagenta": HiMagenta,
		"HiCyan":    HiCyan,
		"HiWhite":   HiWhite,
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			background := generate(HiCyan, tc)
			expectedBackground := fmt.Sprintf("96;%d;", backgroundBase+tc)

			if background != expectedBackground {
				t.Errorf("wrong background: expected %v, received: %v", expectedBackground, background)
			}
		})
	}
}

func TestColors_styles(t *testing.T) {
	tests := map[string]int{
		"BOLD":         Bold,
		"FAINT":        Faint,
		"ITALIC":       Italic,
		"UNDERLINE":    Underline,
		"BLINKSLOW":    BlinkSlow,
		"BLINKRAPID":   BlinkRapid,
		"REVERSEVIDEO": ReverseVideo,
		"CANCEALED":    Concealed,
		"CROSSEDOUT":   CrossedOut,
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			style := generate(HiCyan, Black, tc)
			expectedStyle := fmt.Sprintf("%d;96;40;", tc)

			if style != expectedStyle {
				t.Errorf("wrong style: expected %v, received: %v", expectedStyle, style)
			}
		})
	}
}
