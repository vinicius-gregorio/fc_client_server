package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	database "github.com/vinicius-gregorio/fc_client_server/internal/server/db"
)

// type Server struct {
// 	DB db.DB
// }

// func (s *Server) NewServer(database db.DB) *Server {
// 	return &Server{DB: database}
// }

type Quotation struct {
	ID         int     `json:"id"`
	Code       string  `json:"code"`
	Codein     string  `json:"codein"`
	Name       string  `json:"name"`
	High       float64 `json:"high,string"`
	Low        float64 `json:"low,string"`
	VarBid     float64 `json:"varBid,string"`
	PctChange  float64 `json:"pctChange,string"`
	Bid        float64 `json:"bid,string"`
	Ask        float64 `json:"ask,string"`
	Timestamp  string  `json:"timestamp"`
	CreateDate string  `json:"create_date"`
}

func main() {
	// This is the main function
	fmt.Println("Hello, server!")

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", "root", "root", "localhost", "3306", "challenge01"))

	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	defer db.Close()

	http.HandleFunc("/cotacao", bidHandler(db))
	http.ListenAndServe(":8080", nil)

}

type BidResponse struct {
	Bid float64 `json:"bid"`
}

func bidHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		log.Println("Handler started")
		defer log.Println("Handler ended")

		q, err := getQuotation(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			panic(err)
		}

		err = saveQuotation(db, q)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			panic(err)
		}

		select {
		case <-ctx.Done():
			fmt.Fprintf(w, "done")
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			bidResponse := BidResponse{Bid: q.Bid}
			json.NewEncoder(w).Encode(bidResponse)
		}
	}

}

func getQuotation(ctx context.Context) (*Quotation, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get valid response: %s", res.Status)
	}

	var result map[string]Quotation
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	quotation, exists := result["USDBRL"]
	if !exists {
		return nil, fmt.Errorf("quotation not found in response")
	}

	return &quotation, nil
}

func saveQuotation(db *sql.DB, quotation *Quotation) error {
	qdb := database.NewQuotationDB(db)
	err := qdb.Save(&database.SaveQuotation{
		Name: quotation.Name,
		Bid:  quotation.Bid,
		Ask:  quotation.Ask,
	})
	if err != nil {
		return err
	}
	return nil
}
