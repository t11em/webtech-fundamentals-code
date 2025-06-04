package main

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"time"
)

const cookieSessionId = "sessionId"

// セッションが開始されていることを保証する
func ensureSession(w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie(cookieSessionId)
	if err == http.ErrNoCookie {
		sessionId, err := startSession(w)
		return sessionId, err
	}
	if err == nil {
		sessionId := cookie.Value
		return sessionId, nil
	}
	return "", nil
}

// セッションを開始する
func startSession(w http.ResponseWriter) (string, error) {
	sessionId, err := makeSessionId()
	if err != nil {
		return "", err
	}

	cookie := &http.Cookie{
		Name:     cookieSessionId,
		Value:    sessionId,
		Expires:  time.Now().Add(600 * time.Second),
		HttpOnly: false,
	}
	http.SetCookie(w, cookie)
	return sessionId, nil
}

// セッションIDを発行する
func makeSessionId() (string, error) {
	randBytes := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, randBytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(randBytes), nil
}
