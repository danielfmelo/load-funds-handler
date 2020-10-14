package memory

import (
	"fmt"

	"github.com/danielfmelo/load-funds-handler/domain"
)

type Database struct {
	transactions map[string]map[string]domain.Transaction
	daily        map[string]map[string]domain.DailyTransaction
	weekly       map[string]map[domain.WeeklyTransaction]domain.WeeklyTransactionTotal
}

func New() *Database {
	return &Database{
		transactions: make(map[string]map[string]domain.Transaction),
		daily:        make(map[string]map[string]domain.DailyTransaction),
		weekly:       make(map[string]map[domain.WeeklyTransaction]domain.WeeklyTransactionTotal),
	}
}

func (d *Database) AddTransaction(transaction domain.Transaction) error {
	if transaction.ID == "" {
		return domain.ErrTransactionEmptyID
	}
	t, ok := d.transactions[transaction.ID]
	if !ok {
		d.transactions[transaction.ID] = map[string]domain.Transaction{transaction.CustomerID: transaction}
		return nil
	}
	if _, ok := t[transaction.CustomerID]; !ok {
		t[transaction.CustomerID] = transaction
		if transaction.ID == "6928" {
			fmt.Println("dentro do costumer")
		}
		return nil
	}
	return domain.ErrTransactionAlreadyExist
}

func (d *Database) transactionExist(id string) bool {
	_, ok := d.transactions[id]
	return ok
}

func (d *Database) AddDailyTransaction(customerID, day string, daily domain.DailyTransaction) error {
	//fmt.Println(daily.DailyTotal)
	d.daily[customerID] = map[string]domain.DailyTransaction{day: daily}
	return nil
}

func (d *Database) AddWeeklyTransaction(customerID string, week domain.WeeklyTransaction, total domain.WeeklyTransactionTotal) error {
	d.weekly[customerID] = map[domain.WeeklyTransaction]domain.WeeklyTransactionTotal{week: total}
	return nil
}

func (d *Database) GetDailyTransaction(customerID, day string) (domain.DailyTransaction, error) {
	customer, ok := d.daily[customerID]
	if !ok {
		return domain.DailyTransaction{}, domain.ErrNotFound
	}
	dailyTransaction, ok := customer[day]
	if !ok {
		return domain.DailyTransaction{}, domain.ErrNotFound
	}
	return dailyTransaction, nil
}

func (d *Database) GetWeeklyTransaction(customerID string, week domain.WeeklyTransaction) (domain.WeeklyTransactionTotal, error) {
	customer, ok := d.weekly[customerID]
	if !ok {
		return domain.WeeklyTransactionTotal{}, domain.ErrNotFound
	}
	weeklyTransaction, ok := customer[week]
	if !ok {
		return domain.WeeklyTransactionTotal{}, domain.ErrNotFound
	}
	return weeklyTransaction, nil
}
