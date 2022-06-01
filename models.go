package main

import (
	"encoding/json"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

const DBName = "ems.db"

// type Labels struct {
// 	gorm.Model
// 	Name       string  `json:"cv_name"`
// 	Confidence float64 `json:"cv_confidence"`
// }

type Photograph struct {
	gorm.Model
	Name           string    `json:"name"`
	ParsedName     string    `json:"parsed_name"`
	Artist         string    `json:"artist"`
	CaptureDate    time.Time `json:"capture_date"`
	Description    string    `json:"description"`
	Caption        string    `json:"caption"`
	Sha256         string    `json:"sha256"`
	PerceptualHash string    `json:"perceptual_hash"`

	// Classification struct {
	// 	Labels []Labels
	// }
}

func InitalMigration() {
	DB, err = gorm.Open(sqlite.Open(DBName), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	DB.AutoMigrate(&Photograph{})
	DB.AutoMigrate(&Bookmark{})
	DB.AutoMigrate(&Chatter{})
}

func GetPhotographs(w http.ResponseWriter, r *http.Request) {
}

func GetPhotograph(w http.ResponseWriter, r *http.Request) {
}

func CreatePhotograph(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var photograph Photograph

	json.NewDecoder(r.Body).Decode(&photograph)
	DB.Create(&photograph)
	json.NewEncoder(w).Encode(photograph)
}

func UpdatePhotograph(w http.ResponseWriter, r *http.Request) {
}

func DeletePhotographs(w http.ResponseWriter, r *http.Request) {
}

// Bookmarks

type Bookmark struct {
	gorm.Model
	Href        string    `json:"href"`
	Description string    `json:"description"`
	Extended    string    `json:"extended"`
	Hash        string    `json:"hash"`
	Date        time.Time `json:"date"`
}

func GetBookmarks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var bookmarks []Bookmark
	DB.Find(&bookmarks)
	json.NewEncoder(w).Encode(bookmarks)
}

func GetBookmark(w http.ResponseWriter, r *http.Request) {
}

func CreateBookmark(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var bookmark Bookmark

	json.NewDecoder(r.Body).Decode(&bookmark)

	if bookmark.Href != "" {
		DB.Create(&bookmark)
	}
	
	json.NewEncoder(w).Encode(bookmark)
}

func UpdateBookmark(w http.ResponseWriter, r *http.Request) {
}

func DeleteBookmarks(w http.ResponseWriter, r *http.Request) {
}

// Chatter

type Chatter struct {
	gorm.Model
	Title    string    `json:"title"`
	Body     string    `json:"body"`
	Image    string    `json:"image"`
	CFImage  string    `json:"cf_image"`
	Date     time.Time `json:"date"`
	Location string    `json:"location"`
	Weather  string    `json:"weather"`
}

func GetChatters(w http.ResponseWriter, r *http.Request) {
}

func GetChatter(w http.ResponseWriter, r *http.Request) {
}

func CreateChatter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var Chatter Chatter

	json.NewDecoder(r.Body).Decode(&Chatter)
	DB.Create(&Chatter)
	json.NewEncoder(w).Encode(Chatter)
}

func UpdateChatter(w http.ResponseWriter, r *http.Request) {
}

func DeleteChatters(w http.ResponseWriter, r *http.Request) {
}
