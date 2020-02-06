package logg

import (
	"bytes"
	"fmt"
	"sync"
)

// LevelsManager - managing log levels and checking if data can be printed or not.
type LevelsManager struct {
	List []string // list of all levels
	Min  string   // minimal log level

	bad  map[string]bool
	once sync.Once
	mu   sync.Mutex
}

func (lm *LevelsManager) init() {
	badLevels := make(map[string]bool)
	for _, level := range lm.List {
		if level == lm.Min {
			break
		}
		badLevels[level] = true
	}
	lm.bad = badLevels
}

func (lm *LevelsManager) check(level string) bool {
	lm.once.Do(lm.init)
	_, ok := lm.bad[level]
	return !ok
}

func (lm *LevelsManager) define(b *[]byte) string {
	var level string

	x := bytes.IndexByte(*b, '[')
	if x >= 0 {
		y := bytes.IndexByte((*b)[x:], ']')
		if y >= 0 {
			level = string((*b)[x+1 : x+y])
			if len(level) == 4 {
				*b = bytes.Replace(*b, []byte(fmt.Sprintf("%s]", level)), []byte(fmt.Sprintf("%s] ", level)), -1)
			}
		}
	}

	if level == "" {
		lim := len(*b)
		if len(*b) >= 5 {
			lim = 5
		}

		for _, l := range lm.List {
			if bytes.Contains((*b)[:lim], []byte(l)) {
				level = l
				if len(level) == 4 {
					*b = bytes.Replace(*b, []byte(level), []byte(fmt.Sprintf("%s ", level)), -1)
				}
			}
		}
	}

	return level
}
