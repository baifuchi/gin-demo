package worker

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gin-demo/internal/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

// StartConsumers 启动多个协程消费者，各自写入 data/consumer_<id>.log
func StartConsumers(conn *amqp.Connection, cfg *config.Config) {
	for i := range cfg.App.ConsumerCount {
		go runConsumer(conn, cfg, i)
	}
}

func runConsumer(conn *amqp.Connection, cfg *config.Config, id int) {
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("消费者 %d 打开通道失败: %v", id, err)
		return
	}
	defer ch.Close()

	if err := ch.Qos(1, 0, false); err != nil {
		log.Printf("消费者 %d Qos: %v", id, err)
		return
	}

	tag := fmt.Sprintf("gin-worker-%d", id)
	msgs, err := ch.Consume(cfg.App.QueueName, tag, false, false, false, false, nil)
	if err != nil {
		log.Printf("消费者 %d 订阅失败: %v", id, err)
		return
	}

	outPath := filepath.Join(cfg.App.DataDir, fmt.Sprintf("consumer_%d.log", id))
	f, err := os.OpenFile(outPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		log.Printf("消费者 %d 打开文件 %s: %v", id, outPath, err)
		return
	}
	defer f.Close()

	log.Printf("消费者 %d 就绪，写入 %s", id, outPath)

	for d := range msgs {
		line := fmt.Sprintf("[%s] %s\n", time.Now().Format(time.RFC3339), string(d.Body))
		if _, err := f.WriteString(line); err != nil {
			log.Printf("消费者 %d 写文件: %v", id, err)
			_ = d.Nack(false, true)
			continue
		}
		if err := f.Sync(); err != nil {
			log.Printf("消费者 %d sync: %v", id, err)
		}
		if err := d.Ack(false); err != nil {
			log.Printf("消费者 %d Ack: %v", id, err)
		}
	}
}
