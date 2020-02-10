package logg

import (
	"github.com/rs/zerolog"
	"io/ioutil"
	"log"
	"testing"
	"time"
)

var testMessage = []byte("[INFO] test logging, but use a somewhat realistic message length.")

/******************************************************************************
*                                 Logg.Write                                  *
******************************************************************************/

func BenchmarkLogg_Write(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	logger := New(ioutil.Discard)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = logger.Write(testMessage)
		}
	})
}

func BenchmarkZerolog_Write(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	logger := zerolog.New(ioutil.Discard)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = logger.Write(testMessage)
		}
	})
}

/******************************************************************************
*                           Internal log.Print                                *
******************************************************************************/

func Benchmark_Log(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	log.SetOutput(ioutil.Discard)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Print(testMessage)
		}
	})
}

/******************************************************************************
*                          Logg vs Zerolog (log)                              *
******************************************************************************/

func BenchmarkLogg_Log(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	logger := New(ioutil.Discard)
	log.SetOutput(logger)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Print(testMessage)
		}
	})
}

func BenchmarkZerolog_Log(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	output := zerolog.ConsoleWriter{Out: ioutil.Discard, TimeFormat: time.RFC3339}
	logger := zerolog.New(output)
	log.SetOutput(logger)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Print(testMessage)
		}
	})
}

func BenchmarkLoggLog_Json(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	logger := New(ioutil.Discard)
	logger.format = Json

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Print(testMessage)
		}
	})
}

func BenchmarkZerolog_Json(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	logger := zerolog.New(ioutil.Discard)
	log.SetOutput(logger)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Print(testMessage)
		}
	})
}

/******************************************************************************
*                           Internal functions                                *
******************************************************************************/

func BenchmarkMessage_defineTime(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	m := &message{
		t: time.Now(),
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.defineTime()
		}
	})
}
