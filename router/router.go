package router

import (
	"github.com/gin-gonic/gin"
)

var ginEngine *gin.Engine
func init() {
	ginEngine = gin.New()
}

func InitRouter() *gin.Engine  {
	ginEngine.GET("/ping", func(c *gin.Context) {

		//


		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return ginEngine
}
