package main

import (
	"time"
)

type Post struct {
	Userkey string `json:"user_key" gorm:"primary_key"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

type Playing struct {
	Userkey      string    `sql:"not null" gorm:"primary_key"`
	CreatedAt    time.Time `sql:"not null"`
	TryCount     int       `sql:"not null"`
	AnswerNumber string    `sql:"not null"`
}

type Record struct {
	Userkey     string    `sql:"not null" gorm:"primary_key"`
	CreatedAt   time.Time `sql:"not null"`
	TryCount    int       `sql:"not null"`
	SpendedTime int       `sql:"not null"`
	Nickname    string    `sql:"not null"`
	Score       int       `sql:"not null"`
	rank        int
}

type RecordForShow struct {
	Userkey  string
	Rank     int
	NickName string
	Score    int
}

type UserInfo struct {
	Userkey     string    `sql:"not null"`
	CreatedAt   time.Time `sql:"not null"`
	Nickname    string
	IsItUpdated bool `sql:"not null"`
}
