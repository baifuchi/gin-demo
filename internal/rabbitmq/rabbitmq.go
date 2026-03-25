package rabbitmq

import (
	"fmt"
	"net/url"

	"gin-demo/internal/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Dial 建立连接并声明队列
func Dial(cfg *config.Config) (*amqp.Connection, error) {
	uri := amqpURI(cfg)
	conn, err := amqp.DialConfig(uri, amqp.Config{
		Heartbeat: cfg.RabbitMQ.HeartbeatDuration(),
		Locale:    "en_US",
	})
	if err != nil {
		return nil, fmt.Errorf("连接 RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("打开通道: %w", err)
	}
	_, err = ch.QueueDeclare(
		cfg.App.QueueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,
	)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("声明队列: %w", err)
	}
	if err := ch.Close(); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("关闭声明通道: %w", err)
	}
	return conn, nil
}

func amqpURI(cfg *config.Config) string {
	r := cfg.RabbitMQ
	vh := r.VHost
	if vh == "" || vh == "/" {
		vh = "%2F"
	} else {
		vh = url.PathEscape(vh)
	}
	return fmt.Sprintf("amqp://%s:%s@%s:%d/%s", r.Login, r.Password, r.Host, r.Port, vh)
}
