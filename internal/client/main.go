package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type BidResponse struct {
	Bid float64 `json:"bid"`
}

func main() {
	fmt.Println("Hello, client!")
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	var data BidResponse
	json.NewDecoder(res.Body).Decode(&data)
	print(data.Bid)
	saveQuotationToFile(fmt.Sprintf("%.2f", data.Bid))
}

func saveQuotationToFile(quote string) error {
	file, err := os.OpenFile("quotation.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Dolar: %s", quote))
	if err != nil {
		return err
	}

	return nil
}
