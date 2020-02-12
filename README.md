# logg
[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/pkgz/logg)
[![Tests](https://img.shields.io/github/workflow/status/pkgz/logg/Code)](https://github.com/pkgz/logg/actions)
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
| `SetFlags(int) ` | log.Ltime | Set time and caller flags. |
| `MinLevel(level) ` | Debug | Minimum level for logs. Logs lower this level will be not writed. |
| `ToggleColor(bool) ` | true | Enable or disable output colorizing. |
| `DebugMode() ` | | Will enable a debug mode. Debug mode will add milliseconds to timestamp and log caller. |

## Benchmarks

```sh
BenchmarkLog_Print/long_message-8         	        1676541	            908 ns/op	     592 B/op	       2 allocs/op
BenchmarkLog_Print/short_message-8        	        1700326	            690 ns/op	      80 B/op	       2 allocs/op
BenchmarkLog_Print/short_level-8          	        1806188	            710 ns/op	      80 B/op	       2 allocs/op
BenchmarkLog_Print/long_level-8           	        1780255	            781 ns/op	      80 B/op	       2 allocs/op
BenchmarkLogg_Log_Print/long_message-8    	        883708	            1250 ns/op	     593 B/op	       2 allocs/op
BenchmarkLogg_Log_Print/short_message-8   	        1000000	            1034 ns/op	      80 B/op	       2 allocs/op
BenchmarkLogg_Log_Print/short_level-8     	        1419847	            926 ns/op	      80 B/op	       2 allocs/op
BenchmarkLogg_Log_Print/long_level-8      	        1576074	            701 ns/op	      80 B/op	       2 allocs/op
BenchmarkLogg_Write_Pretty/long_message-8 	        12929612	        102 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogg_Write_Pretty/short_message-8         	12415136	        93.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogg_Write_Pretty/short_level-8           	14326122	        79.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogg_Write_Pretty/long_level-8            	14781205	        82.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogg_Write_Json/long_message-8            	 9923089	        112 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogg_Write_Json/short_message-8           	11412260	        106 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogg_Write_Json/short_level-8             	13234461	        94.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogg_Write_Json/long_level-8              	13051018	        94.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogg_Print/long_level-8                   	 8072880	        166 ns/op	     144 B/op	       3 allocs/op
BenchmarkLogg_Print/long_message-8                 	 3646544	        329 ns/op	    1170 B/op	       3 allocs/op
BenchmarkLogg_Print/short_message-8                	 7247292	        174 ns/op	     144 B/op	       3 allocs/op
BenchmarkLogg_Print/short_level-8                  	 6156025	        179 ns/op	     144 B/op	       3 allocs/op
BenchmarkLogg_Printf/short_level-8                 	 7884908	        153 ns/op	     128 B/op	       2 allocs/op
BenchmarkLogg_Printf/long_level-8                  	 8083395	        151 ns/op	     128 B/op	       2 allocs/op
BenchmarkLogg_Printf/long_message-8                	 3014763	        409 ns/op	    1154 B/op	       2 allocs/op
BenchmarkLogg_Printf/short_message-8               	 7392619	        167 ns/op	     128 B/op	       2 allocs/op
```


## Licence
[MIT License](https://github.com/pkgz/logg/blob/master/LICENSE)
