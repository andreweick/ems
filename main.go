package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func initializeRouter() {
	r := mux.NewRouter()

	r.HandleFunc("/photographs", GetPhotographs).Methods("GET")
	r.HandleFunc("/photographs/{id}", GetPhotograph).Methods("GET")
	r.HandleFunc("/photographs", CreatePhotograph).Methods("POST")
	r.HandleFunc("/photographs/{id}", UpdatePhotograph).Methods("PUT")
	r.HandleFunc("/photographs/{id}", DeletePhotographs).Methods("DELETE")

	// Bookmarks
	r.HandleFunc("/bookmarks", GetBookmarks).Methods("GET")
	r.HandleFunc("/bookmarks/{id}", GetBookmark).Methods("GET")
	r.HandleFunc("/bookmarks", CreateBookmark).Methods("POST")
	r.HandleFunc("/bookmarks/{id}", UpdateBookmark).Methods("PUT")
	r.HandleFunc("/bookmarks/{id}", DeleteBookmarks).Methods("DELETE")

	// Bookmarks
	r.HandleFunc("/chatters", GetChatters).Methods("GET")
	r.HandleFunc("/chatters/{id}", GetChatter).Methods("GET")
	r.HandleFunc("/chatters", CreateChatter).Methods("POST")
	r.HandleFunc("/chatters/{id}", UpdateChatter).Methods("PUT")
	r.HandleFunc("/chatters/{id}", DeleteChatters).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}

func main() {
	InitalMigration()
	initializeRouter()
}
