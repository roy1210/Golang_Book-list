package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"database/sql"
	_ "database/sql"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/subosito/gotenv"
)

// Book struct
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   string `json:"year"`
}

var books []Book
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

	pgURL, err := pq.ParseURL(os.Getenv("ELEPHANTSQL_URL"))
	logFatal(err)

	db, err = sql.Open("postgres", pgURL)
	logFatal(err)

	err = db.Ping()
	logFatal(err)

	// ElephantSQL: -`dbname`, `host`, `password`, `port`, `user`
	log.Println(pgURL)

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
	var book Book
	books = []Book{}

	rows, err := db.Query("select * from books")
	logFatal(err)

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		logFatal(err)

		books = append(books, book)
	}
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	params := mux.Vars(r)

	rows := db.QueryRow("select * from books where id=$1", params["id"])

	err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	logFatal(err)

	json.NewEncoder(w).Encode(book)
}

// r: request => such come from Client or Postman
// Book.ID made in this struct is not syncing with the serialized ID in the elephantSQL. Will replace by the ID made by elephantSLQ as serialized.
func addBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	var bookID int
	// Decode: Map to the arg
	json.NewDecoder(r.Body).Decode(&book)

	// log.Println(book) >> 2019/05/12 05:00:48 {19 C++ is old Mr. C++ 2014}
	// log.Println(book)
	// log.Println(reflect.TypeOf(book))
	// >> main.Book => Means struct

	// Scan: likes copy the row to the bookID row
	err := db.QueryRow("insert into books (title, author, year) values ($1, $2, $3) RETURNING id;", book.Title, book.Author, book.Year).Scan(&bookID)
	logFatal(err)

	// ex >> 4    NOT LIKES `ID:99` ALLOCATED BY MY SELF
	json.NewEncoder(w).Encode(bookID)
	log.Println(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	json.NewDecoder(r.Body).Decode(&book)

	// where : which data is going to update
	result, err := db.Exec("update books set title = $1, author=$2, year=$3 where id=$4 RETURNING id", &book.Title, &book.Author, &book.Year, &book.ID)

	rowsUpdated, err := result.RowsAffected()
	logFatal(err)
	json.NewEncoder(w).Encode(rowsUpdated)
}

func removeBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	result, err := db.Exec("delete from books where id = $1", params["id"])
	logFatal(err)

	rowsDeleted, err := result.RowsAffected()
	logFatal(err)

	json.NewEncoder(w).Encode(rowsDeleted)
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
