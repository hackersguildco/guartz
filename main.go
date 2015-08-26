package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/infiniteloopsco/guartz/endpoint"
	"github.com/infiniteloopsco/guartz/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	initEnv()
	models.InitDB()
	models.InitCron()

	router := endpoint.GetMainEngine()
	port := os.Getenv("PORT")

	s := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	s.ListenAndServe()
}

func initEnv() {
	if err := godotenv.Load(".env_dev"); err != nil {
		log.Fatal("Error loading .env_dev file")
	}
}
