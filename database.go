// archiver - database.go
// 2024-02-20 20:00
// Benny <benny.think@gmail.com>

package main

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var db, _ = gorm.Open(sqlite.Open("archiver.db"), &gorm.Config{})

type User struct {
	gorm.Model
	UserID int64  `gorm:"index;unique"` // 假设Telegram的用户ID是整数
	Mode   string // AI模式和普通模式
	Link   string
	Token  string // Future use
}

type Chat struct {
	gorm.Model
	UserID int64  `gorm:"index"` // userID
	Role   string // 发送者角色：model 或 user
	Text   string // 保存的消息文本
}

func getUser(userID int64) *User {
	var user User
	db.Where("user_id = ?", userID).First(&user)
	return &user
}

func enableAI(userID int64, link string) *User {
	user := getUser(userID)
	user.Mode = modeAI
	user.Link = link
	if user.UserID == 0 {
		// create new user
		user.UserID = userID
		db.Create(&user)
	} else {
		// update db
		db.Save(&user)
	}

	return user
}

func disableAI(userId int64) {
	// reset mode to normal, clear link
	user := getUser(userId)
	user.Mode = modeNormal
	user.Link = ""
	db.Save(&user)
}

func getChats(userID int64) []Chat {
	var chats []Chat
	db.Where("user_id = ?", userID).Find(&chats)
	return chats
}

func getChatsCount(userID int64) int64 {
	var count int64
	db.Model(&Chat{}).Where("user_id = ?", userID).Count(&count)
	return count
}

func addChat(userId int64, role, text string) {
	chat := Chat{
		UserID: userId,
		Role:   role,
		Text:   text,
	}
	db.Create(&chat)
}

func deleteChat(userID int64) int64 {
	// delete according to userID
	t := db.Where("user_id = ?", userID).Delete(&Chat{})
	return t.RowsAffected
}
