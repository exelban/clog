# clog
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/exelban/clog)
[![codecov](https://codecov.io/gh/exelban/clog/branch/master/graph/badge.svg?token=A8eLVAj9cU)](https://codecov.io/gh/exelban/clog)

![](https://s3.eu-central-1.amazonaws.com/serhiy/Github_repo/clog/Zrzut+ekranu+2018-10-16+o+18.52.26.png)  
Color logs for your go application.

# Installation
```bash
go get github.com/exelban/clog
```

# Usage
## Example
```golang
package main

import (
	"github.com/exelban/clog"
	"log"
)

func main () {
	w := clog.Install(clog.Cyan)
  
	w.Custom("[ERROR]", clog.Red)

	log.Print("[ERROR] error text")
}
```

# Licence
[MIT License](https://github.com/exelban/clog/blob/master/LICENSE)
