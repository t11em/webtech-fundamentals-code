package main

import (
	"log"
	"net/http"
)

type LoginPageData struct {
	UserId       string
	ErrorMessage string
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	session, err := ensureSession(w, r)
	if err != nil {
		return
	}

	switch r.Method {
	case http.MethodGet:
		showLogin(w, r, session)
		return
	case http.MethodPost:
		login(w, r, session)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func showLogin(w http.ResponseWriter, r *http.Request, session *HttpSession) {
	var pageData LoginPageData
	if p, ok := session.PageData.(LoginPageData); ok {
		pageData = p
	} else {
		pageData = LoginPageData{}
	}
	templates.ExecuteTemplate(w, "login.html", pageData)
	session.ClearPageData()
}

func login(w http.ResponseWriter, r *http.Request, session *HttpSession) {
	r.ParseForm()
	userId := r.Form.Get("userId")
	password := r.Form.Get("password")

	log.Printf("login attempt: userId=%s\n", userId)
	account, err := accountManager.Authenticate(userId, password)
	if err != nil {
		log.Printf("login failed: %s\n", err.Error())
		session.PageData = LoginPageData{
			ErrorMessage: "ログインに失敗しました",
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	session.UserAccount = account

	log.Printf("login success: userId=%s\n", userId)
	http.Redirect(w, r, "/todo", http.StatusSeeOther)
	return
}
