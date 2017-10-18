package main

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type MessageModel struct {
	gorm.Model
	Original     string
	HTML         string
	Created      time.Time
	Tags         []TagModel     `gorm:"ForeignKey:MessageRef"`
	RelatedUsers []UserModel    `gorm:"ForeignKey:MessageRef"`
	Attributes   []AttributeSet `gorm:"ForeignKey:MessageRef"`
	URLs         []URLModel     `gorm:"ForeignKey:MessageRef"`
}

type TagModel struct {
	gorm.Model
	MessageRef uint
	Slug       string
	ScreenName string
}

type UserModel struct {
	gorm.Model
	MessageRef uint
	Slug       string
	ScreenName string
}

type AttributeSet struct {
	gorm.Model
	MessageRef  uint
	Slug        string
	ScreenName  string
	DateValue   time.Time
	IntValue    int64
	FloatValue  float64
	StringValue string
}

type URLModel struct {
	gorm.Model
	MessageRef uint
	URL        string
	Scheme     string
	Domain     string
	Path       string
	Params     string
}

func init() {
	db, _ := gorm.Open("sqlite3", "db.sqlite3")
	defer db.Close()

	db.AutoMigrate(
		&MessageModel{},
		&URLModel{},
		&AttributeSet{},
		&UserModel{},
		&TagModel{},
	)

	log.Println("Automigrate")
}

func getDB() *gorm.DB {
	db, _ := gorm.Open("sqlite3", "db.sqlite3")
	db.LogMode(true)
	return db.Debug()
}

func LoadMessages(limit int, offset int, tags []string, users []string) []MessageModel {
	db := getDB()
	db.LogMode(true)
	defer db.Close()

	response := make([]MessageModel, 0)
	db = db.
		Preload("RelatedUsers").
		Preload("Tags").
		Preload("Attributes").
		Preload("URLs")

	if len(tags) > 0 {
		db = db.Joins(
			`JOIN "tag_models" 
				ON "tag_models"."message_ref" = "message_models"."id" 
				AND "tag_models."slug" in (?)`,
			tags,
		)
	}

	if len(users) > 0 {
		db = db.Joins(
			`JOIN "user_models" 
				ON "user_models"."message_ref" = "message_models"."id" 
				AND "user_models"."slug" in (?)`,
			users,
		)
	}

	db.
		Order(`"message_models"."created_at" DESC`).
		Limit(limit).
		Offset(offset).
		Find(&response)

	return response
}

func LoadMessage(messageID int) (MessageModel, bool) {
	db := getDB()

	defer db.Close()

	response := MessageModel{}
	notFound := db.
		Preload("RelatedUsers").
		Preload("Tags").
		Preload("Attributes").
		Preload("URLs").
		Find(&response, messageID).
		RecordNotFound()

	return response, !notFound

}

func Persist(message Message) *MessageModel {
	db := getDB()
	db.LogMode(true)
	defer db.Close()

	urls := make([]URLModel, 0)
	for _, url := range message.URLs {
		urls = append(urls, URLModel{URL: url})
	}

	tags := make([]TagModel, 0)
	for _, tag := range message.Tags {
		tags = append(tags, TagModel{Slug: tag, ScreenName: tag})
	}

	users := make([]UserModel, 0)
	for _, user := range message.RelatedUsers {
		users = append(users, UserModel{Slug: user, ScreenName: user})
	}

	attrs := make([]AttributeSet, 0)
	for key, val := range message.Attributes {
		attrs = append(attrs, AttributeSet{
			Slug:        key,
			ScreenName:  key,
			StringValue: val.(string),
		})
	}

	m := MessageModel{
		Original:     message.Original,
		HTML:         message.HTML,
		Created:      time.Now(),
		URLs:         urls,
		Attributes:   attrs,
		Tags:         tags,
		RelatedUsers: users,
	}

	db.Create(&m)

	return &m

}
