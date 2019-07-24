# clog
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/exelban/clog)
[![codecov](https://codecov.io/gh/exelban/clog/branch/master/graph/badge.svg)](https://codecov.io/gh/exelban/clog)

![](https://s3.eu-central-1.amazonaws.com/serhiy/Github_repo/clog/Zrzut+ekranu+2018-10-16+o+18.52.26.png)  
Color logs for your go application.

# Installation
```bash
go get github.com/exelban/clog
```

# Usage
## Example

### Simple usage
```golang
package main

import (
	"github.com/exelban/clog"
	"log"
)

func main () {
	clog.Install()
	
	log.Print("[ERROR] error text")
}
```

### Custom level
```golang
package main

import (
	"github.com/exelban/clog"
	"log"
)

func main () {
	w := clog.Install(clog.Cyan)
  
	w.Custom("[CUSTOM]", clog.HiBlue, clog.Black, clog.Bold)
	
	log.Print("[CUSTOM] custom text")
}
```

### Level filter usage
```golang
package main

import (
	"github.com/exelban/clog"
	"log"
)

func main () {
	w := clog.Install()
	filter := &clog.LevelFilter{
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

# What's new
## 1.2.0
- added level filter to log
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
[MIT License](https://github.com/exelban/clog/blob/master/LICENSE)
