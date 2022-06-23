package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err.Error())
	}

	PORT := os.Getenv("PORT")
	MONGO_PORT := os.Getenv("MONGO_PORT")
	MONGO_HOST := os.Getenv("MONGO_HOST")
	REDIS_PORT := os.Getenv("REDIS_PORT")
	REDIS_HOST := os.Getenv("REDIS_HOST")

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if _, err = mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", MONGO_HOST, MONGO_PORT))); err != nil {
		log.Fatalln("Failed connect to MongoDB ", err.Error())
	}
	log.Println("MongoDB Connected")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", REDIS_HOST, REDIS_PORT),
		Password: "",
		DB:       0,
	})

	result, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(result)

	if err = rdb.Set(ctx, "key", "value", 0).Err(); err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	log.Println("key", val)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello World!")
	})

	log.Println("Server running on localhost:5000")
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%v", PORT), nil))
}
