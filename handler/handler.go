package handler

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/danielfmelo/load-funds-handler/domain"
	"github.com/danielfmelo/load-funds-handler/storage"
)

const (
	maximumValuePerDay        float64 = 5000
	maximumTransactionsPerDay int     = 3
	maximumValuePerWeek       float64 = 20000
)

type HandlerTransaction interface {
	Transaction(fund []byte)
}

type HandlerTransactionService struct {
	storage        storage.Database
	chPublisher    chan []byte
	chErrPublisher chan []byte
}

func New(
	storage storage.Database,
	chPublish chan []byte,
	chErrPublish chan []byte,
) *HandlerTransactionService {
	return &HandlerTransactionService{
		storage:        storage,
		chPublisher:    chPublish,
		chErrPublisher: chErrPublish,
	}
}

func (hs *HandlerTransactionService) Transaction(fund []byte) {
	var transaction domain.Transaction
	if err := json.Unmarshal(fund, &transaction); err != nil {
		msg := fmt.Sprintf("error to unmarshal fund %s", string(fund))
		hs.publishError(msg, err)
		return
	}
	if err := hs.storage.AddTransaction(transaction); err != nil {
		msg := fmt.Sprintf("error to add transaction with id: %s", transaction.ID)
		hs.publishError(msg, err)
		return
	}

	valid, daily, err := hs.isLoadPerDayValid(transaction)
	if err != nil {
		hs.publishError("error to validate transaction per day", err)
		return
	}
	if !valid {
		if err := hs.publishInvalidTransaction(transaction); err != nil {
			hs.publishError("error to publish invalid transaction", err)
		}
		return
	}

	valid, weekly, weeklyTotal, err := hs.isLoadPerWeekValid(transaction)
	if err != nil {
		hs.publishError("error to validate transaction per week", err)
		return
	}
	if !valid {
		if err := hs.publishInvalidTransaction(transaction); err != nil {
			hs.publishError("error to publish invalid transaction", err)
		}
		return
	}
	day := convertTimeToDay(transaction.Time)
	if err = hs.storage.AddDailyTransaction(transaction.CustomerID, day, daily); err != nil {
		hs.publishError("error to add daily transaction", err)
		return
	}
	if err = hs.storage.AddWeeklyTransaction(transaction.CustomerID, weekly, weeklyTotal); err != nil {
		hs.publishError("error to add weekly transaction", err)
		return
	}
	if err = hs.publishValidTransaction(transaction); err != nil {
		hs.publishError("error to publish valid transaction", err)
	}
}

func (hs *HandlerTransactionService) isLoadPerDayValid(transaction domain.Transaction) (bool, domain.DailyTransaction, error) {
	day := convertTimeToDay(transaction.Time)
	daily, err := hs.storage.GetDailyTransaction(transaction.CustomerID, day)
	if err != nil {
		if err != domain.ErrNotFound {
			return false, daily, err
		}
	}

	isMaximum, daily, err := isMaximumValueLoadPerDay(daily, transaction.LoadAmount)
	if err != nil {
		return false, daily, err
	}
	if isMaximum {
		return false, daily, nil
	}
	isMaximum, daily = isMaximumLoadPerDay(daily)
	return !isMaximum, daily, nil

}

func handleWeekNotFoundError(transaction domain.Transaction) (domain.WeeklyTransactionTotal, error) {
	transactionAmount, err := convertRawValue(transaction.LoadAmount)
	if err != nil {
		return domain.WeeklyTransactionTotal{}, err
	}
	week := domain.WeeklyTransactionTotal{Value: transactionAmount}
	return week, nil
}

func isMaximumValueLoadPerDay(daily domain.DailyTransaction, loadAmount string) (bool, domain.DailyTransaction, error) {
	transactionAmount, err := convertRawValue(loadAmount)
	if err != nil {
		return false, daily, err
	}
	tt := daily.DailyTotal + transactionAmount
	if tt > maximumValuePerDay {
		return true, daily, nil
	}
	daily.DailyTotal = daily.DailyTotal + transactionAmount
	return false, daily, err
}

func isMaximumLoadPerDay(daily domain.DailyTransaction) (bool, domain.DailyTransaction) {
	if daily.TransactionCount+1 > maximumTransactionsPerDay {
		return true, daily
	}
	daily.TransactionCount++
	return false, daily
}

func (hs *HandlerTransactionService) isLoadPerWeekValid(
	transaction domain.Transaction,
) (
	bool,
	domain.WeeklyTransaction,
	domain.WeeklyTransactionTotal,
	error,
) {
	year, week := transaction.Time.ISOWeek()
	weekly := domain.WeeklyTransaction{Year: year, Week: week}
	weeklyTotal, err := hs.storage.GetWeeklyTransaction(transaction.CustomerID, weekly)
	if err != nil {
		if err == domain.ErrNotFound {
			total, err := handleWeekNotFoundError(transaction)
			return true, weekly, total, err
		}
		return false, weekly, weeklyTotal, err
	}
	transactionAmount, err := convertRawValue(transaction.LoadAmount)
	if err != nil {
		return false, weekly, weeklyTotal, err
	}
	if (weeklyTotal.Value + transactionAmount) > maximumValuePerWeek {
		return false, weekly, weeklyTotal, nil
	}
	weeklyTotal.Value = weeklyTotal.Value + transactionAmount
	return true, weekly, weeklyTotal, nil
}

func convertTimeToDay(dateTime time.Time) string {
	return dateTime.Format(domain.DateLayout)
}

func convertRawValue(raw string) (float64, error) {
	value := strings.Replace(raw, "$", "", 1)
	valueFloat, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return float64(0.0), err
	}
	return valueFloat, nil
}

func (hs *HandlerTransactionService) publishValidTransaction(transaction domain.Transaction) error {
	valid := domain.TransactionResponse{
		ID:         transaction.ID,
		CustomerID: transaction.CustomerID,
		Accepted:   true,
	}
	event, err := json.Marshal(valid)
	if err != nil {
		return err
	}
	hs.chPublisher <- event
	return nil
}

func (hs *HandlerTransactionService) publishInvalidTransaction(transaction domain.Transaction) error {
	invalidTransaction := domain.TransactionResponse{
		ID:         transaction.ID,
		CustomerID: transaction.CustomerID,
		Accepted:   false,
	}
	event, err := json.Marshal(invalidTransaction)
	if err != nil {
		return err
	}
	hs.chPublisher <- event
	return nil
}

func (hs *HandlerTransactionService) publishError(message string, err error) {
	msg := fmt.Sprintf("msg: %s error: %s", message, err)
	hs.chErrPublisher <- []byte(msg)
}
