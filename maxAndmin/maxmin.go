package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/m1m0ry/mom/mq"
	"github.com/nsqio/go-nsq"
)

var max float64
var min float64
var maxmin chan int

func main() {

	mq, err := mq.NewMessageQueue(mq.MessageQueueConfig{
		SupportedTopics: []string{"rand"},
		Channel:         "maxmin",
	})
	if err != nil {
		log.Fatal(err)
	}

	mq.Sub("rand", func(resp []byte) {
		var num float64
		err := json.Unmarshal(resp, &num)
		if err != nil {
			log.Fatal("unmarshal err: ", err)
		}

		if min == 0 && max == 0 {
			max = num
			min = num
		}

		if num > max {
			max = num
		}

		if num < min {
			min = num
		}

		maxmin <- 0
	})

	mq.Run()

	maxmin = make(chan int)

	//处理异步发送的返回状态
	doChan := make(chan *nsq.ProducerTransaction)
	go func() {
		for res := range doChan {
			if res.Error != nil {
				log.Println("send error:", res.Error.Error())
			}
		}
	}()

	go func() {
		//异步发送消息 (byte数组)
		for {
			<-maxmin
			err := mq.PubAsync("maxmin", map[string]float64{"max": max, "min": min}, doChan)
			if err != nil {
				log.Fatal("could not pulish:", err)
			}
		}
	}()

	// wait for signal to exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	mq.StopProducer()
	mq.StopConsumer("rand")
}
