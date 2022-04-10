package main

import (
	"fmt"
	"os"

	"github.com/monkeswag33/golang-gin/global"
	"github.com/monkeswag33/golang-gin/routes"

	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getenv(key string, defaultkey string) string {
	var value string = os.Getenv(key)
	if len(value) == 0 {
		logger.Warnf("Could not find %s, defaulting to '%s'", key, defaultkey)
		value = defaultkey
	} else {
		logger.Infof("Found %s: '%s'", key, value)
	}
	return value
}

// Parse POSTGRES_HOST, POSTGRES_PORT, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DATABASE
func individualParams() string {
	var connectionString string
	connectionString += "host=" + getenv("POSTGRES_HOST", "localhost") + " "
	connectionString += "port=" + getenv("POSTGRES_PORT", "5432") + " "
	connectionString += "user=" + getenv("POSTGRES_USER", "postgres") + " "
	connectionString += "password=" + getenv("POSTGRES_PASSWORD", "postgres") + " "
	connectionString += "database=" + getenv("POSTGRES_DATABASE", "postgres")
	return connectionString
}

func initDB() *gorm.DB {
	var databaseUri string = os.Getenv("POSTGRES_URI")
	if len(databaseUri) == 0 {
		logger.Warn("Could not find POSTGRES_URI environment variable, trying POSTGRES_HOST, POSTGRES_PORT, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DATABASE")
		databaseUri = individualParams()
	} else {
		logger.Info("Found POSTGRES_URI: ", databaseUri)
	}
	logger.Info("Connecting with: ", databaseUri)
	db, err := gorm.Open(postgres.Open(databaseUri), &gorm.Config{})
	if err != nil {
		logger.Fatal("Error while connecting to database: ", err)
	}
	db.AutoMigrate(&global.User{})
	return db
}

func SetupRouter() (*gin.Engine, *gorm.DB, string) {
	var db *gorm.DB = initDB()
	routes.Db = db
	var PORT string = os.Getenv("PORT")
	if PORT == "" {
		logger.Warn("Could not find PORT, using default port 8080")
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

func setLogLevel() {
	level, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		level = "debug"
	}
	ll, err := logger.ParseLevel(level)
	if err != nil {
		logger.Warn(err)
		ll = logger.DebugLevel
	}
	logger.SetLevel(ll)
}

func main() {
	setLogLevel()
	router, _, PORT := SetupRouter()
	router.Run(PORT)
}
