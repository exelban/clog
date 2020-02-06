# logg
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/pkgz/logg)
[![codecov](https://codecov.io/gh/pkgz/logg/branch/master/graph/badge.svg)](https://codecov.io/gh/pkgz/logg)

![](https://serhiy.s3.eu-central-1.amazonaws.com/Github_repo/logg/v2.0.0-1.png)  
![](https://serhiy.s3.eu-central-1.amazonaws.com/Github_repo/logg/v2.0.0-2.png)  
Better log experience in golang.

## Installation
```bash
go get github.com/pkgz/logg
```

## Usage

### Example
#### Simple usage
```golang
package main

import (
	_ "github.com/pkgz/logg"
	"log"
)

func main () {
	log.Print("[ERROR] error text")
}
```

#### Json logs
```golang
package main

import (
	"github.com/pkgz/logg"
	"log"
)

func main () {
	logg.SetFormat(logg.Json)
	
	log.Print("message")
}
```

#### Level filter usage
```golang
package main

import (
	"github.com/pkgz/logg"
	"log"
)

func main () {
	logg.SetMinLevel("INFO")
	
	log.Print("[DEBUG] will not be printed")
	log.Print("[INFO] will be printed")
	log.Print("[WARN] will be printed")
	log.Print("[ERROR] will be printed")
}
```

### Configuration

| Function | Default | Description |
| --- | --- | --- |
`SetOutput(io.Writer) ` | os.Stderr | Sets the output destination for the standard logger. |
`SetFormat(logg.format) ` | Pretty | Sets the output format (`Pretty` or `Json`) for the logger. |
| `SetFlags(int) ` | log.Ltime | Sets the output flags for the logger. Accept the dafault log flags. |
| `SetDebug() ` | false | Sets the output flags prepared to debug for the logger. |
| `SetLevel([]string) ` | `DEBUG, INFO, WARN, ERROR` | Sets the levels of logs. |
| `SetMinLevel(string) ` | `INFO` | Set the minimum levels of logs. |
| `CustomColor(string, int) ` | | Allow to set custom colors for prefix |

## Benchmarks

```sh
BenchmarkLog-8             	 3000000	       506 ns/op	     272 B/op	       2 allocs/op
BenchmarkLoggWrite-8       	 2000000	       914 ns/op	     160 B/op	       4 allocs/op
BenchmarkLoggLogPretty-8   	 3000000	       541 ns/op	     272 B/op	       2 allocs/op
BenchmarkLoggLogJson-8     	 2000000	       576 ns/op	     272 B/op	       2 allocs/op
```

`BenchmarkLog` - log.Print without installed Logg.  
`BenchmarkLoggWrite` - writer to empty buffer by LoggWriter.  
`BenchmarkLoggLogPretty` - log.Print in pretty format.  
`BenchmarkLoggLogJson` - log.Print in JSON format.


## Licence
[MIT License](https://github.com/pkgz/logg/blob/master/LICENSE)
