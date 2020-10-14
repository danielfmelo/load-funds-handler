package handler_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/danielfmelo/load-funds-handler/domain"

	"github.com/stretchr/testify/assert"

	"github.com/danielfmelo/load-funds-handler/handler"

	"github.com/danielfmelo/load-funds-handler/storage"
)

type handlerTest struct {
	repo *storage.StorageMock
}

func newSuite() *handlerTest {
	return &handlerTest{repo: &storage.StorageMock{}}
}

func fakeTransaction(t *testing.T, amount string) (domain.Transaction, []byte) {
	fund := domain.Transaction{
		ID:         "123",
		CustomerID: "321",
		LoadAmount: "$" + amount,
		Time:       time.Now(),
	}
	transaction, err := json.Marshal(fund)
	assert.Nil(t, err)
	fund.Time.Format(domain.DateLayout)
	return fund, transaction
}

func TestTransactionShouldReceiveUnmarshalError(t *testing.T) {
	suite := newSuite()
	chOut := make(chan []byte, 1)
	chErr := make(chan []byte, 1)
	h := handler.New(suite.repo, chOut, chErr)
	fund := []byte("with error")
	h.Transaction(fund)
	record := <-chErr
	errExpected := "msg: error to unmarshal fund with error error: invalid character 'w' looking for beginning of value"
	assert.Equal(t, errExpected, string(record))
}

func TestTransactionShouldReceiveStorageGetDailyTransactionError(t *testing.T) {
	suite := newSuite()
	chOut := make(chan []byte, 1)
	chErr := make(chan []byte, 1)
	h := handler.New(suite.repo, chOut, chErr)
	transaction, fund := fakeTransaction(t, "100")
	fakeErr := errors.New("some error")
	fakeDaily := domain.DailyTransaction{}
	suite.repo.On("AddTransaction").Return(nil).Once()
	suite.repo.On("GetDailyTransaction", transaction.CustomerID, transaction.Time.Format(domain.DateLayout)).Return(fakeDaily, fakeErr).Once()
	h.Transaction(fund)
	record := <-chErr
	errExpected := "msg: error to validate transaction per day error: some error"
	assert.Equal(t, string(record), errExpected)
}

func TestTransactionShouldReceiveInvalidDailyAmount(t *testing.T) {
	suite := newSuite()
	chOut := make(chan []byte, 1)
	chErr := make(chan []byte, 1)
	h := handler.New(suite.repo, chOut, chErr)
	transaction, fund := fakeTransaction(t, "2500.01")
	fakeDaily := domain.DailyTransaction{DailyTotal: 2500.00}
	suite.repo.On("AddTransaction").Return(nil).Once()
	suite.repo.On("GetDailyTransaction", transaction.CustomerID, transaction.Time.Format(domain.DateLayout)).Return(fakeDaily, nil).Once()
	h.Transaction(fund)
	record := <-chOut
	msgExpected := "{\"id\":\"123\",\"customer_id\":\"321\",\"accepted\":false}"
	assert.Equal(t, string(record), msgExpected)
}

func TestTransactionShouldReceiveValidDailyAmount(t *testing.T) {
	suite := newSuite()
	chOut := make(chan []byte, 1)
	chErr := make(chan []byte, 1)
	h := handler.New(suite.repo, chOut, chErr)
	transaction, fund := fakeTransaction(t, "2500.00")
	fakeDaily := domain.DailyTransaction{DailyTotal: 2500.00}
	year, week := transaction.Time.ISOWeek()
	fakeWeeklyTransaction := domain.WeeklyTransaction{Year: year, Week: week}
	fakeWeeklyTotal := domain.WeeklyTransactionTotal{Value: fakeDaily.DailyTotal}
	day := transaction.Time.Format(domain.DateLayout)
	suite.repo.On("GetDailyTransaction", transaction.CustomerID, day).Return(fakeDaily, nil).Once()
	suite.repo.On("GetWeeklyTransaction", transaction.CustomerID, fakeWeeklyTransaction).Return(fakeWeeklyTotal, nil).Once()
	suite.repo.On("AddTransaction").Return(nil).Once()
	suite.repo.On("AddDailyTransaction", transaction.CustomerID, day).Return(nil).Once()
	weeklyTotalExpected := domain.WeeklyTransactionTotal{Value: 5000}
	suite.repo.On("AddWeeklyTransaction", transaction.CustomerID, fakeWeeklyTransaction, weeklyTotalExpected).Return(nil).Once()
	h.Transaction(fund)
	record := <-chOut
	msgExpected := "{\"id\":\"123\",\"customer_id\":\"321\",\"accepted\":true}"
	assert.Equal(t, string(record), msgExpected)
}

