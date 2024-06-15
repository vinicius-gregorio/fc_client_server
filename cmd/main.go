package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// This is the main function
	fmt.Println("Hello, World!")

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", "root", "root", "mysql", "3306", "wallet"))

	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	defer db.Close()

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)

}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("Handler started")
	defer log.Println("Handler ended")
	select {
	case <-ctx.Done():
		err := ctx.Err()
		log.Println("Handler err:", err)
	case <-time.After(5 * time.Second):
		fmt.Fprintf(w, "Hello, World!")

	}
}
