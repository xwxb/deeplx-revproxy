package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/xwxb/deeplx-revproxy/config"
	"github.com/xwxb/deeplx-revproxy/handler"
	"log"
	"net/http"
)

func main() {
	config.InitConfig()
	config.LogConfigedEndpoints()

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "DeepL Free API Proxy, derived from github.com/OwO-Network/DeepLX. Go to /translate with POST. Proxy is created by github.com/xwxb/deeplx-revproxy",
		})
	})

	r.POST("/translate", handler.ProxyHandler)

	serverPortStr := fmt.Sprintf(":%d", config.Global.Server.Port)
	err := r.Run(serverPortStr)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
