package main

import (
	"html/template"
	"log"
	"net/http"
)

var todoList []string

func handleTodo(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/todo.html")
	t.Execute(w, todoList)
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	todo := r.Form.Get("todo")
	todoList = append(todoList, todo)
	handleTodo(w, r)
}

func main() {
	todoList = append(todoList, "洗顔", "朝食", "歯磨き")

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/todo", handleTodo)
	http.HandleFunc("/add", handleAdd)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("failed to start : ", err)
	}
}
