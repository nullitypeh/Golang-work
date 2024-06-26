package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"Golangxy/config"
	"Golangxy/handlers"
	"Golangxy/utils"

	"github.com/gorilla/mux"
)

var localCache *utils.Cache

func main() {
	// 设置日志文件路径，确保 logs 文件夹与 main 文件夹同级
	logFilePath := filepath.Join("..", "logs", "app.log")

	// 创建 logs 文件夹
	if err := os.MkdirAll(filepath.Dir(logFilePath), 0755); err != nil {
		log.Fatal(err)
	}

	// 初始化 Logger
	utils.InitLogger(logFilePath)

	// 初始化数据库和 Redis
	config.InitDB()
	config.InitRedis()
	localCache = utils.NewCache()

	r := mux.NewRouter()

	r.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		defer handlePanic(w)
		handlers.CreateItemWithLock(w, r, localCache)
	}).Methods("PUT")

	r.HandleFunc("/items/{item_id}", func(w http.ResponseWriter, r *http.Request) {
		defer handlePanic(w)
		handlers.UpdateItemWithLock(w, r, localCache)
	}).Methods("POST")

	r.HandleFunc("/items/{item_id}", func(w http.ResponseWriter, r *http.Request) {
		defer handlePanic(w)
		handlers.GetItemWithCache(w, r, localCache)
	}).Methods("GET")

	r.HandleFunc("/items/{item_id}", func(w http.ResponseWriter, r *http.Request) {
		defer handlePanic(w)
		handlers.DeleteItem(w, r, localCache)
	}).Methods("DELETE")

	http.Handle("/", r)
	fmt.Println("Server is running at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handlePanic(w http.ResponseWriter) {
	if r := recover(); r != nil {
		http.Error(w, fmt.Sprint(r), http.StatusInternalServerError)
		utils.Logger.Println("Recovered from panic:", r)
	}
}
