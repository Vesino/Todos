package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Todo struct {
	Todo        string `json:"todo"`
	Description string `json:"description"`
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
}

func ListTodos(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	todos := []Todo{{"queacer", "Hacer queacer"}, {"Comprar cosas", "ir al super a comprar cosas"}}
	bytes, err := json.Marshal(todos)
	if err != nil {
		log.Println("Not possible marshall todos", err)
		return
	}
	w.Write(bytes)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/todos", ListTodos)

	log.Fatal(http.ListenAndServe(":3333", mux))
}
