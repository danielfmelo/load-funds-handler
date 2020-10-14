package domain

import "time"

const DateLayout = "2006-02-01"

type Transaction struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customer_id"`
	LoadAmount string    `json:"load_amount"`
	Time       time.Time `json:"time"`
}

type DailyTransaction struct {
	Transaction      Transaction
	TransactionCount int
	DailyTotal       float64
}

type WeeklyTransaction struct {
	Year int
	Week int
}

type WeeklyTransactionTotal struct {
	Value float64
}

type TransactionResponse struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	Accepted   bool   `json:"accepted"`
}
