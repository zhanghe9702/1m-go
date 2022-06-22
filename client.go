package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var (
	ip = flag.String("ip", "127.0.0.1", "server ip")
	connections = flag.Int("conn", 1, "number of websocket connections")
)

func main() {
	flag.Usage = func ()  {
		io.WriteString(os.Stderr, `Websockets client generator
		Example usage: ./client -ip=172.17.0.1 -conn=10
		`)
		flag.PrintDefaults()
	}
	flag.Parse()
	u := url.URL{
		Scheme: "ws",
		Host: *ip + ":8000",
		Path: "/",
	}
	log.Printf("Connecting to %s", u.String())
	var conns []*websocket.Conn
	for i := 0; i < *connections; i++ {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			fmt.Println("Failed to connect", i, err)
			break
		}
		conns = append(conns, c)
		defer func() {
			c.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseAbnormalClosure, ""), time.Now().Add(time.Second))
			time.Sleep(time.Second)
			c.Close()
		}()
	}
	log.Printf("Finished initilizing %d connections", len(conns))
	tts := time.Second
	if *connections > 1000 {
		tts = time.Microsecond * 5
	}

	for {
		for i := 0; i < len(conns); i++ {
			time.Sleep(tts)
			conn := conns[i]
			log.Printf("Conn %d sending message", i)
			if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second*5)); err != nil {
				fmt.Printf("Failed to receive pong: %v", err)
			}
			conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Hello from conn %v", i)))
		}
	}

}