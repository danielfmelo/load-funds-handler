package storage

import (
	"github.com/danielfmelo/load-funds-handler/domain"
	"github.com/stretchr/testify/mock"
)

type StorageMock struct {
	mock.Mock
}

func (sm *StorageMock) AddTransaction(transaction domain.Transaction) error {
	args := sm.Called()
	return args.Error(0)
}

func (sm *StorageMock) AddDailyTransaction(customerID, day string, daily domain.DailyTransaction) error {
	args := sm.Called(customerID, day)
	return args.Error(0)
}

func (sm *StorageMock) AddWeeklyTransaction(customerID string, week domain.WeeklyTransaction, total domain.WeeklyTransactionTotal) error {
	args := sm.Called(customerID, week, total)
	return args.Error(0)
}

func (sm *StorageMock) GetDailyTransaction(customerID, day string) (domain.DailyTransaction, error) {
	args := sm.Called(customerID, day)
	return args.Get(0).(domain.DailyTransaction), args.Error(1)
}

func (sm *StorageMock) GetWeeklyTransaction(customerID string, week domain.WeeklyTransaction) (domain.WeeklyTransactionTotal, error) {
	args := sm.Called(customerID, week)
	return args.Get(0).(domain.WeeklyTransactionTotal), args.Error(1)
}
