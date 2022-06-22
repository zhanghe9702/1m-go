package main

import (
	"reflect"
	"sync"

	"github.com/gorilla/websocket"
	"golang.org/x/sys/unix"
)

type epoll struct {
	fd int
	connections map[int]*websocket.Conn
	lock sync.RWMutex
}

func NewEpoll() (*epoll, error) {
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	return &epoll{
		fd: fd,
		lock: sync.RWMutex{},
		connections: make(map[int]*websocket.Conn),
	},nil
}

func getWebSocketFD(conn *websocket.Conn) int {
	tcpConn := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn")
	fd := reflect.Indirect(tcpConn.FieldByName("fd")).FieldByName("pfd").FieldByName("Sysfd").Int()
	return int(fd)
}

func (e *epoll) Add(conn *websocket.Conn) error {

	return nil
}