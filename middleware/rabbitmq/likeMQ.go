package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
	"fmt"
	"strings"
	"strconv"
	"douyin/models"
)

type LikeMQ struct {
	RabbitMQ
	channel   *amqp.Channel
	queueName string
	exchange  string
	key       string
}

// NewLikeRabbitMQ 获取likeMQ的对应队列。
func NewLikeRabbitMQ(queueName string) *LikeMQ {
	likeMQ := &LikeMQ{
		RabbitMQ:  *Rmq,
		queueName: queueName,
	}
	cha, err := likeMQ.conn.Channel()
	likeMQ.channel = cha
	if err != nil {
		panic(fmt.Sprintf("%s:%s\n", err, "获取通道失败"))
	}
	return likeMQ
}


var RmqLikeAdd *LikeMQ
var RmqLikeDel *LikeMQ

// InitLikeRabbitMQ 初始化rabbitMQ连接。
func InitLikeRabbitMQ() {
	RmqLikeAdd = NewLikeRabbitMQ("like_add")
	go RmqLikeAdd.Consumer()

	RmqLikeDel = NewLikeRabbitMQ("like_del")
	go RmqLikeDel.Consumer()
}

// Publish like操作的发布配置。
func (l *LikeMQ) Publish(message string) {

	_, err := l.channel.QueueDeclare(
		l.queueName,false, false, false, false, nil)
	if err != nil {
		panic(fmt.Sprintf("%s:%s\n", err, "声明like队列失败"))
	}

	err1 := l.channel.Publish(
		l.exchange,
		l.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err1 != nil {
		panic(fmt.Sprintf("%s:%s\n", err1, "like队列publish失败"))
	}

}

// Consumer like关系的消费逻辑。
func (l *LikeMQ) Consumer() {

	_, err := l.channel.QueueDeclare(l.queueName, false, false, false, false, nil)

	if err != nil {
		panic(fmt.Sprintf("%s:%s\n", err, "声明like队列失败"))
	}

	//2、接收消息
	messages, err1 := l.channel.Consume(
		l.queueName,	//队列名
		"",			//消费者名，用来区分多个消费者，以实现公平分发或均等分发策略
		true,		//是否自动应答
		false,		//是否具有排他性
		false,		//是否接收同一个连接中的消息，若为true，则只能接收别的conn中发送的消息
		false,		//消息队列是否阻塞
		nil,
	)
	if err1 != nil {
		panic(fmt.Sprintf("%s:%s\n", err, "获取like消息失败"))
	}

	forever := make(chan bool)
	switch l.queueName {
	case "like_add":
		//点赞消费队列
		go l.consumerLikeAdd(messages)
	case "like_del":
		//取消赞消费队列
		go l.consumerLikeDel(messages)

	}

	log.Printf("[*] Waiting for messagees,To exit press CTRL+C")

	<-forever

}

//consumerLikeAdd 赞关系添加的消费方式。
func (l *LikeMQ) consumerLikeAdd(messages <-chan amqp.Delivery) {
	for d := range messages {
		// 参数解析。
		params := strings.Split(fmt.Sprintf("%s", d.Body), " ")
		userId, _ := strconv.Atoi(params[0])
		videoId, _ := strconv.Atoi(params[1])

		user := models.User{}
		user.ID = uint(userId)
		video := models.Video{}
		video.ID = uint(videoId)
		err := models.DB.Model(&user).Association("LikeVideo").Append(&video)
		if err!=nil {fmt.Println(err) }
	
	}
}

//consumerLikeDel 赞关系删除的消费方式。
func (l *LikeMQ) consumerLikeDel(messages <-chan amqp.Delivery) {
	for d := range messages {
		// 参数解析。
		params := strings.Split(fmt.Sprintf("%s", d.Body), " ")
		userId, _ := strconv.Atoi(params[0])
		videoId, _ := strconv.Atoi(params[1])

		user := models.User{}
		user.ID = uint(userId)
		video := models.Video{}
		video.ID = uint(videoId)
		err := models.DB.Model(&user).Association("LikeVideo").Delete(&video)
		if err!=nil {fmt.Println(err) }
	}
}
