package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

// Book json struct
type Book struct {
	ID     string `json:"id"`
	Isbn   string `json:"isbn"`
	Title  string `json:"title"`
	Author `json:"author"`
}

// Author struct for Book struct
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var books []Book

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("getBooks request")

	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("getBook request")

	params := mux.Vars(r)

	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	// not found
	json.NewEncoder(w).Encode(Book{})
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("createBook request")

	var book Book

	_ = json.NewDecoder(r.Body).Decode(&book)

	book.ID = strconv.Itoa(rand.Intn(10000000))

	books = append(books, book)
	json.NewEncoder(w).Encode(book)

}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("updateBook request")

	params := mux.Vars(r)

	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)

			var book Book

			_ = json.NewDecoder(r.Body).Decode(&book)

			book.ID = params["id"]

			books = append(books, book)

			json.NewEncoder(w).Encode(book)

			return
		}
	}
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("deleteBook request")

	params := mux.Vars(r)

	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}

		json.NewEncoder(w).Encode(books)
	}
}

func main() {
	r := mux.NewRouter()

	// mock data
	books = append(books, Book{
		ID:    "1",
		Isbn:  "789234",
		Title: "pejuang",
		Author: Author{
			Firstname: "Logan",
			Lastname:  "Paul",
		},
	})

	books = append(books, Book{
		ID:    "2",
		Isbn:  "234234",
		Title: "pemburu",
		Author: Author{
			Firstname: "kera",
			Lastname:  "sakti",
		},
	})

	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	port := "8000"

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	fmt.Println("Listening on port", port)

	log.Fatal(http.ListenAndServe(":"+port, r))

}
