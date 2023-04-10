//Laboratory 3
//For debuggin, run the prog in the terminal "go run TokenRing.go"

package main

import (
	"fmt"
	"sync"
)

type token struct {
	data      string
	recipient int
	ttl       int
}

// definition of goroutine group
var wg sync.WaitGroup

func node(thread_index int, chan_resv chan token, chan_send chan token) {
	defer wg.Done()

	// waiting for a some signal
	for {
		message := <-chan_resv
		if message.data == "close" {
			break
		}
		//signal received
		fmt.Println("Узел", thread_index, "получил сообшение")
		if message.recipient == thread_index {
			fmt.Println("Узел ", thread_index, ":", message.data)
		} else {
			if message.ttl > 0 {
				message.ttl -= 1
				chan_send <- message
			} else {
				fmt.Println("В узле", thread_index, "сообшение истекло ")
			}
		}
	}
}

func main() {
	fmt.Println("ПРОТОКОЛ TokenRing")
	fmt.Println("Вводите каличество N узлов")
	var threads_number int
	fmt.Scanln(&threads_number)

	// creating channels, 0 channel is reserved for talking with the main thread
	var channels []chan token = make([]chan token, threads_number)
	var quit = make(chan int)

	// creating threads
	for i := 0; i < threads_number; i++ {
		channels[i] = make(chan token)
		wg.Add(1)
	}

	for i := 0; i < threads_number; i++ {
		if i == threads_number-1 {

			// the message hasn't set to other channels, is that we sending it to the main (same) chan.
			go node(i, channels[i], channels[i])
		} else {
			//sending the message to other channels
			go node(i, channels[i], channels[i+1])
		}
	}

	var new_token token

	for threads_number > 0 {
		fmt.Println("Вводите сообщение или close чтобы заканчивать заупск программы")
		fmt.Scanln(&new_token.data)
		if new_token.data == "close" {
			channels[0] <- new_token
			break
		}
		fmt.Println("Вводите id-recipient")
		fmt.Scanln(&new_token.recipient)
		for {
			if new_token.recipient > len(channels)-1 {
				fmt.Println("Error, Вводите целое число (начиная с 1 до ", len(channels))
				fmt.Scanln(&new_token.recipient)
			} else {
				break
			}
		}

		fmt.Println(" message timeout")
		fmt.Scanln(&new_token.ttl)

		channels[0] <- new_token
	}

	for i := 1; i < threads_number; i++ {
		channels[i] <- new_token
	}

	wg.Wait()

	// creating threads
	for i := 1; i < threads_number; i++ {
		close(channels[i])
	}

	close(quit)
}
