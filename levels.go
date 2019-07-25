package logg

import (
	"bytes"
	"sync"
)

// LevelFilter - keeps levels list, minimum level. Most part of this was borrowed from logutils library from hashicorp.
// https://github.com/hashicorp/logutils
type LevelFilter struct {
	Levels []string
	MinLevel string

	badLevels map[string]struct{}
	once sync.Once
}

// Check - check if string have a minimum level to be write to output.
func (lf *LevelFilter) Check(b []byte) bool {
	lf.once.Do(lf.init)

	var level string
	x := bytes.IndexByte(b, '[')
	if x >= 0 {
		y := bytes.IndexByte(b[x:], ']')
		if y >= 0 {
			level = string(b[x+1 : x+y])
		}
	}

	_, ok := lf.badLevels[level]
	return !ok
}

// init - init function allows to create list of levels which must be missed.
func (lf *LevelFilter) init() {
	badLevels := make(map[string]struct{})
	for _, level := range lf.Levels {
		if level == lf.MinLevel {
			break
		}
		badLevels[level] = struct{}{}
	}
	lf.badLevels = badLevels
}