package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	sse "github.com/zackshen/go-sse"
)

func main() {
	pid := os.Getpid()

	handler := sse.NewSSEHandler()
	handler.Listen()
	go func() {
		for {
			time.Sleep(time.Second)

			loopFlag := "-bn"
			pidFlag := "-p"
			if runtime.GOOS == "darwin" {
				loopFlag = "-l"
				pidFlag = "-pid"
			}
			cmd := exec.Command("top", loopFlag, "1", pidFlag, strconv.Itoa(pid))
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				log.Fatal(err)
			}
			cmd.Start()
			bytes, err := ioutil.ReadAll(stdout)
			if err != nil {
				log.Fatal(err)
			}
			handler.Broadcast(string(bytes))
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
