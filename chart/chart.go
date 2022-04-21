package main

import (
	"encoding/binary"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/nsqio/go-nsq"
)

const (
	N = 300
)

var globalNums []float64

type myMessageHandler struct{}

//ByteToFloat64 byte转Float64
func ByteToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	return math.Float64frombits(bits)
}

func (h *myMessageHandler) processMessage(m []byte) error {
	num := ByteToFloat64(m)
	// fmt.Println(num)
	if len(globalNums) > N {
		globalNums = globalNums[1:]
	}
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

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//webSocket请求ping 返回pong
func websocketapi(c *gin.Context) {
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	for {
		//写入ws数据
		err = ws.WriteJSON(globalNums)
		if err != nil {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func main() {
	// Instantiate a consumer that will subscribe to the provided channel.

	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer("rand", "chart", config)
	if err != nil {
		log.Fatal(err)
	}

	consumer.AddHandler(&myMessageHandler{})
	err = consumer.ConnectToNSQLookupd("localhost:4161")
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.LoadHTMLGlob("index.html")

	r.GET("/ws/data", websocketapi)
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	r.Run(":8080")

	// wait for signal to exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Gracefully stop the consumer.
	consumer.Stop()
}
