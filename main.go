package main

import (
	"log"
	"os"

	"gin-demo/internal/config"
	"gin-demo/internal/router"
	rmq "gin-demo/internal/rabbitmq"
	"gin-demo/internal/worker"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("加载配置: %v", err)
	}

	gin.SetMode(cfg.Server.Mode)

	if err := os.MkdirAll(cfg.App.DataDir, 0o755); err != nil {
		log.Fatalf("创建数据目录: %v", err)
	}

	conn, err := rmq.Dial(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	worker.StartConsumers(conn, cfg)

	r := gin.Default()
	router.Register(r, conn, cfg)

	log.Printf("监听 %s，POST /publish {\"message\":\"...\"}，消费者写入 %s/consumer_*.log",
		cfg.Server.Addr, cfg.App.DataDir)
	if err := r.Run(cfg.Server.Addr); err != nil {
		log.Fatal(err)
	}
}
