package gomsglog

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jblawatt/gomsglog/gomsglog/parsers"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"

	jww "github.com/spf13/jwalterweatherman"

	// Datenbank Dialekt einbinden.
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

func (a *AttributeSet) String() string {
	return a.StringValue
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

func AutoMigrate() {
	db := GetDB()
	defer db.Close()

	db.AutoMigrate(
		&MessageModel{},
		&URLModel{},
		&AttributeSet{},
		&UserModel{},
		&TagModel{},
	)

}

func GetDB() *gorm.DB {
	dbpath := viper.GetString("database.connectionstring")
	dbabspath, pathError := filepath.Abs(dbpath)
	if pathError != nil {
		fmt.Fprintf(os.Stderr, "Error loading absolute database path %s.", dbpath)
		os.Exit(1)
	}
	jww.DEBUG.Printf("Loading Database from %s", dbabspath)
	db, _ := gorm.Open(
		viper.GetString("database.dialect"),
		dbabspath,
	)
	debug := viper.GetBool("database.debug")
	if debug {
		db.LogMode(true)
		return db.Debug()
	}
	return db

}

func LoadMessages(limit int, offset int, tags []string, users []string) []MessageModel {
	db := GetDB()
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
				AND "tag_models"."slug" in (?)`,
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
	db := GetDB()

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

func Persist(message parsers.Message) *MessageModel {
	db := GetDB()
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

func DeleteMessage(messageID int) bool {
	db := GetDB()
	notFound := db.Delete(new(MessageModel), messageID).RecordNotFound()
	return !notFound
}

type UserSet struct {
	Slug       string `gorm:"slug"`
	ScreenName string `gorm:"screenname"`
}

func LoadUsers() []UserSet {
	db := GetDB()
	defer db.Close()

	rows, err := db.Table("user_models").Select(`DISTINCT "slug", "screen_name"`).Order(`"user_models"."slug"`).Rows()
	if err != nil {
		panic(err)
	}

	response := make([]UserSet, 0)

	for rows.Next() {
		var u UserSet
		if err := rows.Scan(&u.Slug, &u.ScreenName); err != nil {
			panic(err)
		}

		response = append(response, u)
	}

	return response
}

type TagSet struct {
	Slug       string `gorm:"slug"`
	ScreenName string `gorm:"screenname"`
}

func LoadTags() []TagSet {
	db := GetDB()
	defer db.Close()

	rows, err := db.Table("tag_models").Select(`DISTINCT "slug", "screen_name"`).Order(`"tag_models"."slug"`).Rows()
	if err != nil {
		panic(err)
	}

	response := make([]TagSet, 0)

	for rows.Next() {
		var t TagSet
		if err := rows.Scan(&t.Slug, &t.ScreenName); err != nil {
			panic(err)
		}

		response = append(response, t)
	}

	return response
}
