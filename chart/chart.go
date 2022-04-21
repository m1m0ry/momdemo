package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/nsqio/go-nsq"

	"github.com/m1m0ry/mom/web"
)

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
		str, _ := json.Marshal(web.Data)
		err = ws.WriteMessage(1, str)
		if err != nil {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func main() {
	// Instantiate a consumer that will subscribe to the provided channel.

	config := nsq.NewConfig()

	consumer_sign, err := nsq.NewConsumer("rand", "chart", config)
	if err != nil {
		log.Fatal(err)
	}
	consumer_mean, err := nsq.NewConsumer("mean", "chart", config)
	if err != nil {
		log.Fatal(err)
	}
	consumer_variance, err := nsq.NewConsumer("variance", "chart", config)
	if err != nil {
		log.Fatal(err)
	}
	consumer_maxmin, err := nsq.NewConsumer("maxmin", "chart", config)
	if err != nil {
		log.Fatal(err)
	}

	consumer_sign.AddHandler(&web.MyMessageHandler{})
	consumer_mean.AddHandler(&web.MeanHandler{})
	consumer_variance.AddHandler(&web.VarianceHandler{})
	consumer_maxmin.AddHandler(&web.MaxAndminHandler{})

	err = consumer_sign.ConnectToNSQLookupd("localhost:4161")
	if err != nil {
		log.Fatal(err)
	}
	err = consumer_variance.ConnectToNSQLookupd("localhost:4161")
	if err != nil {
		log.Fatal(err)
	}
	err = consumer_mean.ConnectToNSQLookupd("localhost:4161")
	if err != nil {
		log.Fatal(err)
	}
	err = consumer_maxmin.ConnectToNSQLookupd("localhost:4161")
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.LoadHTMLGlob("index.html")

	r.GET("/ws", websocketapi)
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	r.GET("/data", func(c *gin.Context) {
		str, _ := json.Marshal(web.Data)
		c.String(http.StatusOK, string(str))
	})
	r.Run(":8080")

	// wait for signal to exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Gracefully stop the consumer.
	consumer_sign.Stop()
	consumer_mean.Stop()
	consumer_variance.Stop()
	consumer_maxmin.Stop()
}
