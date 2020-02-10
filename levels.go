package logg

import (
	"bytes"
)

// Levels - managing log levels and checking if data can be printed or not.
type Levels struct {
	List []string // list of all levels
	Min  string   // minimal log level

	bad map[string]bool
}

// Allows to generate a list of levels which are allows to log.
func (l *Levels) init() {
	badLevels := make(map[string]bool)

	for _, level := range l.List {
		if level == l.Min {
			break
		}
		badLevels[level] = true
	}

	l.bad = badLevels
}

// Define the log level and checks if log level is allowed to print (more than minimum level).
func (l *Levels) define(m *message) {
	x := bytes.IndexByte(m.data, '[')
	if x >= 0 {
		y := bytes.IndexByte(m.data[x:], ']')
		if y >= 0 {
			for i := 0; i < len(l.List); i++ {
				l := l.List[i]
				if l == string(m.data[x+1:x+y]) {
					m.level = m.data[x+1 : x+y]
					m.data = m.data[x+y+1:] // removing level from data
					m.brackets = true
					break
				}
			}
		}
	} else {
		for i := 0; i < len(l.List); i++ {
			l := l.List[i]
			if bytes.Contains(m.data, []byte(l)) {
				x := bytes.Index(m.data, []byte(l))
				m.level = []byte(l)
				m.data = m.data[x+len([]byte(l)):] // removing level from data
				break
			}
		}
	}

	_, ok := l.bad[string(m.level)]
	m.allowed = !ok
}
