package seqgo

import (
	"fmt"
	"github.com/li-zheng-hao/seqgo/queue"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type MessageBatchSender struct {
	messageQueue *queue.EsQueue
	messagePool  *sync.Pool
}

func (hook *SeqHook) Flush() {
	hook.messageSender.Send()
}

func NewMessageBatchSender(poolSize uint32) *MessageBatchSender {
	return &MessageBatchSender{
		messageQueue: queue.NewQueue(poolSize),
		messagePool: &sync.Pool{
			New: func() interface{} {
				return &strings.Builder{}
			},
		},
	}
}

func (ms *MessageBatchSender) Send() {

	combinedMessage := ms.messagePool.Get().(*strings.Builder)

	defer ms.messagePool.Put(combinedMessage)

	for i := 0; i < SeqHookOption.BatchSize; i++ {
		msg, ok, _ := ms.messageQueue.Get()
		if !ok {
			break
		}
		combinedMessage.Write(msg.([]byte))
	}

	if combinedMessage.Len() == 0 {
		return
	}

	str := combinedMessage.String()

	req, err := http.NewRequest("POST", SeqHookOption.Endpoint, strings.NewReader(str))
	if err != nil {
		fmt.Errorf("error seq post %v", err)
	}

	req.Header.Add("Content-Type", "application/vnd.serilog.clef")

	if SeqHookOption.ApiKey != "" {
		req.Header.Add("X-Seq-ApiKey", SeqHookOption.ApiKey)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Errorf("error creating http request: %v", err)

	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Errorf("error creating seq event and read response body error : %v", err)

		}

		fmt.Errorf("error creating seq event: %v", string(data))
	}
}

// Push message
func (ms *MessageBatchSender) Push(message []byte) {
	ms.messageQueue.Put(message)
}

func ScheduleSend(sender *MessageBatchSender) {
	for {
		if sender.messageQueue.Quantity() == 0 {
			timeSec := 1
			if SeqHookOption == nil {
				time.Sleep(time.Duration(timeSec) * time.Second)
			} else {
				timeSec = SeqHookOption.Period
			}

			time.Sleep(time.Duration(timeSec) * time.Second)
		} else {
			sender.Send()
		}
	}
}
