package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/m1m0ry/mom/mq"
	"github.com/nsqio/go-nsq"
)

// 产生期望a，标准差b，服从正态分布的随机数
func generateNumber(a float64, b float64) float64 {
	return rand.NormFloat64()*b + a
}

func main() {

	//初始化生产者
	producer, err := mq.NewMessageQueue(mq.MessageQueueConfig{})
	if err != nil {
		log.Fatal(err)
	}
	producer.Run()

	//主题 rand
	topicName := "rand"

	//处理异步发送的返回状态
	doChan := make(chan *nsq.ProducerTransaction)
	go func() {
		for res := range doChan {
			if res.Error != nil {
				log.Println("send error:", res.Error.Error())
			}
		}
	}()

	//异步发送消息
	for {
		num := generateNumber(20, 10)
		err := producer.PubAsync(topicName, num, doChan)
		if err != nil {
			log.Fatal("could not pulish:", err)
		}

		// 休眠100毫秒
		time.Sleep(time.Millisecond * 100)
	}
}
