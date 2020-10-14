package listener

import "github.com/danielfmelo/load-funds-handler/handler"

type Transaction struct {
	handle handler.HandlerTransaction
}

func New(handle handler.HandlerTransaction) *Transaction {
	return &Transaction{
		handle: handle,
	}
}

func (t *Transaction) Receiver(chFunds chan []byte) {
	go func() {
		for {
			select {
			case record := <-chFunds:
				t.handle.Transaction(record)
			}
		}
	}()
}
