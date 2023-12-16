# Seqgo

A Seq hook for Logrus

# Sample

```go

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestLog(t *testing.T) {
	log.AddHook(NewSeqHook(func(options *SeqHookOptions) {
		options.fields = map[string]string{
			"System": "Test",
			"Env":    "Dev",
		}
		options.endpoint = "http://localhost:5341"

	}))
	log.Info("hello world1")
	log.Info("hello world2")
	log.Info("hello world3")
    
    // must flush when exit
	Flush()
}

func TestLogWithAdditionProperty(t *testing.T) {
	log.AddHook(NewSeqHook(func(options *SeqHookOptions) {
		options.fields = map[string]string{
			"System": "Test",
			"Env":    "Dev",
		}
		options.endpoint = "http://localhost:5341"

	}))
	log.WithField("NewField", "test").Info("hello world1")

    // must flush when exit
	Flush()
}

```

# Reference Resources

1. [go-queue](https://github.com/yireyun/go-queue)
2. [logruseq](https://github.com/alxyng/logruseq)
3. [seq official document](https://docs.datalust.co/docs/an-overview-of-seq)