package main

import (
	"net/http"
	"time"
)

const cookieNameSessionId = "sessionId"

type HttpSession struct {
	SessionId   string
	Expires     time.Time
	PageData    any
	UserAccount *UserAccount
}

func NewHttpSession(sessionId string, validitidyTime time.Duration) *HttpSession {
	session := &HttpSession{
		SessionId: sessionId,
		Expires:   time.Now().Add(validitidyTime),
		PageData:  "",
	}
	return session
}

func (s *HttpSession) ClearPageData() {
	s.PageData = ""
}

func (s HttpSession) SetCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     cookieNameSessionId,
		Value:    s.SessionId,
		Expires:  s.Expires,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}
