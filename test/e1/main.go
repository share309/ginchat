package main

import (
	"github.com/gin-gonic/gin"
	handler "github.com/share309/ginchat"
)

func main() {
	router := gin.Default()
	router.GET("/chat", handler.Chat)
	router.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
