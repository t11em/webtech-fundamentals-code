package main

import (
	"log"
	"net/http"
)

type CreateUserAccountPageData struct {
	ErrorMessage string
}

func handleCreateUserAccount(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		session, err := ensureSession(w, r)
		if err != nil {
			return
		}
		showCreateUserAccount(w, session)
		return
	case http.MethodPost:
		session, err := ensureSession(w, r)
		if err != nil {
			return
		}
		CreateNewUserAccount(w, r, session)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func showCreateUserAccount(w http.ResponseWriter, session *HttpSession) {
	pageData := CreateUserAccountPageData{}

	if p, ok := session.PageData.(CreateUserAccountPageData); ok {
		pageData.ErrorMessage = p.ErrorMessage
	}
	templates.ExecuteTemplate(w, "create-user-account.html", pageData)
	session.ClearPageData()
}

func CreateNewUserAccount(w http.ResponseWriter, r *http.Request, session *HttpSession) {
	r.ParseForm()
	userId := r.Form.Get("userId")
	password := MakePassowrd()

	user, err := accountManager.NewUserAccount(userId, password)
	if err != nil {
		pageData := CreateUserAccountPageData{}
		log.Printf("create user failed : userId=%s, err=%s", userId, err)
		if err == ErrUserAlreadyExists {
			pageData.ErrorMessage = "すでに使われているユーザIDです"
		} else if err == ErrInvalidUserIdFormat {
			pageData.ErrorMessage = "ユーザIDの形式が不正です"
		} else {
			pageData.ErrorMessage = err.Error()
		}
		session.PageData = pageData
		http.Redirect(w, r, "/create-user-account", http.StatusSeeOther)
		return
	}

	session.PageData = NewUserAccountPageData{
		UserId:   user.Id,
		Password: password,
		Expires:  user.ExpiresText(),
	}
	http.Redirect(w, r, "/user-account", http.StatusSeeOther)
	return
}
