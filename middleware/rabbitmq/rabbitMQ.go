package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"douyin/config"
)

// const MQURL = "amqp://guest:guest@127.0.0.1:5672/"

type RabbitMQ struct {
	conn  *amqp.Connection
	mqurl string
}

var Rmq *RabbitMQ

// InitRabbitMQ 初始化RabbitMQ的连接和通道。
func InitRabbitMQ() error {
	MQURL:= fmt.Sprintf("amqp://%s:%s@%s:%d/",
		config.GlobalConfig.Rabbitmq.Username,
		config.GlobalConfig.Rabbitmq.Password,
		config.GlobalConfig.Rabbitmq.Host,
		config.GlobalConfig.Rabbitmq.Port)
	Rmq = &RabbitMQ{
		mqurl: MQURL,
	}
	dial, err := amqp.Dial(Rmq.mqurl)
	if err != nil {
		return err
	}
	Rmq.conn = dial

	InitFollowRabbitMQ()
	InitLikeRabbitMQ()
	InitCommentRabbitMQ()
	return nil

}

// 关闭mq通道和mq的连接。
func (r *RabbitMQ) destroy() {
	r.conn.Close()
}