package database

import (
	"context"
	"database/sql"
	"time"
)

type SaveQuotation struct {
	Name string
	Bid  float64
	Ask  float64
}
type QuotationDB struct {
	DB *sql.DB
}

func NewQuotationDB(db *sql.DB) *QuotationDB {
	return &QuotationDB{
		DB: db,
	}
}

func (qdb *QuotationDB) Save(quotation *SaveQuotation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	stmt, err := qdb.DB.PrepareContext(ctx, "INSERT INTO quotes (name, bid, ask) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, quotation.Name, quotation.Bid, quotation.Ask)
	if err != nil {
		return err
	}
	return nil
}
