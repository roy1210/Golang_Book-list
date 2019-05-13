package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/roy1210/Study/book-list/models"
	bookRepository "github.com/roy1210/Study/book-list/repository/book"
)

type Controller struct{}

var books []models.Book

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// GetBooks is method which wrapped getBooks function
// return http.HandlerFunc type
// Wrapする理由は、HandleFuncにてw と r以外にdbを引数として取りたいから
func (c Controller) GetBooks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book models.Book
		books = []models.Book{}

		bookRepo := bookRepository.BookRepository{}
		books = bookRepo.GetBooks(db, book, books)

		json.NewEncoder(w).Encode(books)
	}
}

func (c Controller) GetBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book models.Book
		params := mux.Vars(r)

		books = []models.Book{}
		bookRepo := bookRepository.BookRepository{}

		id, err := strconv.Atoi(params["id"])
		logFatal(err)

		book = bookRepo.GetBook(db, book, id)

		json.NewEncoder(w).Encode(book)
	}
}

func (c Controller) AddBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book models.Book
		var bookID int
		// Decode: Map to the arg
		json.NewDecoder(r.Body).Decode(&book)

		//Scan: likes copy the row to the bookID row
		err := db.QueryRow("insert into books (title, author, year) values ($1, $2, $3) RETURNING id;", book.Title, book.Author, book.Year).Scan(&bookID)
		logFatal(err)

		//ex >> 4    NOT LIKES `ID:99` ALLOCATED BY MY SELF
		json.NewEncoder(w).Encode(bookID)
	}
}

func (c Controller) UpdateBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book models.Book
		json.NewDecoder(r.Body).Decode(&book)

		bookRepo := bookRepository.BookRepository{}
		rowsUpdated := bookRepo.UpdateBook(db, book)

		json.NewEncoder(w).Encode(rowsUpdated)
	}
}

func (c Controller) RemoveBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		bookRepo := bookRepository.BookRepository{}

		id, err := strconv.Atoi(params["id"])
		logFatal(err)

		rowsDeleted := bookRepo.RemoveBook(db, id)

		json.NewEncoder(w).Encode(rowsDeleted)
	}
}

// ラップの方法
// Wrap by GetBooks method. returnに貼り付けで、`getBooks`を`func`に書き換える

// func getBooks(w http.ResponseWriter, r *http.Request) {
// 	var book models.Book
// 	books = []models.Book{}

// 	rows, err := db.Query("select * from books")
// 	logFatal(err)

// 	defer rows.Close()

// 	for rows.Next() {
// 		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
// 		logFatal(err)

// 		books = append(books, book)
// 	}
// 	json.NewEncoder(w).Encode(books)
// }

// func getBook(w http.ResponseWriter, r *http.Request) {
// 	var book models.Book
// 	params := mux.Vars(r)

// 	rows := db.QueryRow("select * from books where id=$1", params["id"])

// 	err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
// 	logFatal(err)

// 	json.NewEncoder(w).Encode(book)
// }

// r: request => such come from Client or Postman
// Book.ID made in this struct is not syncing with the serialized ID in the elephantSQL. Will replace by the ID made by elephantSLQ as serialized.
// func addBook(w http.ResponseWriter, r *http.Request) {
// 	var book models.Book
// 	var bookID int
// 	// Decode: Map to the arg
// 	json.NewDecoder(r.Body).Decode(&book)

// log.Println(book) >> 2019/05/12 05:00:48 {19 C++ is old Mr. C++ 2014}
// log.Println(book)
// log.Println(reflect.TypeOf(book))
// >> main.Book => Means struct

// Scan: likes copy the row to the bookID row
// err := db.QueryRow("insert into books (title, author, year) values ($1, $2, $3) RETURNING id;", book.Title, book.Author, book.Year).Scan(&bookID)
// logFatal(err)

// ex >> 4    NOT LIKES `ID:99` ALLOCATED BY MY SELF
// 	json.NewEncoder(w).Encode(bookID)
// 	log.Println(book)
// }

// func updateBook(w http.ResponseWriter, r *http.Request) {
// 	var book models.Book
// 	json.NewDecoder(r.Body).Decode(&book)

// where : which data is going to update
// 	result, err := db.Exec("update books set title = $1, author=$2, year=$3 where id=$4 RETURNING id", &book.Title, &book.Author, &book.Year, &book.ID)

// 	rowsUpdated, err := result.RowsAffected()
// 	logFatal(err)
// 	json.NewEncoder(w).Encode(rowsUpdated)
// }

// func removeBook(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	result, err := db.Exec("delete from books where id = $1", params["id"])
// 	logFatal(err)

// 	rowsDeleted, err := result.RowsAffected()
// 	logFatal(err)

// 	json.NewEncoder(w).Encode(rowsDeleted)
// }
