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

const (
	N = 300
)

var globalNums []float64
var variance float64
var sum float64
var sum2 float64
var mean float64
var variance_chan chan float64

func main() {

	mq, err := mq.NewMessageQueue(mq.MessageQueueConfig{
		SupportedTopics: []string{"rand"},
		Channel:         "variance",
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

		if len(globalNums) > N {
			pre := globalNums[0]
			globalNums = globalNums[1:]
			old_mean := mean
			mean = mean + (num-pre)/N
			variance = variance + (mean-old_mean)/N*((N-1)*(pre+num)-2*N*old_mean+2*pre)
		} else {
			sum += num
			mean = sum / (float64(len(globalNums)) + 1)
			sum2 += (num - mean) * (num - mean)
			variance = sum2 / (float64(len(globalNums)) + 1)
		}

		variance_chan <- variance

		globalNums = append(globalNums, num)
	})

	mq.Run()

	//处理异步发送的返回状态
	doChan := make(chan *nsq.ProducerTransaction)
	go func() {
		for res := range doChan {
			if res.Error != nil {
				log.Println("send error:", res.Error.Error())
			}
		}
	}()

	variance_chan = make(chan float64)

	go func() {
		//异步发送消息 (byte数组)
		for {
			err := mq.PubAsync("variance", <-variance_chan, doChan)
			if err != nil {
				log.Fatal("could not pulish:", err)
			}
		}
	}()

	// wait for signal to exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Gracefully stop the consumer.
	mq.StopProducer()
	mq.StopConsumer("rand")
}
