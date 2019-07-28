package logg

import (
	"reflect"
	"testing"
)

func TestLevels_init(t *testing.T) {
	tests := map[string]struct {
		lm              LevelsManager
		expectedBadList map[string]bool
	}{
		"test1": {
			lm: LevelsManager{
				List: []string{""},
				Min:  "",
			},
			expectedBadList: map[string]bool{},
		},
		"test2": {
			lm: LevelsManager{
				List: []string{"DEBUG", "INFO", "WARN", "ERROR"},
				Min:  "",
			},
			expectedBadList: map[string]bool{
				"DEBUG": true,
				"INFO":  true,
				"WARN":  true,
				"ERROR": true,
			},
		},
		"test3": {
			lm: LevelsManager{
				List: []string{"DEBUG", "INFO", "WARN", "ERROR"},
				Min:  "INFO",
			},
			expectedBadList: map[string]bool{
				"DEBUG": true,
			},
		},
		"test4": {
			lm: LevelsManager{
				List: []string{"DEBUG", "INFO", "WARN", "ERROR"},
				Min:  "ERROR",
			},
			expectedBadList: map[string]bool{
				"DEBUG": true,
				"INFO":  true,
				"WARN":  true,
			},
		},
		"test5": {
			lm: LevelsManager{
				List: []string{"DEBUG", "INFO", "WARN", "ERROR"},
				Min:  "TEST",
			},
			expectedBadList: map[string]bool{
				"DEBUG": true,
				"INFO":  true,
				"WARN":  true,
				"ERROR": true,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.lm.init()

			if !reflect.DeepEqual(tc.lm.bad, tc.expectedBadList) {
				t.Errorf("wrong bad list: expected %v, received: %v", tc.expectedBadList, tc.lm.bad)
			}
		})
	}
}

func TestLevels_check(t *testing.T) {
	tests := map[string]struct {
		lm            LevelsManager
		testLevel     string
		expectedCheck bool
	}{
		"test1": {
			lm: LevelsManager{
				List: []string{""},
				Min:  "",
			},
			testLevel:     "",
			expectedCheck: true,
		},
		"test2": {
			lm: LevelsManager{
				List: []string{"DEBUG", "INFO", "WARN", "ERROR"},
				Min:  "",
			},
			testLevel:     "",
			expectedCheck: true,
		},
		"test3": {
			lm: LevelsManager{
				List: []string{"DEBUG", "INFO", "WARN", "ERROR"},
				Min:  "",
			},
			testLevel:     "INFO",
			expectedCheck: false,
		},
		"test4": {
			lm: LevelsManager{
				List: []string{"DEBUG", "INFO", "WARN", "ERROR"},
				Min:  "INFO",
			},
			testLevel:     "INFO",
			expectedCheck: true,
		},
		"test5": {
			lm: LevelsManager{
				List: []string{"DEBUG", "INFO", "WARN", "ERROR"},
				Min:  "INFO",
			},
			testLevel:     "DEBUG",
			expectedCheck: false,
		},
		"test6": {
			lm: LevelsManager{
				List: []string{"DEBUG", "INFO", "WARN", "ERROR"},
				Min:  "INFO",
			},
			testLevel:     "ERROR",
			expectedCheck: true,
		},
		"test7": {
			lm: LevelsManager{
				List: []string{"DEBUG", "INFO", "WARN", "ERROR"},
				Min:  "ERROR",
			},
			testLevel:     "",
			expectedCheck: true,
		},
		"test8": {
			lm: LevelsManager{
				List: []string{"DEBUG", "INFO", "WARN", "ERROR"},
				Min:  "ERROR",
			},
			testLevel:     "WARN",
			expectedCheck: false,
		},
		"test9": {
			lm: LevelsManager{
				List: []string{"DEBUG", "INFO", "WARN", "ERROR"},
				Min:  "ERROR",
			},
			testLevel:     "ERROR",
			expectedCheck: true,
		},
		"test10": {
			lm: LevelsManager{
				List: []string{"DEBUG", "INFO", "WARN", "ERROR"},
				Min:  "TEST",
			},
			testLevel:     "",
			expectedCheck: true,
		},
		"test11": {
			lm: LevelsManager{
				List: []string{"DEBUG", "INFO", "WARN", "ERROR"},
				Min:  "TEST",
			},
			testLevel:     "INFO",
			expectedCheck: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			val := tc.lm.check(tc.testLevel)

			if val != tc.expectedCheck {
				t.Errorf("wrong check value: expected %v, received: %v", tc.expectedCheck, val)
			}
		})
	}
}

func TestLevels_define(t *testing.T) {
	lm := LevelsManager{
		List: []string{"DEBUG", "INFO", "WARN", "ERROR"},
	}

	tests := map[string]struct {
		data          []byte
		expectedLevel string
	}{
		"test1": {
			data:          []byte(""),
			expectedLevel: "",
		},
		"test2": {
			data:          []byte("info"),
			expectedLevel: "",
		},
		"test3": {
			data:          []byte("INFO"),
			expectedLevel: "INFO",
		},
		"test4": {
			data:          []byte("ERROr"),
			expectedLevel: "",
		},
		"test5": {
			data:          []byte("ERRO"),
			expectedLevel: "",
		},
		"test6": {
			data:          []byte("[INFO]"),
			expectedLevel: "INFO",
		},
		"test7": {
			data:          []byte("TEST"),
			expectedLevel: "",
		},
		"test8": {
			data:          []byte("[TEST]"),
			expectedLevel: "TEST",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			lvl := lm.define(tc.data)

			if lvl != tc.expectedLevel {
				t.Errorf("wrong level: expected %v, received: %v", tc.expectedLevel, lvl)
			}
		})
	}
}
