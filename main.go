package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/monkeswag33/golang-gin/global"
	"github.com/monkeswag33/golang-gin/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func initDB() *gorm.DB {
	var databaseUri string = os.Getenv("POSTGRES_URI")
	if len(databaseUri) == 0 {
		log.Fatal("Could not find POSTGRES_URI environment variable")
	}
	fmt.Println("Found POSTGRES_URI: " + databaseUri)
	var newLogger logger.Interface = logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	db, err := gorm.Open(postgres.Open(databaseUri), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal("Error while connecting to database")
	}
	db.AutoMigrate(&global.User{})
	return db
}

func SetupRouter() (*gin.Engine, *gorm.DB, string) {
	var db *gorm.DB = initDB()
	routes.Db = db
	var PORT string = os.Getenv("PORT")
	if PORT == "" {
		log.Println("Could not find PORT, using default port 8080")
		PORT = "8080"
	} else {
		fmt.Println("Found PORT: " + PORT)
	}
	var router *gin.Engine = gin.New()
	if os.Getenv("GIN_MODE") == "debug" || os.Getenv("GIN_MODE") == "" {
		router.Use(gin.Logger())
	}
	router.Use(gin.Recovery())
	router.SetTrustedProxies([]string{"0.0.0.0"})
	routes.Routes(router)
	return router, db, fmt.Sprintf(":%s", PORT)
}

func main() {
	router, _, PORT := SetupRouter()
	router.Run(PORT)
}
