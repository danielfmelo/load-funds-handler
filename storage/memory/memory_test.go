package memory_test

import (
	"testing"
	"time"

	"github.com/danielfmelo/load-funds-handler/domain"
	"github.com/danielfmelo/load-funds-handler/storage/memory"
	"github.com/stretchr/testify/assert"
)

func TestAddTransaction(t *testing.T) {
	testCases := []struct {
		name        string
		fund        domain.Transaction
		errExpected error
	}{
		{
			name:        "add should work",
			fund:        domain.Transaction{ID: "123", CustomerID: "1234", LoadAmount: "$1", Time: time.Now()},
			errExpected: nil,
		},
		{
			name:        "should return empty ID",
			fund:        domain.Transaction{ID: "", CustomerID: "1234", LoadAmount: "$1", Time: time.Now()},
			errExpected: domain.ErrTransactionEmptyID,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := memory.New()
			err := m.AddTransaction(tc.fund)
			assert.Equal(t, tc.errExpected, err)
		})
	}
}

func TestAddTransactionShouldWorkWithDifferentCustomers(t *testing.T) {
	fund := domain.Transaction{
		ID:         "123",
		CustomerID: "1234",
		LoadAmount: "$1",
		Time:       time.Now(),
	}
	m := memory.New()
	err := m.AddTransaction(fund)
	assert.Nil(t, err)
	fund.CustomerID = "567"
	err = m.AddTransaction(fund)
	assert.Nil(t, err)
}

func TestAddTransactionShouldReturnAlreadyExist(t *testing.T) {
	fund := domain.Transaction{
		ID:         "123",
		CustomerID: "1234",
		LoadAmount: "$1",
		Time:       time.Now(),
	}
	m := memory.New()
	err := m.AddTransaction(fund)
	assert.Nil(t, err)
	err = m.AddTransaction(fund)
	assert.Equal(t, domain.ErrTransactionAlreadyExist, err)
}

func TestAddDailyTransaction(t *testing.T) {
	fund := domain.Transaction{
		ID:         "123",
		CustomerID: "1234",
		LoadAmount: "$1",
		Time:       time.Now(),
	}
	m := memory.New()
	dailyTransaction := domain.DailyTransaction{TransactionCount: 1, Transaction: fund, DailyTotal: 10}
	err := m.AddDailyTransaction(fund.CustomerID, fund.Time.Format(domain.DateLayout), dailyTransaction)
	assert.Nil(t, err)
}

func TestAddWeeklyTransaction(t *testing.T) {
	fund := domain.Transaction{
		ID:         "123",
		CustomerID: "1234",
		LoadAmount: "$1",
		Time:       time.Now(),
	}
	m := memory.New()
	year, week := fund.Time.ISOWeek()
	weeklyTransaction := domain.WeeklyTransaction{Year: year, Week: week}
	weeklyTransactionTotal := domain.WeeklyTransactionTotal{Value: 1}
	err := m.AddWeeklyTransaction(fund.CustomerID, weeklyTransaction, weeklyTransactionTotal)
	assert.Nil(t, err)
}

func TestGetDailyTransaction(t *testing.T) {
	fund := domain.Transaction{
		ID:         "123",
		CustomerID: "1234",
		LoadAmount: "$10",
		Time:       time.Now(),
	}
	m := memory.New()
	dailyTransaction := domain.DailyTransaction{TransactionCount: 1, Transaction: fund, DailyTotal: 10}
	day := fund.Time.Format(domain.DateLayout)
	err := m.AddDailyTransaction(fund.CustomerID, day, dailyTransaction)
	assert.Nil(t, err)
	daily, err := m.GetDailyTransaction(fund.CustomerID, day)
	assert.Nil(t, err)
	assert.Equal(t, float64(10), daily.DailyTotal)
	assert.Equal(t, 1, daily.TransactionCount)
	assert.Equal(t, fund, daily.Transaction)
}

func TestGetDailyTransactionShouldReturnNotFound(t *testing.T) {
	fund := domain.Transaction{
		ID:         "123",
		CustomerID: "1234",
		LoadAmount: "$10",
		Time:       time.Now(),
	}
	m := memory.New()
	dailyTransaction := domain.DailyTransaction{TransactionCount: 1, Transaction: fund, DailyTotal: 10}
	day := fund.Time.Format(domain.DateLayout)
	err := m.AddDailyTransaction(fund.CustomerID, day, dailyTransaction)
	assert.Nil(t, err)
	_, err = m.GetDailyTransaction("888", day)
	assert.Equal(t, domain.ErrNotFound, err)
}

func TestGetWeeklyTransaction(t *testing.T) {
	fund := domain.Transaction{
		ID:         "123",
		CustomerID: "1234",
		LoadAmount: "$1",
		Time:       time.Now(),
	}
	m := memory.New()
	year, week := fund.Time.ISOWeek()
	weeklyTransaction := domain.WeeklyTransaction{Year: year, Week: week}
	weeklyTransactionTotal := domain.WeeklyTransactionTotal{Value: 1}
	err := m.AddWeeklyTransaction(fund.CustomerID, weeklyTransaction, weeklyTransactionTotal)
	assert.Nil(t, err)
	weekly, err := m.GetWeeklyTransaction(fund.CustomerID, weeklyTransaction)
	assert.Nil(t, err)
	assert.Equal(t, float64(1), weekly.Value)
}

func TestGetWeeklyTransactionShouldReturnNotFound(t *testing.T) {
	fund := domain.Transaction{
		ID:         "123",
		CustomerID: "1234",
		LoadAmount: "$1",
		Time:       time.Now(),
	}
	m := memory.New()
	year, week := fund.Time.ISOWeek()
	weeklyTransaction := domain.WeeklyTransaction{Year: year, Week: week}
	weeklyTransactionTotal := domain.WeeklyTransactionTotal{Value: 1}
	err := m.AddWeeklyTransaction(fund.CustomerID, weeklyTransaction, weeklyTransactionTotal)
	assert.Nil(t, err)
	_, err = m.GetWeeklyTransaction("888", weeklyTransaction)
	assert.Equal(t, domain.ErrNotFound, err)
}
