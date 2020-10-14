package handler

import "github.com/stretchr/testify/mock"

type HandlerMock struct {
	mock.Mock
}

func (h *HandlerMock) Transaction(fund []byte) {
	h.Called(fund)
}
