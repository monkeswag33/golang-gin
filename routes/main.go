package routes

import (
	"context"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DbPool *gorm.DB
var Context context.Context

type User struct {
	ID        int    `json:"id"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
}

func Routes(router *gin.Engine) {
	var routerGroup *gin.RouterGroup = router.Group("/")
	routerGroup.GET("/get", GetHandler)
	routerGroup.POST("/post", PostHandler)
	routerGroup.PATCH("/update/:id", UpdateHandler)
	routerGroup.DELETE("/delete/:id", DeleteHandler)
}
