package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"fizzbuzz/internal/controller"
	"fizzbuzz/internal/core"
	"fizzbuzz/internal/repo"
)

type FizzBuzzRequest struct {
	Numbers []int    `json:"numbers" binding:"required,dive,gt=0"`
	Value   []string `json:"value" binding:"required,dive,notblank"`
	Limit   int      `json:"limit" binding:"required,gt=0"`
}

func main() {
	cfg := loadEnv()

	mongoDB := initMongoDB(cfg.MongoURI, cfg.MongoDB)
	repository := repo.NewMongoRepo[mongo.Collection]("fizzbuzz_records", &mongoDB)

	core := core.NewCore[mongo.Collection](repository)
	controller := controller.NewFizzBuzzController(core)

	r := gin.Default()

	r.POST("/fizzbuzz", controller.FizzBuzz)
	r.GET("/stats", controller.Stats)

	r.Run(":8080")
}

func initMongoDB(uri, dbName string) mongo.Collection {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	return *client.Database(dbName).Collection("fizzbuzz_records")
}

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	MongoURI   string
	MongoDB    string
}

func loadEnv() Config {
	return Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		MongoURI:   os.Getenv("MONGO_URI"),
		MongoDB:    os.Getenv("MONGO_DB"),
	}
}
