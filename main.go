package main

import (
	"log"
	"net/http"

	"database/sql"
	_ "database/sql"

	"github.com/gorilla/mux"
	"github.com/roy1210/Study/book-list/controllers"
	"github.com/roy1210/Study/book-list/driver"
	"github.com/roy1210/Study/book-list/models"
	"github.com/subosito/gotenv"
)

var books []models.Book
var db *sql.DB

func init() {
	gotenv.Load()
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	db = driver.ConnectDB()
	controller := controllers.Controller{}

	// PUT: replace resource information
	router := mux.NewRouter()
	router.HandleFunc("/books", controller.GetBooks(db)).Methods("GET")
	router.HandleFunc("/books/{id}", controller.GetBook(db)).Methods("GET")
	router.HandleFunc("/books", controller.AddBook(db)).Methods("POST")
	router.HandleFunc("/books", controller.UpdateBook(db)).Methods("PUT")
	router.HandleFunc("/books/{id}", controller.RemoveBook(db)).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}

/*
BEFORE USE DB

func main() {

	books = append(books,
		Book{ID: 1, Title: "Golang Pointers", Author: "Mr. Golang", Year: "2010"},
		Book{ID: 2, Title: "Goroutine", Author: "Mr. Goroutine", Year: "2011"},
		Book{ID: 3, Title: "Golang routers", Author: "Mr. Router", Year: "2012"},
		Book{ID: 4, Title: "Golang concurrency", Author: "Mr. Currency", Year: "2013"},
		Book{ID: 5, Title: "Golang good parts", Author: "Mr. Good", Year: "2014"},
	)

	// PUT: replace resource information
	router := mux.NewRouter()
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", addBook).Methods("POST")
	router.HandleFunc("/books", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", removeBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	// log.Println("Gets all books")
	// json.encode : Marshal to json format
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	// Vars(r): return  r of STRING Map >>  map[id:3]
	// params["id"] : "id" is the key name anyname can take it. Value is in the request of URL.
	params := mux.Vars(r)
	i, _ := strconv.Atoi(params["id"])

	// Check the type of var => log.Println(reflect.TypeOf(i))

	for _, book := range books {
		if book.ID == i {
			json.NewEncoder(w).Encode(&book)
		}
	}
}

// r: request => such come from Client or Postman
func addBook(w http.ResponseWriter, r *http.Request) {
	var book Book

	// _ = json.NewDecoder(r.Body).Decode(&book)
	json.NewDecoder(r.Body).Decode(&book)

	books = append(books, book)
	json.NewEncoder(w).Encode(books)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	json.NewDecoder(r.Body).Decode(&book)

	for i, item := range books {
		if item.ID == book.ID {
			books[i] = book
		}
	}

	json.NewEncoder(w).Encode(books)
}

func removeBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	for i, item := range books {
		if item.ID == id {
			// books[i+1:]...: 追加される要素については、スライスではなく一つ一つという意味での...だと思う
			books = append(books[:i], books[i+1:]...)

				ID 1 2 3 4 5
				i	 0 1 2 3 4
					     D

		}
	}
	json.NewEncoder(w).Encode(books)
}
*/
