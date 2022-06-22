package main

import (
	"log"
	"net/http"
	"syscall"
	_ "net/http/pprof"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

var epoller *epoll

func handle(w http.ResponseWriter, r *http.Request) {
	conn, _,_, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		return
	}
	if err := epoller.Add(conn); err != nil {
		log.Printf("Failed to add connection")
		conn.Close()
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
	var err error
	epoller, err = NewEpoll()
	if err != nil {
		panic(err)
	}
	go Start()
	http.HandleFunc("/", handle)
	http.ListenAndServe(":8000", nil)
}

func Start() {
	for {
		connections, err := epoller.Wait()
		if err != nil {
			log.Printf("Failed to epoll wait %v", err)
			continue
		}
		for _, conn := range connections {
			if conn == nil {
				break
			}
		        msg, _,  err := wsutil.ReadClientData(conn)
			if err != nil {
				if err := epoller.Remove(conn); err != nil {
					log.Printf("Failed to remove %v", err)

				}
				conn.Close()
			} else {
				log.Printf("msg: %s", string(msg))
			}
		}
	}
}