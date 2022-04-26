package analysis

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/m1m0ry/mom/mq"
)

const (
	N = 300
)

type RowData struct {
	Datas    []float64 `json:"data"`
	Min      float64   `json:"min"`
	Max      float64   `json:"max"`
	Variance float64   `json:"variance"`
	Mean     float64   `json:"mean"`
}

var Data RowData

func rand(resp []byte) {
	var num float64
	err := json.Unmarshal(resp, &num)
	if err != nil {
		if err != nil {
			log.Fatal("unmarshal err: ", err)
		}
	}
	if len(Data.Datas) > N {
		Data.Datas = Data.Datas[1:]
	}
	Data.Datas = append(Data.Datas, num)
}

func variance(resp []byte) {
	var num float64
	err := json.Unmarshal(resp, &num)
	if err != nil {
		if err != nil {
			log.Fatal("unmarshal err: ", err)
		}
	}
	Data.Variance = num
}

func mean(resp []byte) {
	var num float64
	err := json.Unmarshal(resp, &num)
	if err != nil {
		if err != nil {
			log.Fatal("unmarshal err: ", err)
		}
	}
	Data.Mean = num
}

func maxmin(resp []byte) {
	var maxmin map[string]float64
	err := json.Unmarshal(resp, &maxmin)
	if err != nil {
		if err != nil {
			log.Fatal("unmarshal err: ", err)
		}
	}
	Data.Max = maxmin["max"]
	Data.Min = maxmin["min"]
}

func Init() {
	mq, err := mq.NewMessageQueue(mq.MessageQueueConfig{
		SupportedTopics: []string{"rand", "mean", "variance", "maxmin"},
		Channel:         "chart",
	})
	if err != nil {
		log.Fatal(err)
	}

	mq.Sub("rand", func(resp []byte) { rand(resp) })
	mq.Sub("mean", func(resp []byte) { mean(resp) })
	mq.Sub("variance", func(resp []byte) { variance(resp) })
	mq.Sub("maxmin", func(resp []byte) { maxmin(resp) })

	mq.Run()

	// wait for signal to exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Gracefully stop the consumer.
	mq.StopConsumer("rand")
	mq.StopConsumer("mean")
	mq.StopConsumer("variance")
	mq.StopConsumer("maxmin")
}
