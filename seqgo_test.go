package seqgo

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestLog(t *testing.T) {
	log.AddHook(NewSeqHook(func(options *SeqHookOptions) {
		options.batchSize = 10
		options.fields = map[string]string{
			"System": "Test",
			"Env":    "Dev",
		}
		options.endpoint = "http://localhost:5341"

	}))
	log.Info("hello world1")
	log.Info("hello world2")
	log.Info("hello world3")
	Flush()
}

func TestLogWithAdditionProperty(t *testing.T) {
	log.AddHook(NewSeqHook(func(options *SeqHookOptions) {
		options.batchSize = 10
		options.fields = map[string]string{
			"System": "Test",
			"Env":    "Dev",
		}
		options.endpoint = "http://localhost:5341"

	}))
	log.WithField("NewField", "test").Info("hello world1")
	Flush()
}
