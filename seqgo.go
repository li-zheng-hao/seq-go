package seqgo

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

// SeqHook sends logs to Seq via HTTP.
type SeqHook struct {
	endpoint      string
	apiKey        string
	levels        []logrus.Level
	messageSender *MessageBatchSender
}

func (hook *SeqHook) Levels() []logrus.Level {
	return hook.levels
}

var SeqHookOption *SeqHookOptions = &SeqHookOptions{
	Levels: []logrus.Level{
		logrus.TraceLevel,
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	},
	Period:       2,
	BatchSize:    10,
	MaxLimitSize: 10000,
}

func NewSeqHook(configure func(*SeqHookOptions)) *SeqHook {

	configure(SeqHookOption)

	endpoint := fmt.Sprintf("%v/api/events/raw", SeqHookOption.Endpoint)

	SeqHookOption.Endpoint = endpoint

	seqHook := &SeqHook{
		endpoint:      endpoint,
		apiKey:        SeqHookOption.ApiKey,
		levels:        SeqHookOption.Levels,
		messageSender: NewMessageBatchSender(SeqHookOption.MaxLimitSize),
	}

	go ScheduleSend(seqHook.messageSender)

	return seqHook
}

// Fire sends a log entry to Seq.
func (hook *SeqHook) Fire(entry *logrus.Entry) error {
	formatter := logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg:   "@mt",
			logrus.FieldKeyLevel: "@l",
			logrus.FieldKeyTime:  "@t",
		},
	}
	for k, v := range SeqHookOption.Fields {
		entry.Data[k] = v
	}
	data, err := formatter.Format(entry)

	if err != nil {
		return err
	}

	hook.messageSender.Push(data)

	return nil
}

// SeqHookOptions collects non-default Seq hook options.
type SeqHookOptions struct {
	ApiKey       string
	Levels       []logrus.Level
	Period       int
	Fields       map[string]string
	BatchSize    int
	Endpoint     string
	MaxLimitSize uint32
}
