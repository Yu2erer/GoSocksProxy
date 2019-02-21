package main

import (
	"net"

	"github.com/yu2erer/GoSocksProxy/socket"
)

func main() {
	listenAddr := &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 1234,
	}
	serverAddr := &net.TCPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 3344,
	}
	localConn := socket.NewSocks(listenAddr, serverAddr, "des", "haohaio!")
	localConn.Listen()
}
