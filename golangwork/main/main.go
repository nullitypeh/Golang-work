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
	localCache = utils.NewCache() // 初始化全局变量

	r := mux.NewRouter()

	r.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateItemWithLock(w, r, localCache)
	}).Methods("PUT")

	r.HandleFunc("/items/{item_id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateItemWithLock(w, r, localCache)
	}).Methods("POST")

	r.HandleFunc("/items/{item_id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetItemWithCache(w, r, localCache)
	}).Methods("GET")

	r.HandleFunc("/items/{item_id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteItem(w, r, localCache)
	}).Methods("DELETE")

	http.Handle("/", r)
	fmt.Println("Server is running at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
