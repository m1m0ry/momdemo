package main

import (
	"encoding/binary"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/nsqio/go-nsq"
)

// 产生期望a，标准差b，服从正态分布的随机数
func generateNumber(a float64, b float64) float64 {
	return rand.NormFloat64()*b + a
}

//Float64ToByte Float64转byte
func Float64ToByte(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

//ByteToFloat64 byte转Float64
func ByteToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	return math.Float64frombits(bits)
}

func main() {

	//初始化生产者
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer("127.0.0.1:4150", config)
	if err != nil {
		log.Fatal(err)
	}

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

	//异步发送消息 (byte数组)
	for {
		num := generateNumber(20, 10)
		err := producer.PublishAsync(topicName, Float64ToByte(num), doChan)
		if err != nil {
			log.Fatal("could not pulish:", err)
		}

		// 休眠100毫秒
		time.Sleep(time.Millisecond * 100)
	}
}
