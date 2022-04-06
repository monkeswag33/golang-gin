package main

import (
	"context"
	"fmt"
	"ishank/webserver/src/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

func initDB(context context.Context) *pgxpool.Pool {
	var databaseUri string = os.Getenv("POSTGRES_URI")
	if len(databaseUri) == 0 {
		log.Fatal("Could not find POSTGRES_URI environment variable")
	}
	fmt.Println("Found POSTGRES_URI: " + databaseUri)
	dbPool, err := pgxpool.Connect(context, databaseUri)
	if err != nil {
		log.Fatal("Error while connecting to database")
	}
	return dbPool
}

func main() {
	var context context.Context = context.Background()
	var dbPool *pgxpool.Pool = initDB(context)
	defer dbPool.Close()
	routes.DbPool = dbPool
	routes.Context = context

	var PORT string = os.Getenv("PORT")
	if PORT == "" {
		log.Println("Could not find PORT, using default port 8080")
		PORT = "8080"
	}
	var router *gin.Engine = gin.New()
	if os.Getenv("GIN_MODE") == "debug" || os.Getenv("GIN_MODE") == "" {
		router.Use(gin.Logger())
	}
	router.Use(gin.Recovery())
	router.SetTrustedProxies([]string{"0.0.0.0"})
	routes.Routes(router)
	router.Run(fmt.Sprintf(":%s", PORT))
}
