package main

import (
	"log"
	"net/http"
	"sync/atomic"
	"syscall"
	_ "net/http/pprof"
	"github.com/gorilla/websocket"
)
var count int64
func ws(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	n := atomic.AddInt64(&count, 1)
	if n % 100 == 0 {
		log.Printf("Total number of connections: %v", n)
	}
	defer func ()  {
		n := atomic.AddInt64(&count, -1)
		if n%100 == 0 {
			log.Printf("Total number of connections: %v", n)
		}
		conn.Close()
	}()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Failed to read message %v", err)
			conn.Close()
			return
		}
		log.Println(string(msg))
	}
}

func main() {
	var rlimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlimit); err != nil {
		panic(err)
	}
	rlimit.Cur = rlimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rlimit); err != nil {
		panic(err)
	}
	log.Printf("now we could open %d sockets", rlimit.Max)
	// Enable pprof hooks
	go func() {
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatalf("Pprof failed: %v", err)
		}
	}()
	http.HandleFunc("/", ws)
	http.ListenAndServe(":8000", nil)
}