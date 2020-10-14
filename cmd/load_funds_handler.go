package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/danielfmelo/load-funds-handler/handler"
	"github.com/danielfmelo/load-funds-handler/listener"
	"github.com/danielfmelo/load-funds-handler/storage/memory"
)

func main() {

	outputCh := make(chan []byte)
	errCh := make(chan []byte)
	inputCh := make(chan []byte)
	database := memory.New()
	handle := handler.New(database, outputCh, errCh)
	listening := listener.New(handle)

	listening.Receiver(inputCh)
	var wgOrderControl sync.WaitGroup
	var wgReadAllControl sync.WaitGroup
	readOutput(outputCh, errCh, &wgOrderControl, &wgReadAllControl)
	readFile(inputCh, &wgOrderControl, &wgReadAllControl)

	wgOrderControl.Wait()
}

func readFile(inputCh chan []byte, wgOrderControl *sync.WaitGroup, wgReadAllControl *sync.WaitGroup) {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		event := scanner.Text()
		wgOrderControl.Add(1)
		wgReadAllControl.Add(1)
		inputCh <- []byte(event)
		wgReadAllControl.Wait()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func readOutput(
	outputCh chan []byte,
	errCh chan []byte,
	wgOrderControl *sync.WaitGroup,
	wgReadAllControl *sync.WaitGroup,
) {
	go func() {
		for {
			select {
			case record := <-outputCh:
				fmt.Println(string(record))
				wgOrderControl.Done()
				wgReadAllControl.Done()
			case _ = <-errCh:
				wgOrderControl.Done()
				wgReadAllControl.Done()
			}
		}
	}()
}
