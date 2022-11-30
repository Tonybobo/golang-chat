package router

import (
	"net/http"

	"github.com/tonybobo/go-chat/pkg/global/log"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine{
	gin.SetMode("debug")

	server := gin.Default()
	server.Use(CorsMiddleware())
	server.Use(Recovery)

	socket := ServeWs

	group := server.Group("/api")
	{
		group.GET("/healthcheck" , func(c *gin.Context){
			c.JSON(http.StatusOK , gin.H{"Message" : "Server is healthy"})
		})
		group.GET("/socket.io",socket)
	}

	return server
}

func Recovery(c *gin.Context){
	defer func(){
		if r:=recover(); r!= nil{
			log.Logger.Error("Gin Error" , log.Any("error" , r))
			c.JSON(http.StatusBadGateway , gin.H{"Error" : "System Error"})
		}
	}()

	c.Next()
}

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context){
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") 
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*") 
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Logger.Error("HttpError", log.Any("HttpError", err))
			}
		}()

		c.Next()
	}
	
}