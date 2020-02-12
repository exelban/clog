# logg
[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/pkgz/logg)
[![Tests](https://img.shields.io/github/workflow/status/pkgz/logg/Code%20coverage)](https://github.com/pkgz/logg/actions)
[![codecov](https://img.shields.io/codecov/c/gh/pkgz/logg)](https://codecov.io/gh/pkgz/logg)

![](https://serhiy.s3.eu-central-1.amazonaws.com/Github_repo/logg/v3_pretty.png)  
![](https://serhiy.s3.eu-central-1.amazonaws.com/Github_repo/logg/v3_json.png)  
Better log experience in golang.

## Installation
```bash
go get github.com/pkgz/logg
```

## Features
- color logs
- different output formats (pretty/json)
- different levels (debug/info/error/warning/panic)
- can be used with internal log library
- zero allocations

## Usage
There is two way how you can use this logger:

- install globally (this option will set own writer to default log);
- using a local logger instance;

### Global logger
In this mode, the logger will set own writer to log the library. So you can use the default log library for logging. But with better output.

```golang
package main

import (
    "github.com/pkgz/logg"
    "log"
    "os"
)

func main () {
    logg.NewGlobal(os.Stdout)

    log.Print("ERR some message to log")
}
```

### Local logger
In this case, you must define a log instance in each context (function). But in this case, you can define different settings for different scopes.

```golang
package main

import (
    "github.com/pkgz/logg"
    "os"
)

func main () {
    log := logg.New(os.Stdout)

    log.Print("ERR some message to log")
}
```


### Json log
```golang
package main

import (
    "github.com/pkgz/logg"
    "os"
)

func main () {
    log := logg.New(os.Stdout)
    log.SetFormat(logg.Json)

    log.Print("ERR some message to log")
}
```

### Settings
There are a few parameters which you can set:

- flags (define time and caller format. Using the format from internal log library)
- format (output log format. Pretty or Json)
- color (colorize output or not)

#### Levels
- Debug: `DBG | DEBUG | [DBG] | [DEBUG]`
- Info: `INF | INFO | [INF] | [INFO]`
- Error: `ERR | ERROR | [ERR] | [ERROR]`
- Warning: `WRN | WARN | [WRN] | [WARN]`
- Panic: `PNC | PANIC | [PNC] | [PANIC]`

#### API
| Function | Default | Description |
| --- | --- | --- |
| `SetWriter(io.Writer) ` | ioutil.Discard | Set writer. |
| `SetFormat(logg.format) ` | Pretty | Set output format. Can be pretty or json. |
| `SetFlags(int) ` | int | Set time and caller flags. |
| `MinLevel(level) ` | Debug | Minimum level for logs. Logs lower this level will be not writed. |
| `ToggleColor(bool) ` | true | Enable or disable output colorizing. |
| `DebugMode() ` | | Will enable a debug mode. Debug mode will add milliseconds to timestamp and log caller. |

## Benchmarks

```sh
BenchmarkLog_Print/long_message-8         	 1799361	       703 ns/op	     592 B/op	       2 allocs/op
BenchmarkLog_Print/short_message-8        	 2444176	       491 ns/op	      80 B/op	       2 allocs/op
BenchmarkLog_Print/short_level-8          	 2451745	       494 ns/op	      80 B/op	       2 allocs/op
BenchmarkLog_Print/long_level-8           	 2364008	       493 ns/op	      80 B/op	       2 allocs/op
BenchmarkLogg_Log_Print/long_message-8    	 1443676	       825 ns/op	     593 B/op	       2 allocs/op
BenchmarkLogg_Log_Print/short_message-8   	 1701645	       713 ns/op	      80 B/op	       2 allocs/op
BenchmarkLogg_Log_Print/short_level-8     	 1882834	       636 ns/op	      80 B/op	       2 allocs/op
BenchmarkLogg_Log_Print/long_level-8      	 1862190	       676 ns/op	      80 B/op	       2 allocs/op
BenchmarkLogg_Write_Pretty/long_message-8 	13835029	        90.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogg_Write_Pretty/short_message-8         	13607654	        86.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogg_Write_Pretty/short_level-8           	16418979	        73.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogg_Write_Pretty/long_level-8            	15362362	        76.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogg_Write_Json/long_message-8            	11209905	       107 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogg_Write_Json/short_message-8           	11571448	       105 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogg_Write_Json/short_level-8             	13749747	        87.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogg_Write_Json/long_level-8              	13415761	        89.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogg_Print/long_message-8                 	 4293176	       311 ns/op	    1170 B/op	       3 allocs/op
BenchmarkLogg_Print/short_message-8                	 7748691	       155 ns/op	     144 B/op	       3 allocs/op
BenchmarkLogg_Print/short_level-8                  	 8459841	       141 ns/op	     144 B/op	       3 allocs/op
BenchmarkLogg_Print/long_level-8                   	 8349727	       143 ns/op	     144 B/op	       3 allocs/op
BenchmarkLogg_Printf/long_message-8                	 3314443	       358 ns/op	    1154 B/op	       2 allocs/op
BenchmarkLogg_Printf/short_message-8               	 7906260	       153 ns/op	     128 B/op	       2 allocs/op
BenchmarkLogg_Printf/short_level-8                 	 8624889	       144 ns/op	     128 B/op	       2 allocs/op
BenchmarkLogg_Printf/long_level-8                  	 8253075	       145 ns/op	     128 B/op	       2 allocs/op
```


## Licence
[MIT License](https://github.com/pkgz/logg/blob/master/LICENSE)
