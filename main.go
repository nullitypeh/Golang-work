package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/yourusername/product-management/config"
	"github.com/yourusername/product-management/handlers"
)

func main() {
	config.InitDB()
	config.InitRedis()

	r := mux.NewRouter()
	r.HandleFunc("/items", handlers.GetItems).Methods("GET")
	r.HandleFunc("/items", handlers.CreateItem).Methods("POST")
	r.HandleFunc("/items/{id}", handlers.GetItem).Methods("GET")
	r.HandleFunc("/items/{id}", handlers.UpdateItem).Methods("PUT")
	r.HandleFunc("/items/{id}", handlers.DeleteItem).Methods("DELETE")

	http.Handle("/", r)
	fmt.Println("Server is running at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
