package listener_test

import (
	"testing"

	"github.com/danielfmelo/load-funds-handler/handler"
	"github.com/danielfmelo/load-funds-handler/listener"
)

type transactionListenerSuite struct {
	handle *handler.HandlerMock
}

func newSuite() transactionListenerSuite {
	return transactionListenerSuite{
		handle: &handler.HandlerMock{},
	}
}

func TestReceiver(t *testing.T) {
	suite := newSuite()
	record := []byte("some data")
	suite.handle.On("Transaction", record).Return().Once()
	ch := make(chan []byte)
	lf := listener.New(suite.handle)
	lf.Receiver(ch)
	ch <- record
}
