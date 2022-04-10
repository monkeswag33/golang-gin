package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var Db *gorm.DB

func Routes(router *gin.Engine) {
	var routerGroup *gin.RouterGroup = router.Group("/")
	routerGroup.GET("/ping", PingHandler)
	routerGroup.GET("/get", GetHandler)
	routerGroup.POST("/post", PostHandler)
	routerGroup.PATCH("/update/:id", UpdateHandler)
	routerGroup.DELETE("/delete/:id", DeleteHandler)
}
