package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

const DBName = "ems.db"
const CfAccountNumber = "5930846a5870031c415bb26e42e38833"

//type Labels struct {
//	gorm.Model
//	Name       string  `json:"cv_name"`
//	Confidence float64 `json:"cv_confidence"`
//}
//
//type Photograph struct {
//	gorm.Model
//	Name           string    `json:"name"`
//	ParsedName     string    `json:"parsed_name"`
//	Artist         string    `json:"artist"`
//	CaptureDate    time.Time `json:"capture_date"`
//	Description    string    `json:"description"`
//	Caption        string    `json:"caption"`
//	Sha256         string    `json:"sha256"`
//	PerceptualHash string    `json:"perceptual_hash"`
//
//	Classification struct {
//		Labels []Labels
//	}
//}

type CfImage struct {
	Result     Result        `json:"result,omitempty"`
	ResultInfo interface{}   `json:"result_info,omitempty"`
	Success    bool          `json:"success,omitempty"`
	Errors     []interface{} `json:"errors,omitempty"`
	Messages   []interface{} `json:"messages,omitempty"`
}

type Photograph struct {
	ID               string    `json:"Id,omitempty"`
	Caption          string    `json:"Caption,omitempty"`
	City             string    `json:"City,omitempty"`
	Country          string    `json:"Country,omitempty"`
	State            string    `json:"State,omitempty"`
	DateTimeOriginal time.Time `json:"DateTimeOriginal,omitempty"`
	Headline         string    `json:"Headline,omitempty"`
	Keywords         string    `json:"Keywords,omitempty"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type Images struct {
	ID                string     `json:"id,omitempty"`
	Filename          string     `json:"filename,omitempty"`
	Photograph        Photograph `json:"meta,omitempty"`
	Uploaded          time.Time  `json:"uploaded,omitempty"`
	RequireSignedURLs bool       `json:"requireSignedURLs,omitempty"`
	Variants          []string   `json:"variants,omitempty"`
}

type Result struct {
	Images []Images `json:"images,omitempty"`
}

func InitialMigration() {
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
	err := DB.Create(&photograph)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(photograph)
}

func UpdatePhotograph(w http.ResponseWriter, r *http.Request) {
}

// UpdatePhotographs by reading all the images in Cloudflare and storing them in
// the database if they don't already exist
func UpdatePhotographs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	perPage := 100
	keepReading := true

	for page := 1; keepReading; page++ {

		req, err := http.NewRequest(http.MethodGet, "https://api.cloudflare.com/client/v4/accounts/"+CfAccountNumber+"/images/v1?page="+strconv.Itoa(page)+"&per_page="+strconv.Itoa(perPage), nil)
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+os.Getenv("CF_IMAGES_STREAM_READ_ONLY"))
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		body, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}
		resp.Body.Close()

		cfi := CfImage{}

		jsonErr := json.Unmarshal(body, &cfi)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}

		if cfi.Success == false {
			log.Fatal("Cloudflare returned an error")
		}

		log.Printf("Page: %v", page)
		log.Printf("Number returned: %v", len(cfi.Result.Images))

		for i, img := range cfi.Result.Images {
			err := DB.Create(&img.Photograph)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Creating: Page: %v, i: %v: id: %v, PhotographID: %v\n", page, i, img.ID, img.Photograph.ID)
		}

		keepReading = len(cfi.Result.Images) == perPage
	}

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
