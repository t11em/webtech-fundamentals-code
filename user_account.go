package main

import (
	"time"
	_ "time/tzdata"

	"golang.org/x/crypto/bcrypt"
)

const UserAccountLimitInMinute = 60

const PasswordLength = 10

const PasswordChars = "23456789abcdefghijkmnpqrstuvwxyz"

type UserAccount struct {
	Id             string
	HashedPassword string
	Expires        time.Time
	ToDoList       []string
}

func NewUserAccount(userId string, plainPassword string, expires time.Time) *UserAccount {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	account := &UserAccount{
		Id:             userId,
		HashedPassword: string(hashedPassword),
		Expires:        expires,
		ToDoList:       make([]string, 1, 10),
	}
	return account
}

func (u UserAccount) ExpiresText() string {
	return u.Expires.Format("2006/01/02 15:04:05")
}
