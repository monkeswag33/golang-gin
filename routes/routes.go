package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/monkeswag33/golang-gin/global"
	logger "github.com/sirupsen/logrus"
)

func PingHandler(ctx *gin.Context) {
	ctx.String(http.StatusOK, "OK")
}

func GetHandler(ctx *gin.Context) {
	var users []global.User
	Db.Find(&users)
	ctx.JSON(http.StatusOK, users)
}

func PostHandler(ctx *gin.Context) {
	var user global.User
	if err := ctx.BindJSON(&user); err != nil {
		logger.Fatal("Error converting JSON body to struct: ", err)
	}
	Db.Create(&user)
	ctx.JSON(http.StatusCreated, user)
}

func UpdateHandler(ctx *gin.Context) {
	var user global.User
	var updates map[string]interface{}
	if err := ctx.BindJSON(&updates); err != nil {
		logger.Fatal("Error while converting JSON body to struct: ", err)
	}
	id, _ := strconv.Atoi(ctx.Param("id"))
	Db.Find(&user, id)
	Db.Model(&user).Updates(updates)
	ctx.JSON(http.StatusOK, user)
}

func DeleteHandler(ctx *gin.Context) {
	var user global.User
	id, _ := strconv.Atoi(ctx.Param("id"))
	Db.Find(&user, id)
	Db.Delete(&global.User{}, id)
	ctx.JSON(http.StatusOK, user)
}
