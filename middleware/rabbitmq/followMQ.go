package rabbitmq

import(
	"github.com/streadway/amqp"
	"log"
	"fmt"
	"strings"
	"strconv"
	"douyin/models"
)

type FollowMQ struct {
	RabbitMQ
	channel   *amqp.Channel
	queueName string	//队列名称
	exchange  string	//交换机
	key       string	//routing Key
}

// NewFollowRabbitMQ 获取followMQ的对应队列。
func NewFollowRabbitMQ(queueName string) *FollowMQ {
	followMQ := &FollowMQ{
		RabbitMQ:  *Rmq,
		queueName: queueName,
	}

	cha, err := followMQ.conn.Channel()
	followMQ.channel = cha
	if err != nil {
		panic(fmt.Sprintf("%s:%s\n", err, "获取通道失败"))
	}
	return followMQ
}

// 关闭mq通道和mq的连接。
func (f *FollowMQ) destroy() {
	f.channel.Close()
}

var RmqFollowAdd *FollowMQ
var RmqFollowDel *FollowMQ

// InitFollowRabbitMQ 初始化rabbitMQ连接。
func InitFollowRabbitMQ() {
	RmqFollowAdd = NewFollowRabbitMQ("follow_add")
	go RmqFollowAdd.Consumer()

	RmqFollowDel = NewFollowRabbitMQ("follow_del")
	go RmqFollowDel.Consumer()
}


// Publish生产者  follow关系的发布配置。
func (f *FollowMQ) Publish(message string) {
	//1、声明队列
	_, err := f.channel.QueueDeclare(
		f.queueName, // 队列名
		false,	//是否持久化
		false,	//是否为自动删除
		false,	//是否具有排他性
		false,	//是否阻塞
		nil,	//额外属性
	)
	if err != nil {
		panic(fmt.Sprintf("%s:%s\n", err, "声明关注队列失败"))
	}
	//2、发送消息
	err1 :=f.channel.Publish(
		f.exchange,
		f.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err1 != nil {
		panic(fmt.Sprintf("%s:%s\n", err1, "follow队列publish失败"))
	}
}

// Consumer消费者  follow关系的消费逻辑。
func (f *FollowMQ) Consumer() {
	// 1、声明队列（生产者和消费者 两端都要声明）
	_, err := f.channel.QueueDeclare(f.queueName, false, false, false, false, nil)

	if err != nil {
		panic(fmt.Sprintf("%s:%s\n", err, "声明关注队列失败"))
	}

	//2、从队列接收消息
	msgs, err := f.channel.Consume(
		f.queueName,	//队列名
		"",			//消费者名，用来区分多个消费者，以实现公平分发或均等分发策略
		true,		//是否自动应答
		false,		//是否具有排他性
		false,		//是否接收同一个连接中的消息，若为true，则只能接收别的conn中发送的消息
		false,		//消息队列是否阻塞
		nil,
	)
	if err != nil {
		panic(fmt.Sprintf("%s:%s\n", err, "获取关注消息失败"))
	}

	forever := make(chan bool)
	switch f.queueName {
	case "follow_add":
		go f.consumerFollowAdd(msgs)
	case "follow_del":
		go f.consumerFollowDel(msgs)

	}

	log.Printf("[*] Waiting for messagees,To exit press CTRL+C")

	<-forever

}

// 消费者 添加follow关系的具体实现
func (f *FollowMQ) consumerFollowAdd(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		// 参数解析。
		params := strings.Split(fmt.Sprintf("%s", d.Body), " ")
		toId, _ := strconv.Atoi(params[0])
		fromId, _ := strconv.Atoi(params[1])
		// 日志记录。
		fmt.Printf("CALL FollowAction(%v,%v)", fromId,toId)
		//执行FollowAction关注操作
		relation := models.Relation{
			FollowedId:   uint(toId),
			FollowerId: uint(fromId),
		}
		if err := models.DB.Table("user_follows").Create(&relation).Error; err != nil { //创建记录
			fmt.Println(err) 
		}
	}
}

// 关系删除的消费方式。
func (f *FollowMQ) consumerFollowDel(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		// 参数解析。
		params := strings.Split(fmt.Sprintf("%s", d.Body), " ")
		fromId, _ := strconv.Atoi(params[0])
		toId, _ := strconv.Atoi(params[1])
		// 日志记录。
		fmt.Printf("CALL delFollowRelation(%v,%v)", fromId, toId)
		user1 := models.User{}
		user1.ID = uint(fromId)
		user2 := models.User{}
		user2.ID = uint(toId)
		err := models.DB.Model(&user1).Association("Follow").Delete(user2)
		fmt.Println(err) 
	}
}


