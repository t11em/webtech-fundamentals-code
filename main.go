package main

import (
	"errors"
	"html/template"
	"fmt"
	"log"
	"net/http"
)

var (
	sessionManager *HttpSessionManager
	accountManager *UserAccountManager
	templates *template.Template

	ErrMethodNotAllowed = errors.New("method not allowed")
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/todo", handleTodo)
	http.HandleFunc("/add", handleAdd)

	port := getPortNumber()
	fmt.Printf("listening port : %d\n", port)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("failed to start : ", err)
	}
}
