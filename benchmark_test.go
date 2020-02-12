package logg

import (
	"io/ioutil"
	"log"
	"testing"
)

var (
	testMessage           = []byte("[INFO] test logging, but use a somewhat realistic message length.")
	shortLevelTestMessage = []byte("ERR test logging, but use a somewhat realistic message length.")
	longLevelTestMessage  = []byte("PANIC test logging, but use a somewhat realistic message length.")
)

func BenchmarkLogg_New(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger := New(ioutil.Discard)
			logger.minLevel = Info
		}
	})
}

func BenchmarkLogg_Write(b *testing.B) {
	logger := New(ioutil.Discard)

	b.Run("test", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for n := 0; n < b.N; n++ {
			_, _ = logger.Write(testMessage)
		}
	})

	b.Run("short", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for n := 0; n < b.N; n++ {
			_, _ = logger.Write(shortLevelTestMessage)
		}
	})

	b.Run("long", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for n := 0; n < b.N; n++ {
			_, _ = logger.Write(longLevelTestMessage)
		}
	})
}

func BenchmarkLogg_Write_JSON(b *testing.B) {
	logger := New(ioutil.Discard)
	logger.format = Json

	b.Run("test", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for n := 0; n < b.N; n++ {
			_, _ = logger.Write(testMessage)
		}
	})

	b.Run("short", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for n := 0; n < b.N; n++ {
			_, _ = logger.Write(shortLevelTestMessage)
		}
	})

	b.Run("long", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for n := 0; n < b.N; n++ {
			_, _ = logger.Write(longLevelTestMessage)
		}
	})
}

func BenchmarkLogg_Print(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	logger := New(ioutil.Discard)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Print(testMessage)
		}
	})
}

func BenchmarkLogg_Printf(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	logger := New(ioutil.Discard)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Printf("%s", testMessage)
		}
	})
}

func BenchmarkLogg_Debug(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	logger := New(ioutil.Discard)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Debug(testMessage)
		}
	})
}

func BenchmarkLogg_Debugf(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	logger := New(ioutil.Discard)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Debugf("%s", testMessage)
		}
	})
}

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

func BenchmarkLoggLog_Json(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	logger := New(ioutil.Discard)
	logger.format = Json
	log.SetOutput(logger)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Print(testMessage)
		}
	})
}
