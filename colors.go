package logg

import (
	"bytes"
	"fmt"
	"sync"
)

const escape = "\x1b"
const textBase = 30
const backgroundBase = 40

// Base colors
const (
	Black int = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// High intensity colors
const (
	HiBlack int = iota + 60
	HiRed
	HiGreen
	HiYellow
	HiBlue
	HiMagenta
	HiCyan
	HiWhite
)

// Base attributes
const (
	Reset int = iota
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	ReverseVideo
	Concealed
	CrossedOut
)

type ColorsManager struct {
	list map[string]string
	mu   sync.Mutex
}

func (cm *ColorsManager) define(m *message) string {
	var color string

	for p, c := range cm.list {
		if bytes.Contains(m.data, []byte(p)) {
			color = c
		}
	}

	return color
}

// Custom - allow to set custom colors for prefix.
// Accept parameters in next configuration: [textColor, backgroundColor, style].
func (cm *ColorsManager) CustomColor(prefix string, v ...interface{}) {
	if len(v) == 0 {
		panic(fmt.Sprintf("logg: missed configuration for %s", prefix)) // TODO: remove panic
	}

	switch v[0].(type) {
	case int:
	default:
		panic(fmt.Sprintf("logg: wrong configuration for %s (%v)", prefix, v)) // TODO: remove panic
	}

	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.list[prefix] = generate(v...)
}

// Base colors
type Colors interface {
	Black() string
	Red() string
	Green() string
	Yellow() string
	Blue() string
	Magenta() string
	Cyan() string
	White() string

	HiBlack() string
	HiRed() string
	HiGreen() string
	HiYellow() string
	HiBlue() string
	HiMagenta() string
	HiCyan() string
	HiWhite() string
}

type colors struct{}

// Black text color.
func (c *colors) Black() string { return generate(Black) }

// Red text color.
func (c *colors) Red() string { return generate(Red) }

// Green text color.
func (c *colors) Green() string { return generate(Green) }

// Yellow text color.
func (c *colors) Yellow() string { return generate(Yellow) }

// Blue text color.
func (c *colors) Blue() string { return generate(Blue) }

// Magenta text color.
func (c *colors) Magenta() string { return generate(Magenta) }

// Cyan text color.
func (c *colors) Cyan() string { return generate(Cyan) }

// White text color.
func (c *colors) White() string { return generate(White) }

// Black high intense color.
func (c *colors) HiBlack() string { return generate(HiBlack) }

// Red high intense color.
func (c *colors) HiRed() string { return generate(HiRed) }

// Green high intense color.
func (c *colors) HiGreen() string { return generate(HiGreen) }

// Yellow high intense color.
func (c *colors) HiYellow() string { return generate(HiYellow) }

// Blue high intense color.
func (c *colors) HiBlue() string { return generate(HiBlue) }

// Magenta high intense color.
func (c *colors) HiMagenta() string { return generate(HiMagenta) }

// Cyan high intense color.
func (c *colors) HiCyan() string { return generate(HiCyan) }

// White high intense color.
func (c *colors) HiWhite() string { return generate(HiWhite) }

func generate(v ...interface{}) string {
	var textBase = 30
	var backgroundBase = 40

	var color string

	switch len(v) {
	case 1:
		text := textBase + v[0].(int)
		color = fmt.Sprintf("%d;", text)
	case 2:
		text := textBase + v[0].(int)
		background := backgroundBase + v[1].(int)
		color = fmt.Sprintf("%d;%d;", text, background)
	case 3:
		text := textBase + v[0].(int)
		background := backgroundBase + v[1].(int)
		style := v[2].(int)
		color = fmt.Sprintf("%d;%d;%d;", style, text, background)
	}

	return color
}
