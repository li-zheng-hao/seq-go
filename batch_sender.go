package seqgo

import (
	"fmt"
	"github/li-zheng-hao/seqgo/queue"
	"io"
	"net/http"
	"strings"
	"time"
)

type MessageBatchSender struct {
	messages *queue.EsQueue
}

var messageBatchSender = &MessageBatchSender{
	messages: queue.NewQueue(10000),
}

func Flush() {
	send()
}

func send() {

	messages := make([]interface{}, SeqHookOption.batchSize)

	messageBatchSender.messages.Gets(messages)
	if messages[0] == nil {
		return
	}
	combinedMessage := strings.Builder{}

	for i := range messages {
		if messages[i] == nil {
			continue
		}
		combinedMessage.Write(messages[i].([]byte))
	}

	str := combinedMessage.String()
	req, err := http.NewRequest("POST", SeqHookOption.endpoint, strings.NewReader(str))
	if err != nil {
		fmt.Errorf("error seq post %v", err)
	}

	req.Header.Add("Content-Type", "application/vnd.serilog.clef")

	if SeqHookOption.apiKey != "" {
		req.Header.Add("X-Seq-ApiKey", SeqHookOption.apiKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Errorf("error creating seq event: %v", err)

		}

		fmt.Errorf("error creating seq event: %v", string(data))
	}
}

// Push message
func Push(message []byte) {
	messageBatchSender.messages.Put(message)
}

func ScheduleSend() {
	for {
		if messageBatchSender.messages.Quantity() == 0 {
			timeSec := 1
			if SeqHookOption == nil {
				time.Sleep(time.Duration(timeSec) * time.Second)
			} else {
				timeSec = SeqHookOption.period
			}

			time.Sleep(time.Duration(timeSec) * time.Second)
		} else {
			send()
		}
	}
}
