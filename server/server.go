package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/yu2erer/GoSocksProxy/cipher"
)

func testToRead(conn *cipher.CipherConn) {
	defer conn.Close()
	buf := make([]byte, 263)
	n, err := io.ReadAtLeast(conn, buf, 5)
	if err != nil {
		return
	}

	var dstIP []byte
	switch buf[0] {
	case 0x01: // ipv4
		dstIP = buf[1 : net.IPv4len+1]
		fmt.Println("访问地址:", string(buf[1:net.IPv4len+1]))
	case 0x03: // domainname
		ipAddr, err := net.ResolveIPAddr("ip", string(buf[2:n-2]))
		if err != nil {
			log.Println("err ", err)
			log.Println(string(buf))
			return
		}
		fmt.Println("访问域名:", string(buf[2:n-2]))
		dstIP = ipAddr.IP
	case 0x04: // ipv6
		dstIP = buf[1 : net.IPv6len+1]
	default:
		return
	}

	dstPort := buf[n-2:]
	dstAddr := &net.TCPAddr{
		IP:   dstIP,
		Port: int(binary.BigEndian.Uint16(dstPort)),
	}
	client, err := net.DialTCP("tcp", nil, dstAddr)
	if err != nil {
		return
	}
	defer client.Close()
	client.SetLinger(0)

	go func() {
		defer conn.Close()
		defer client.Close()
		io.Copy(conn, client)
	}()
	io.Copy(client, conn)
}

func main() {
	tcpAddr := &net.TCPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 3344,
	}
	listener, _ := net.ListenTCP("tcp", tcpAddr)

	for {
		conn, _ := listener.AcceptTCP()
		dstServer, err := cipher.NewCipherConn(conn, "des", "haohaio!")
		if err != nil {
			continue
		}
		go testToRead(dstServer)
	}

}
