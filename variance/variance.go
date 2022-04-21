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

const (
	N = 300
)

var globalNums []float64
var variance float64
var sum float64
var sum2 float64
var mean float64
var variance_chan chan float64

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

	if len(globalNums) > N {
		pre := globalNums[0]
		globalNums = globalNums[1:]
		old_mean:=mean
		mean = mean + (num-pre)/N
		variance=variance+(mean-old_mean)/N*((N-1)*(pre+num)-2*N*old_mean+2*pre)
	} else {
		sum += num
		mean = sum / (float64(len(globalNums)) + 1)
		sum2+=(num-mean)*(num-mean)
		variance=sum2/ (float64(len(globalNums)) + 1)
	}

	variance_chan <- variance

	globalNums = append(globalNums, num)

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

	variance_chan = make(chan float64)

	go func() {
		//异步发送消息 (byte数组)
		for {
			err := producer.PublishAsync("variance", Float64ToByte(<-variance_chan), doChan)
			if err != nil {
				log.Fatal("could not pulish:", err)
			}
		}
	}()

	// Instantiate a consumer that will subscribe to the provided channel.
	config_c := nsq.NewConfig()
	consumer, err := nsq.NewConsumer("rand", "variance", config_c)
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