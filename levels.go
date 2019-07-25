package logg

import (
	"bytes"
	"sync"
)

// LevelsManager - managing log levels and checking if data can be printed or not.
type LevelsManager struct {
	List []string // list of all levels
	Min  string   // minimal log level

	bad  map[string]struct{}
	once sync.Once
}

// SetLevels - set the levels of logs.
func SetLevels(list []string) {
	Logger.mu.Lock()
	Logger.levels.List = list
	Logger.mu.Unlock()
}

// SetMinLevel - set the minimum levels of logs.
func SetMinLevel(minLevel string) {
	Logger.mu.Lock()
	Logger.levels.Min = minLevel
	Logger.mu.Unlock()
}

func (lm *LevelsManager) init() {
	badLevels := make(map[string]struct{})
	for _, level := range lm.List {
		if level == lm.Min {
			break
		}
		badLevels[level] = struct{}{}
	}
	lm.bad = badLevels
}

func (lm *LevelsManager) check(level string) bool {
	lm.once.Do(lm.init)
	_, ok := lm.bad[level]
	return !ok
}

func (lm *LevelsManager) define(b []byte) string {
	var level string
	x := bytes.IndexByte(b, '[')
	if x >= 0 {
		y := bytes.IndexByte(b[x:], ']')
		if y >= 0 {
			level = string(b[x+1 : x+y])
		}
	}

	if level == "" {
		for _, l := range lm.List {
			if bytes.Contains(b, []byte(l)) {
				level = l
			}
		}
	}

	return level
}
