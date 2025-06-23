package main

import (
	"errors"
	"log"
	"math/rand"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists   = errors.New("user account already exists")
	ErrInvalidUserIdFormat = errors.New("invalid user id format")
	ErrLoginFailed         = errors.New("login failed")
	RegexAccountId         = regexp.MustCompile(`^[A-Za-z0-9_+@-]{1,32}$`)
)

type UserAccountManager struct {
	userAccounts map[string]*UserAccount
	location     *time.Location
}

func NewUserAccountManager() *UserAccountManager {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic("failed to load Asia/Tokyo location: " + err.Error())
	}

	manager := &UserAccountManager{
		userAccounts: make(map[string]*UserAccount),
		location:     jst,
	}
	return manager
}

func (m *UserAccountManager) ValidateUserId(userId string) bool {
	return RegexAccountId.MatchString(userId)
}

func (m *UserAccountManager) NewUserAccount(userId, password string) (*UserAccount, error) {
	if !m.ValidateUserId(userId) {
		return nil, ErrInvalidUserIdFormat
	}
	_, exists := m.GetUserAccount(userId)
	if exists {
		return nil, ErrUserAlreadyExists
	}
	expires := time.Now().In(m.location).Add(time.Minute * UserAccountLimitInMinute)
	account := NewUserAccount(userId, password, expires)

	m.userAccounts[userId] = account
	log.Printf("user account created: %s", userId)
	return account, nil
}

func (m *UserAccountManager) GetUserAccount(userId string) (*UserAccount, bool) {
	account, exists := m.userAccounts[userId]
	return account, exists
}

func (m *UserAccountManager) Authenticate(userId string, password string) (*UserAccount, error) {
	account, exists := m.GetUserAccount(userId)
	if !exists {
		return nil, ErrLoginFailed
	}

	err := bcrypt.CompareHashAndPassword([]byte(account.HashedPassword), []byte(password))
	if err != nil {
		log.Printf("login failed for password %s: %s", userId, err.Error())
		return nil, ErrLoginFailed
	}
	return account, nil
}

func MakePassword() string {
	password := make([]byte, PasswordLength)
	for i := 0; i < PasswordLength; i++ {
		password[i] = PasswordChars[rand.Intn(len(PasswordChars))]
	}
	return string(password)
}
