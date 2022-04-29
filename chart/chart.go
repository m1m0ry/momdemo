package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/m1m0ry/mom/analysis"
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
		str, _ := json.Marshal(analysis.Data)
		err = ws.WriteMessage(1, str)
		if err != nil {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func main() {

	//初始化数据接受
	go analysis.Init()

	r := gin.Default()
	r.LoadHTMLGlob("index.html")

	r.GET("/ws", websocketapi)
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	r.GET("/data", func(c *gin.Context) {
		str, _ := json.Marshal(analysis.Data)
		c.String(http.StatusOK, string(str))
	})
	r.Run(":3000")
}
