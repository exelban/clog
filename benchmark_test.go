package logg

import (
	"io/ioutil"
	"log"
	"testing"
)

var testMessage = []byte("[INFO] test logging, but use a somewhat realistic message length.")

func BenchmarkLog(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	log.SetOutput(ioutil.Discard)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Print(testMessage)
		}
	})
}

func BenchmarkLoggWrite(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	Logger.out = ioutil.Discard
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = Logger.Write(testMessage)
		}
	})
}

func BenchmarkLoggLogPretty(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	Logger.out = ioutil.Discard
	Logger.format = Pretty
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Print(testMessage)
		}
	})
}

func BenchmarkLoggLogJson(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	Logger.out = ioutil.Discard
	Logger.format = Json
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Print(testMessage)
		}
	})
}

//func BenchmarkZerologWrite(b *testing.B) {
//	b.ReportAllocs()
//	b.ResetTimer()
//
//	logger := zerolog.New(ioutil.Discard)
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			_, _ = logger.Write(testMessage)
//		}
//	})
//}
//
//func BenchmarkZerologLog(b *testing.B) {
//	b.ReportAllocs()
//	b.ResetTimer()
//
//	logger := zerolog.New(ioutil.Discard)
//	log.SetOutput(logger)
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			log.Print(testMessage)
//		}
//	})
//}
