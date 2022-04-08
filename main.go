package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/monkeswag33/golang-gin/routes"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	var databaseUri string = os.Getenv("POSTGRES_URI")
	if len(databaseUri) == 0 {
		log.Fatal("Could not find POSTGRES_URI environment variable")
	}
	fmt.Println("Found POSTGRES_URI: " + databaseUri)
	db, err := gorm.Open(postgres.Open(databaseUri), &gorm.Config{})
	if err != nil {
		log.Fatal("Error while connecting to database")
	}
	return db
}

func SetupRouter() (*gin.Engine, *pgxpool.Pool, string) {
	var context context.Context = context.Background()
	var db *gorm.DB = initDB()
	routes.DbPool = db
	routes.Context = context
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
	return router, dbPool, fmt.Sprintf(":%s", PORT)
}

type User struct {
	ID        int
	Firstname string
	Lastname  string
}

func main() {
	// router, pool, PORT := SetupRouter()
	// defer pool.Close()
	// router.Run(PORT)
}
