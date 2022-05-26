package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var db *redis.Client
var ctx = context.Background()

func main() {
	port := DataBase()
	http.HandleFunc("/", requestHandler)
	http.HandleFunc("/favicon.ico", loadFavicon)
	http.ListenAndServe(":"+port, nil)
	fmt.Println("Server started on port:", port)

}

func DataBase() string {
	server, password, DB, port := loadVars()

	Db, _ := strconv.Atoi(DB)
	db = redis.NewClient(&redis.Options{
		Addr:     server,
		Password: password,
		DB:       Db,
	})

	fmt.Println("Connected to Database")
	return port
}

func loadVars() (string, string, string, string) {
	if _, err := os.Stat(".env"); err == nil {
		godotenv.Load()
	}
	port := os.Getenv("PORT")
	RedisServer := os.Getenv("REDIS_SERVER")
	RedisPassword := os.Getenv("REDIS_PASSWORD")
	RedisDatabase := os.Getenv("REDIS_DATABASE")

	if RedisServer == "" {
		panic("REDIS_SERVER is not set")
	}
	if RedisPassword == "" {
		panic("REDIS_PASSWORD is not set")
	}
	if RedisDatabase == "" {
		RedisDatabase = "0"
	}
	if port == "" {
		port = "8080"
	}

	return RedisServer, RedisPassword, RedisDatabase, port

}
func requestHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len("/"):]

	if r.URL.Path == "/" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	value, err := db.Get(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	if value == "" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, value, http.StatusFound)
}
func loadFavicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "favicon.ico")
}
