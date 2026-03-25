package router

import (
	"gin-demo/internal/config"
	"gin-demo/internal/handler"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Register 注册路由与依赖
func Register(r *gin.Engine, conn *amqp.Connection, cfg *config.Config) {
	r.GET("/hello", handler.Hello)

	pub := &handler.Publish{Conn: conn, Cfg: cfg}
	r.POST("/publish", pub.Handle)
}
