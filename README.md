# Seqgo

A Seq hook for Logrus

# Usage

```go
import (
	log "github.com/sirupsen/logrus"
	. "github/li-zheng-hao/seqgo"
	"time"
)

func main() {
	hook := NewSeqHook(func(options *SeqHookOptions) {
		options.BatchSize = 10
		options.Fields = map[string]string{
			"System": "Test",
			"Env":    "Dev",
		}
		options.Endpoint = "http://localhost:5341"

	})
	log.AddHook(hook)

	for i := 0; i < 10; i++ {
		log.Info(time.Now().String())

	}
	hook.Flush()

}
```

# Reference Resources

1. [go-queue](https://github.com/yireyun/go-queue)
2. [logruseq](https://github.com/alxyng/logruseq)
3. [seq official document](https://docs.datalust.co/docs/an-overview-of-seq)