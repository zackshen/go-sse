package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	sse "github.com/zackshen/go-sse"
)

func main() {

	handler := sse.NewSSEHandler()
	handler.Listen()
	go func() {
		for {
			time.Sleep(time.Second * 3)
			handler.Broadcast("hello")
		}
	}()

	fmt.Println("Listen and server at 0.0.0.0:3000: http://127.0.0.1:3000")
	http.Handle("/", http.FileServer(http.Dir("./html")))
	http.HandleFunc("/sse", handler.HttpHandler)
	err := http.ListenAndServe("0.0.0.0:3000", nil)
	if err != nil {
		log.Fatal("run server", err)
	}

}
