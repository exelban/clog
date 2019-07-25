# logg
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/exelban/logg)
[![codecov](https://codecov.io/gh/exelban/logg/branch/master/graph/badge.svg)](https://codecov.io/gh/exelban/logg)

![](https://s3.eu-central-1.amazonaws.com/serhiy/Github_repo/clog/Zrzut+ekranu+2018-10-16+o+18.52.26.png)  
Color logs for your go application.

# Installation
```bash
go get github.com/exelban/logg
```

# Usage
## Example

### Simple usage
```golang
package main

import (
	"github.com/exelban/logg"
	"log"
)

func main () {
	logg.Install()
	
	log.Print("[ERROR] error text")
}
```

### Custom level
```golang
package main

import (
	"github.com/exelban/logg"
	"log"
)

func main () {
	w := logg.Install(logg.Cyan)
  
	w.Custom("[CUSTOM]", logg.HiBlue, logg.Black, logg.Bold)
	
	log.Print("[CUSTOM] custom text")
}
```

### Level filter usage
```golang
package main

import (
	"github.com/exelban/logg"
	"log"
)

func main () {
	w := logg.Install()
	filter := &logg.LevelFilter{
		Levels: []string{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: "WARN",
	}
	w.SetFilters(filter)
	
	log.Print("[DEBUG] will not be printed")
	log.Print("[INFO] will not be printed")
	log.Print("[WARN] will not printed")
	log.Print("[ERROR] will not printed")
}
```

## Benchmarks

```sh
BenchmarkDiscard-4     	100000000	       12.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoggWrite-4   	 3000000	       448 ns/op	      32 B/op	       2 allocs/op
BenchmarkLogg-4        	 2000000	       825 ns/op	     172 B/op	       3 allocs/op
BenchmarkLog-4         	 3000000	       568 ns/op	      80 B/op	       2 allocs/op
```

`BenchmarkDiscard` - writer to empty buf.  
`BenchmarkLoggWrite` - writer to empty buffer by LoggWriter.  
`BenchmarkLogg` - log using log.Print and installed logg.  
`BenchmarkLog` - log using log.Print (without logg).


# What's new
## 2.0.0
- renamed to logg

## 1.2.0
- added level filter to log
- added benchmarks
- removed blocking goroutine
- moved colors to separate folder
- small fixes

## 1.0.2
- added one more example
- added benchmark if someone want to compare logging to log package
- added one more test

## 1.0.1
- added preinstalled colors for [ERROR], [INFO], [WARN] and [DEBUG]

## 1.0.0
- first release

# Licence
[MIT License](https://github.com/exelban/logg/blob/master/LICENSE)
