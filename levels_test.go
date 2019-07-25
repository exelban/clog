package logg

import (
	"github.com/rs/zerolog"
	"io/ioutil"
	"log"
	"testing"
)

var messages = [][]byte{
	[]byte("[TRACE] foo"),
	[]byte("[DEBUG] foo"),
	[]byte("[INFO] foo"),
	[]byte("[WARN] foo"),
	[]byte("[ERROR] foo"),
}

func BenchmarkDiscard(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = ioutil.Discard.Write(messages[i%len(messages)])
	}
}

func BenchmarkLog(b *testing.B) {
	b.ReportAllocs()

	log.SetOutput(ioutil.Discard)
	for i := 0; i < b.N; i++ {
		log.Print(messages[i%len(messages)])
	}
}

func BenchmarkLoggWrite(b *testing.B) {
	b.ReportAllocs()

	Logger.out = ioutil.Discard

	for i := 0; i < b.N; i++ {
		_, _ = Logger.Write(messages[i%len(messages)])
	}
}

func BenchmarkLoggLog(b *testing.B) {
	b.ReportAllocs()

	Logger.out = ioutil.Discard
	for i := 0; i < b.N; i++ {
		log.Print(messages[i%len(messages)])
	}
}

func BenchmarkZerologWrite(b *testing.B) {
	b.ReportAllocs()

	logger := zerolog.New(ioutil.Discard)
	for i := 0; i < b.N; i++ {
		_, _ = logger.Write(messages[i%len(messages)])
	}
}

func BenchmarkZerologLog(b *testing.B) {
	b.ReportAllocs()

	logger := zerolog.New(ioutil.Discard)
	log.SetOutput(logger)
	for i := 0; i < b.N; i++ {
		log.Print(messages[i%len(messages)])
	}
}
