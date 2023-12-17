package seqgo

import (
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	hook := NewSeqHook(func(options *SeqHookOptions) {
		options.batchSize = 10
		options.fields = map[string]string{
			"System": "Test",
			"Env":    "Dev",
		}
		options.endpoint = "http://localhost:5341"

	})
	log.AddHook(hook)

	for i := 0; i < 10; i++ {
		log.Info(time.Now().String())
		hook.Flush()

	}

}

func TestLogWithAdditionProperty(t *testing.T) {
	hook := NewSeqHook(func(options *SeqHookOptions) {
		options.batchSize = 10
		options.fields = map[string]string{
			"System": "Test",
			"Env":    "Dev",
		}
		options.endpoint = "http://localhost:5341"

	})
	log.AddHook(hook)
	log.WithField("NewField", "test").Info("hello world1")
	hook.Flush()
}
