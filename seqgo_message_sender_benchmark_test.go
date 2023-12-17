package seqgo

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
	"testing"
	"time"
)

func BenchmarkLogByte(b *testing.B) {

	for i := 0; i < b.N; i++ {
		combinedMessage := &bytes.Buffer{}
		combinedMessage.Write([]byte("hello world"))
	}

}
func BenchmarkLogString(b *testing.B) {

	for i := 0; i < b.N; i++ {
		combinedMessage := &strings.Builder{}
		combinedMessage.Write([]byte("hello world"))
	}

}

func BenchmarkLogBytePool(b *testing.B) {

	pool := &sync.Pool{New: func() interface{} { return new(bytes.Buffer) }}
	for i := 0; i < b.N; i++ {
		combinedMessage := pool.Get().(*bytes.Buffer)
		combinedMessage.Write([]byte("hello world"))
		pool.Put(combinedMessage)

	}

}
func BenchmarkLogStringPool(b *testing.B) {
	pool := &sync.Pool{New: func() interface{} { return new(strings.Builder) }}

	for i := 0; i < b.N; i++ {
		combinedMessage := pool.Get().(*strings.Builder)
		combinedMessage.Write([]byte("hello world"))
		pool.Put(combinedMessage)
	}

}

func BenchmarkSendMessage(b *testing.B) {
	hook := NewSeqHook(func(options *SeqHookOptions) {
		options.batchSize = 10
		options.fields = map[string]string{
			"System": "Test",
			"Env":    "Dev",
		}
		options.endpoint = "http://localhost:5341"

	})
	log.AddHook(hook)

	for i := 0; i < b.N; i++ {
		log.Info(time.Now().String())
		hook.Flush()
	}

}
