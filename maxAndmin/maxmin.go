package main

import (
	"encoding/binary"
	"log"
	"math"
	"os"
	"os/signal"
	"syscall"

	"github.com/nsqio/go-nsq"
)

var flag bool =true
var max float64
var min float64
var sign chan int

type myMessageHandler struct{}

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

func (h *myMessageHandler) processMessage(m []byte) error {
	num := ByteToFloat64(m)

	if flag {
		max = num
		min = num
		flag=false
	}

	if num > max {
		max = num
	}

	if num < min {
		min = num
	}

	sign <- 0

	return nil
}

// HandleMessage implements the Handler interface.
func (h *myMessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	err := h.processMessage(m.Body)
	return err
}

func main() {

	//初始化生产者
	config_p := nsq.NewConfig()
	producer, err := nsq.NewProducer("127.0.0.1:4150", config_p)
	if err != nil {
		log.Fatal(err)
	}

	//处理异步发送的返回状态
	doChan := make(chan *nsq.ProducerTransaction)
	go func() {
		for res := range doChan {
			if res.Error != nil {
				log.Println("send error:", res.Error.Error())
			}
		}
	}()

	sign = make(chan int)

	go func() {
		//异步发送消息 (byte数组)
		for {
			<-sign

			bytes := make([][]byte, 2)

			bytes[0]=Float64ToByte(max)
			bytes[1]=Float64ToByte(min)

			err := producer.MultiPublishAsync("maxmin", bytes, doChan)
			if err != nil {
				log.Fatal("could not pulish:", err)
			}
		}
	}()

	// Instantiate a consumer that will subscribe to the provided channel.
	config_c := nsq.NewConfig()
	consumer, err := nsq.NewConsumer("rand", "maxmin", config_c)
	if err != nil {
		log.Fatal(err)
	}

	consumer.AddHandler(&myMessageHandler{})
	err = consumer.ConnectToNSQLookupd("localhost:4161")
	if err != nil {
		log.Fatal(err)
	}

	// wait for signal to exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Gracefully stop the consumer.
	consumer.Stop()
	producer.Stop()
}
