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
		log.Printf("rand error decoding sakura response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
		log.Printf("sakura response: %q", resp)
	} else {
		if len(Data.Datas) > N {
			Data.Datas = Data.Datas[1:]
		}
		Data.Datas = append(Data.Datas, num)
	}
}

func variance(resp []byte) {
	var num float64
	err := json.Unmarshal(resp, &num)
	if err != nil {
		log.Printf("variance error decoding sakura response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
		log.Printf("sakura response: %q", resp)
	} else {
		Data.Variance = num
	}
}

func mean(resp []byte) {
	var num float64
	err := json.Unmarshal(resp, &num)
	if err != nil {
		log.Printf("mean error decoding sakura response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
		log.Printf("sakura response: %q", resp)
	} else {
		Data.Mean = num
	}
}

func maxmin(resp []byte) {
	var maxmin map[string]float64
	err := json.Unmarshal(resp, &maxmin)
	if err != nil {
		log.Printf("maxmin error decoding sakura response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
		log.Printf("sakura response: %q", resp)
	} else {
		Data.Max = maxmin["max"]
		Data.Min = maxmin["min"]
	}
}

func Init() {

	file, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("cant open log: ", err)
	}
	defer file.Close()

	log.SetOutput(file)
	log.SetPrefix("chart init(): ")

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
