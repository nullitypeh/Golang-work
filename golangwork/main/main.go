package main

import (
	"fmt"
	"log"
	"net/http"

	"Golangxy/config"
	"Golangxy/handlers"
	"Golangxy/utils"

	"github.com/gorilla/mux"
)

var localCache *utils.Cache

func main() {
	config.InitDB()
	config.InitRedis()
	localCache = utils.NewCache()

	r := mux.NewRouter()

	r.HandleFunc("/items", handlers.CreateItemWithLock).Methods("PUT")
	r.HandleFunc("/items/{item_id}", handlers.UpdateItemWithLock).Methods("POST")
	r.HandleFunc("/items/{item_id}", GetItemWithCache).Methods("GET") // 使用新的处理程序函数
	r.HandleFunc("/items/{item_id}", DeleteItem).Methods("DELETE")    // 添加 DELETE 请求的处理程序

	http.Handle("/", r)
	fmt.Println("Server is running at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
func GetItemWithCache(w http.ResponseWriter, r *http.Request) {
	handlers.GetItemWithCache(w, r, localCache)
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	handlers.DeleteItem(w, r, localCache)
}
