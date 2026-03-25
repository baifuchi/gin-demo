package handler

import (
	"net/http"
	"time"

	"gin-demo/internal/config"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Publish 发布消息到队列
type Publish struct {
	Conn *amqp.Connection
	Cfg  *config.Config
}

type publishBody struct {
	Message string `json:"message"`
}

func (h *Publish) Handle(c *gin.Context) {
	var body publishBody
	if err := c.ShouldBindJSON(&body); err != nil || body.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供 JSON: {\"message\":\"内容\"}"})
		return
	}

	ch, err := h.Conn.Channel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer ch.Close()

	err = ch.PublishWithContext(c.Request.Context(), "", h.Cfg.App.QueueName, false, false, amqp.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: amqp.Persistent,
		Body:         []byte(body.Message),
		Timestamp:    time.Now(),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true, "queued": body.Message})
}
