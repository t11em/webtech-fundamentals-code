package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	ErrSessionExpired   = errors.New("session expired")
	ErrSessionNotFound  = errors.New("session not found")
	ErrInvalidSessionId = errors.New("invalid session id")
)

type HttpSessionManager struct {
	sessions map[string]*HttpSession
}

func NewHttpSessionManager() *HttpSessionManager {
	mgr := &HttpSessionManager{
		sessions: make(map[string]*HttpSession),
	}
	return mgr
}

func (m *HttpSessionManager) StartSession(w http.ResponseWriter) (*HttpSession, error) {
	sessionId, err := m.makeSessionId()
	if err != nil {
		return nil, err
	}

	log.Printf("start session: %s", sessionId)
	session := NewHttpSession(sessionId, 10*time.Minute)
	m.sessions[sessionId] = session
	session.SetCookie(w)

	return session, nil

}

func (m *HttpSessionManager) makeSessionId() (string, error) {
	randBytes := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, randBytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(randBytes), nil
}

func (m *HttpSessionManager) revokeSession(w http.ResponseWriter, sessionId string) {
	delete(m.sessions, sessionId)
	log.Printf("session revoked: %s", sessionId)
	if w == nil {
		return
	}
	cookie := &http.Cookie{
		Name:    CookieNameSessionId,
		MaxAge:  -1,
		Expires: time.Unix(1, 0),
	}
	http.SetCookie(w, cookie)
}

func checkSession(w http.ResponseWriter, r *http.Request) (*HttpSession, error) {
	session, err := sessionManager.GetValidSession(r)
	if err == nil {
		return session, nil
	}
	orgErr := err

	log.Printf("session check failed: %s", err.Error())
	session, err = sessionManager.StartSession(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	if r.Referer() != "" {
		page := LoginPageData{}
		page.ErrorMessage = "セッションが不正です"
		session.PageData = page
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return nil, orgErr
}

func (m *HttpSessionManager) GetValidSession(r *http.Request) (*HttpSession, error) {

	c, err := r.Cookie(CookieNameSessionId)
	if err == http.ErrNoCookie {
		return nil, ErrSessionNotFound
	}
	if err == nil {
		sessionId := c.Value
		session, err := m.getSession(sessionId)
		return session, err
	}
	return nil, err
}

func (m *HttpSessionManager) getSession(sessionId string) (*HttpSession, error) {
	if session, exists := m.sessions[sessionId]; exists {
		if time.Now().After(session.Expires) {
			delete(m.sessions, sessionId)
			return nil, ErrSessionExpired
		}
		return session, nil
	} else {
		return nil, ErrSessionExpired
	}
}

func ensureSession(w http.ResponseWriter, r *http.Request) (*HttpSession, error) {
	session, err := sessionManager.GetValidSession(r)
	if err == nil {
		return session, nil
	}

	log.Printf("session check failed: %s", err.Error())
	session, err = sessionManager.StartSession(w)
	if err != nil {
		writeInternalServerError(w, err)
		return nil, err
	}
	return session, err
}
