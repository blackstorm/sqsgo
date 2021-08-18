package main

import (
	"fmt"
	"net/http"
	"time"
)

type message struct {
	value interface{}
}

func main() {
	queue := make(chan message)

	go func(ch chan message) {
		i := 0
		for {
			i++
			msg := fmt.Sprintf("No %d test message", i)
			fmt.Println(msg)
			queue <- message{
				value: msg,
			}
			time.Sleep(time.Second * 2)
		}
	}(queue)

	http.ListenAndServe(":8080", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		for {
			select {
			case msg := <-queue:
				rw.Write([]byte(fmt.Sprintf("get message %s", msg)))
				return
			case <-time.After(time.Second * 5):
				rw.Write([]byte("no message"))
				return
			}
		}
	}))
}