func TestTransactionShouldReceiveValidDailyAmountTwice(t *testing.T) {
	suite := newSuite()
	chOut := make(chan []byte, 1)
	chErr := make(chan []byte, 1)
	h := handler.New(suite.repo, chOut, chErr)
	transaction, fund := fakeTransaction(t, "2500.00")
	fakeDaily := domain.DailyTransaction{DailyTotal: 2500.00}
	year, week := transaction.Time.ISOWeek()
	fakeWeeklyTransaction := domain.WeeklyTransaction{Year: year, Week: week}
	fakeWeeklyTotal := domain.WeeklyTransactionTotal{Value: fakeDaily.DailyTotal}
	day := transaction.Time.Format(domain.DateLayout)
	suite.repo.On("GetDailyTransaction", transaction.CustomerID, day).Return(fakeDaily, nil).Once()
	suite.repo.On("GetWeeklyTransaction", transaction.CustomerID, fakeWeeklyTransaction).Return(fakeWeeklyTotal, nil).Once()
	suite.repo.On("AddTransaction").Return(nil).Once()
	suite.repo.On("AddDailyTransaction", transaction.CustomerID, day).Return(nil).Once()
	weeklyTotalExpected := domain.WeeklyTransactionTotal{Value: 5000}
	suite.repo.On("AddWeeklyTransaction", transaction.CustomerID, fakeWeeklyTransaction, weeklyTotalExpected).Return(nil).Once()
	h.Transaction(fund)
	record := <-chOut
	msgExpected := "{\"id\":\"123\",\"customer_id\":\"321\",\"accepted\":true}"
	assert.Equal(t, string(record), msgExpected)
}

func TestTransactionShouldReceiveInvalidDailyTransactionCount(t *testing.T) {
	suite := newSuite()
	chOut := make(chan []byte, 1)
	chErr := make(chan []byte, 1)
	h := handler.New(suite.repo, chOut, chErr)
	transaction, fund := fakeTransaction(t, "2")
	fakeDaily := domain.DailyTransaction{TransactionCount: 3}
	suite.repo.On("AddTransaction").Return(nil).Once()
	suite.repo.On("GetDailyTransaction", transaction.CustomerID, transaction.Time.Format(domain.DateLayout)).Return(fakeDaily, nil).Once()
	h.Transaction(fund)
	record := <-chOut
	msgExpected := "{\"id\":\"123\",\"customer_id\":\"321\",\"accepted\":false}"
	assert.Equal(t, string(record), msgExpected)
}

func TestTransactionShouldReceiveValidDailyTransactionCount(t *testing.T) {
	suite := newSuite()
	chOut := make(chan []byte, 1)
	chErr := make(chan []byte, 1)
	h := handler.New(suite.repo, chOut, chErr)
	transaction, fund := fakeTransaction(t, "2500.00")
	fakeDaily := domain.DailyTransaction{TransactionCount: 2}
	year, week := transaction.Time.ISOWeek()
	fakeWeeklyTransaction := domain.WeeklyTransaction{Year: year, Week: week}
	fakeWeeklyTotal := domain.WeeklyTransactionTotal{Value: 0}
	day := transaction.Time.Format(domain.DateLayout)
	suite.repo.On("GetDailyTransaction", transaction.CustomerID, day).Return(fakeDaily, nil).Once()
	suite.repo.On("GetWeeklyTransaction", transaction.CustomerID, fakeWeeklyTransaction).Return(fakeWeeklyTotal, nil).Once()
	suite.repo.On("AddTransaction").Return(nil).Once()
	suite.repo.On("AddDailyTransaction", transaction.CustomerID, day).Return(nil).Once()
	weeklyTotalExpected := domain.WeeklyTransactionTotal{Value: 2500}
	suite.repo.On("AddWeeklyTransaction", transaction.CustomerID, fakeWeeklyTransaction, weeklyTotalExpected).Return(nil).Once()
	h.Transaction(fund)
	record := <-chOut
	msgExpected := "{\"id\":\"123\",\"customer_id\":\"321\",\"accepted\":true}"
	assert.Equal(t, string(record), msgExpected)
}

func TestTransactionShouldReceiveInvalidWeeklyAmount(t *testing.T) {
	suite := newSuite()
	chOut := make(chan []byte, 1)
	chErr := make(chan []byte, 1)
	h := handler.New(suite.repo, chOut, chErr)
	transaction, fund := fakeTransaction(t, "2500.00")
	fakeDaily := domain.DailyTransaction{TransactionCount: 2}
	year, week := transaction.Time.ISOWeek()
	fakeWeeklyTransaction := domain.WeeklyTransaction{Year: year, Week: week}
	fakeWeeklyTotal := domain.WeeklyTransactionTotal{Value: 17501}
	day := transaction.Time.Format(domain.DateLayout)
	suite.repo.On("AddTransaction").Return(nil).Once()
	suite.repo.On("GetDailyTransaction", transaction.CustomerID, day).Return(fakeDaily, nil).Once()
	suite.repo.On("GetWeeklyTransaction", transaction.CustomerID, fakeWeeklyTransaction).Return(fakeWeeklyTotal, nil).Once()
	h.Transaction(fund)
	record := <-chOut
	msgExpected := "{\"id\":\"123\",\"customer_id\":\"321\",\"accepted\":false}"
	assert.Equal(t, string(record), msgExpected)
}
