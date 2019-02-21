package socket

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/yu2erer/GoSocksProxy/cipher"
)

type Socks struct {
	ListenAddr  *net.TCPAddr
	ServerAddr  *net.TCPAddr
	CryptMethod string
	Password    string
}

func NewSocks(listenAddr, serverAddr *net.TCPAddr, cryptMethod, password string) *Socks {
	return &Socks{
		ListenAddr:  listenAddr,
		ServerAddr:  serverAddr,
		CryptMethod: cryptMethod,
		Password: password,
	}
}

func (s *Socks) Listen() {
	listener, err := net.ListenTCP("tcp", s.ListenAddr)
	if err != nil {
		log.Println(err)
		return
	}

	defer listener.Close()
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}
		conn.SetLinger(0)
		go s.handleRequest(conn)
	}
}

func (s *Socks) handleRequest(conn *net.TCPConn) {
	defer conn.Close()
	buf := make([]byte, 263)
	n, err := io.ReadAtLeast(conn, buf, 2)
	if err != nil {
		return
	}

	// 判断 socks 版本 只支持 socks5
	if buf[0] != 0x05 {
		return
	}

	nmethod := int(buf[1])
	msgLen := nmethod + 2

	if n < msgLen {
		if _, err = io.ReadFull(conn, buf[n:msgLen]); err != nil {
			return
		}
	} else if n > msgLen {
		return
	}

	/*
		告诉客户端 不需要验证
		+----+--------+
		|VER | METHOD |
		+----+--------+
		| 1  |   1    |
		+----+--------+
	*/
	conn.Write([]byte{0x05, 0x00})
	if n, err = io.ReadAtLeast(conn, buf, 5); err != nil {
		return
	}
	if buf[0] != 0x05 {
		return
	}
	if buf[1] != 0x01 {
		return
	}

	/*
				n 应当大于 7
		        +----+-----+-------+------+----------+----------+
		        |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
		        +----+-----+-------+------+----------+----------+
		        | 1  |  1  | X'00' |  1   | Variable |    2     |
		        +----+-----+-------+------+----------+----------+
	*/
	dstServer, err := s.DialServer()
	if err != nil {
		log.Println(err)
		return
	}

	defer dstServer.Close()

	/*
		直接转发 socks5协议的半截过去 0x01 为 ipv4 0x03 为domainname 0x04 为ipv6
		后面几位不定长为要访问的地址
		最后两位为端口号
	*/
	dstServer.Write(buf[3:n])

	conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	go func() {
		defer conn.Close()
		defer dstServer.Close()
		io.Copy(conn, dstServer)
	}()
	io.Copy(dstServer, conn)
}

func (s *Socks) DialServer() (*cipher.CipherConn, error) {
	conn, err := net.DialTCP("tcp", nil, s.ServerAddr)
	if err != nil {
		return nil, fmt.Errorf("Can't connect to server: %s, err: %s", s.ServerAddr, err)
	}
	dstServer, err := cipher.NewCipherConn(conn, s.CryptMethod, s.Password)
	if err != nil {
		return nil, err
	}
	return dstServer, nil
}
