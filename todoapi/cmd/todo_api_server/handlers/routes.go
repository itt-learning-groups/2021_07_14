package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter(tl TodoListHandler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", Home).Methods("GET")
	r.HandleFunc("/echo", tl.EchoPost).Methods("POST")
	r.HandleFunc("/echo/{param}", tl.EchoPut).Methods("PUT")
	r.HandleFunc("/todo", tl.Create).Methods("POST")
	r.HandleFunc("/list", tl.GetAll).Methods("GET")
	r.HandleFunc("/todo/{name}", tl.Edit).Methods("PUT")
	r.HandleFunc("/todo/{name}", tl.Delete).Methods("DELETE")
	r.HandleFunc("/clear", tl.DeleteAll).Methods("DELETE")

	return r
}

func Home(w http.ResponseWriter, r *http.Request) {
	msg := "Welcome to the Learn and DevOps ToDo App!"
	data, err := json.Marshal(msg)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Context-Type", "application/json")
	_, err = w.Write(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
