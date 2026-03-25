package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 应用配置（对应 configs/config.yaml）
type Config struct {
	Server   Server   `yaml:"server"`
	RabbitMQ RabbitMQ `yaml:"rabbitmq"`
	App      App      `yaml:"app"`
}

type Server struct {
	Addr string `yaml:"addr"`
	Mode string `yaml:"mode"`
}

type RabbitMQ struct {
	Host      string `yaml:"host"`
	VHost     string `yaml:"vhost"`
	Port      int    `yaml:"port"`
	Login     string `yaml:"login"`
	Password  string `yaml:"password"`
	Heartbeat int    `yaml:"heartbeat"` // 秒
}

type App struct {
	QueueName     string `yaml:"queue_name"`
	DataDir       string `yaml:"data_dir"`
	ConsumerCount int    `yaml:"consumer_count"`
}

// HeartbeatDuration 将配置中的秒转为 amqp 心跳间隔
func (r RabbitMQ) HeartbeatDuration() time.Duration {
	if r.Heartbeat <= 0 {
		return 10 * time.Second
	}
	return time.Duration(r.Heartbeat) * time.Second
}

// Load 从 YAML 文件加载；路径由 CONFIG_PATH 指定，默认 configs/config.yaml
func Load() (*Config, error) {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "configs/config.yaml"
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置 %s: %w", path, err)
	}
	var c Config
	if err := yaml.Unmarshal(raw, &c); err != nil {
		return nil, fmt.Errorf("解析配置: %w", err)
	}
	if err := c.validate(); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *Config) validate() error {
	if c.Server.Addr == "" {
		c.Server.Addr = ":8080"
	}
	if c.App.QueueName == "" {
		return fmt.Errorf("app.queue_name 不能为空")
	}
	if c.App.DataDir == "" {
		return fmt.Errorf("app.data_dir 不能为空")
	}
	if c.App.ConsumerCount < 1 {
		c.App.ConsumerCount = 1
	}
	if c.RabbitMQ.Port == 0 {
		c.RabbitMQ.Port = 5672
	}
	return nil
}
