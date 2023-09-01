package rabbitmq

import (
	"douyin/models"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

type CommentMQ struct {
	RabbitMQ
	channel   *amqp.Channel
	queueName string
	exchange  string
	key       string
}

// NewCommentRabbitMQ 获取CommentMQ的对应队列。
func NewCommentRabbitMQ(queueName string) *CommentMQ {
	CommentMQ := &CommentMQ{
		RabbitMQ:  *Rmq,
		queueName: queueName,
	}
	cha, err := CommentMQ.conn.Channel()
	CommentMQ.channel = cha
	if err != nil {
		panic(fmt.Sprintf("%s:%s\n", err, "获取通道失败"))
	}
	return CommentMQ
}


var RmqCommentAdd *CommentMQ
var RmqCommentDel *CommentMQ

// InitCommentRabbitMQ 初始化rabbitMQ连接。
func InitCommentRabbitMQ() {
	RmqCommentAdd = NewCommentRabbitMQ("comment_add")
	go RmqCommentAdd.Consumer()

	RmqCommentDel = NewCommentRabbitMQ("comment_del")
	go RmqCommentDel.Consumer()
}

// Publish comment操作的发布配置。
func (l *CommentMQ) Publish(message string) {

	_, err := l.channel.QueueDeclare(
		l.queueName,false, false, false, false, nil)
	if err != nil {
		panic(fmt.Sprintf("%s:%s\n", err, "声明Comment队列失败"))
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
		panic(fmt.Sprintf("%s:%s\n", err1, "Comment队列publish失败"))
	}

}

// Consumer Comment关系的消费逻辑。
func (l *CommentMQ) Consumer() {

	_, err := l.channel.QueueDeclare(l.queueName, false, false, false, false, nil)

	if err != nil {
		panic(fmt.Sprintf("%s:%s\n", err, "声明Comment队列失败"))
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
		panic(fmt.Sprintf("%s:%s\n", err, "获取Comment消息失败"))
	}

	forever := make(chan bool)
	switch l.queueName {
	case "comment_add":
		//点赞消费队列
		go l.consumerCommentAdd(messages)
	case "comment_del":
		//取消赞消费队列
		go l.consumerCommentDel(messages)

	}

	log.Printf("[*] Waiting for messagees,To exit press CTRL+C")

	<-forever

}

//consumerCommentAdd 赞关系添加的消费方式。
func (l *CommentMQ) consumerCommentAdd(messages <-chan amqp.Delivery) {
	for d := range messages {
		// 参数解析。
		params := strings.Split(fmt.Sprintf("%s", d.Body), " ")
		Id, _ := strconv.Atoi(params[0])
		text := params[1]
		userId, _ := strconv.Atoi(params[2])
		videoId, _ := strconv.Atoi(params[3])

		comment := models.Comment{
			Model: gorm.Model{
				ID: uint(Id),
			},
			UserId:  uint(userId),
			VideoId: uint(videoId),
			Content: text,
		}

		err := models.DB.Create(&comment).Error
		if err!=nil {fmt.Println(err) }
	
	}
}

//consumerLikeDel 赞关系删除的消费方式。
func (l *CommentMQ) consumerCommentDel(messages <-chan amqp.Delivery) {
	for d := range messages {
		// 参数解析。
		params := fmt.Sprintf("%s", d.Body)
		commentId, _ := strconv.Atoi(params)

		err := models.DB.Delete(&models.Comment{}, uint(commentId)).Error
		if err!=nil {fmt.Println(err) }
	}
}
