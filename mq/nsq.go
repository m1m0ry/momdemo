package mq

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nsqio/go-nsq"
)

type MessageQueueConfig struct {
	NsqAddr         string
	NsqLookupdAddr  string
	SupportedTopics []string
	Channel         string
}

type MessageQueue struct {
	config    MessageQueueConfig
	producer  *nsq.Producer
	consumers map[string]*nsq.Consumer
}

func NewMessageQueue(config MessageQueueConfig) (mq *MessageQueue, err error) {
	if config.NsqAddr == "" {
		config.NsqAddr = "localhost:4150"
	}
	if config.NsqLookupdAddr == "" {
		config.NsqLookupdAddr = "localhost:4161"
	}
	if config.Channel == "" {
		config.Channel = "default"
	}

	producer, err := initProducer(config.NsqAddr)
	if err != nil {
		return nil, err
	}
	consumers := make(map[string]*nsq.Consumer)
	for _, topic := range config.SupportedTopics {
		nsq.Register(topic, config.Channel)
		consumers[topic], err = initConsumer(topic, config.Channel, config.NsqAddr)
		if err != nil {
			return
		}
	}
	return &MessageQueue{
		config:    config,
		producer:  producer,
		consumers: consumers,
	}, nil
}

//连接nsq
func (mq *MessageQueue) Run() {
	for _, c := range mq.consumers {
		// c.ConnectToNSQLookupd(mq.config.NsqLookupdAddr)
		c.ConnectToNSQD(mq.config.NsqAddr)
	}
}

func initProducer(addr string) (producer *nsq.Producer, err error) {
	config := nsq.NewConfig()
	producer, err = nsq.NewProducer(addr, config)
	return
}

func initConsumer(topic string, channel string, address string) (c *nsq.Consumer, err error) {
	config := nsq.NewConfig()
	config.LookupdPollInterval = 15 * time.Second //设置重连时间
	c, err = nsq.NewConsumer(topic, channel, config)
	return
}

//同步发送消息
func (mq *MessageQueue) Pub(name string, data interface{}) (err error) {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return mq.producer.Publish(name, body)
}

//异步发送消息
func (mq *MessageQueue) PubAsync(name string, data interface{}, doneChan chan *nsq.ProducerTransaction) (err error) {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return mq.producer.PublishAsync(name, body, doneChan)
}

//消息处理函数
type Messagehandler func(v []byte)

//接受消息
func (mq *MessageQueue) Sub(topic string, handler Messagehandler) (err error) {
	v, ok := mq.consumers[topic]
	if !ok {
		err = fmt.Errorf("No such topic: " + topic)
		return
	}
	v.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		handler(message.Body)
		return nil
	}))
	return nil
}

func (mq *MessageQueue) StopProducer() {
	mq.producer.Stop()
}

func (mq *MessageQueue) StopConsumer(topic string) (err error) {
	v, ok := mq.consumers[topic]
	if !ok {
		err = fmt.Errorf("No such topic: " + topic)
		return
	}
	v.Stop()
	return nil
}
