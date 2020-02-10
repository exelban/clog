package logg

import (
	"fmt"
	"sync"
)

const escapeClose = "\x1b[0m"

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

type Colors struct {
	list map[string][]byte
	mu   sync.Mutex
}

func (c *Colors) init() {
	for k, color := range c.list {
		if color[0] != 27 && color[len(color)-1] != 'm' {
			col := []byte{27, 91}
			col = append(col, color...)
			col = append(col, 'm')
			c.list[k] = col
		}
	}
}

// define the log color based on log level.
func (c *Colors) define(m *message) {
	m.color = c.list[string(m.level)]
}

// CustomColor - allow to set custom colors for prefix.
func (c *Colors) CustomColor(prefix string, v int) {
	c.mu.Lock()
	c.list[prefix] = generate(v)
	c.mu.Unlock()
}

func generate(v int) []byte {
	return []byte(fmt.Sprintf("%d", 30+v))
}
