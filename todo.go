package main

import (
	"html"
	"html/template"
	"net/http"
	"strings"
)

var todoLists = make(map[string][]string)

func getTodoList(sessionId string) []string {
	todos, ok := todoLists[sessionId]
	if !ok {
		todos = []string{}
		todoLists[sessionId] = todos
	}
	return todos
}

func handleTodo(w http.ResponseWriter, r *http.Request) {
	sessionId, err := ensureSession(w, r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	todos := getTodoList(sessionId)
	t, _ := template.ParseFiles("templates/todo.html")
	t.Execute(w, todos)
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	sessionId, err := ensureSession(w, r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	todos := getTodoList(sessionId)

	r.ParseForm()
	todo := strings.TrimSpace(html.EscapeString(r.Form.Get("todo")))
	if todo != "" {
		todoLists[sessionId] = append(todos, todo)
	}
	http.Redirect(w, r, "/todo", 303)
}
