package storage

import "github.com/danielfmelo/load-funds-handler/domain"

type Database interface {
	AddTransaction(transaction domain.Transaction) error
	AddDailyTransaction(customerID, day string, daily domain.DailyTransaction) error
	AddWeeklyTransaction(customerID string, week domain.WeeklyTransaction, total domain.WeeklyTransactionTotal) error
	GetDailyTransaction(customerID, day string) (domain.DailyTransaction, error)
	GetWeeklyTransaction(customerID string, week domain.WeeklyTransaction) (domain.WeeklyTransactionTotal, error)
}
