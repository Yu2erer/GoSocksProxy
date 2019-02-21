package cipher

import (
	"io"
	"net"
	"strings"
)

type CipherConn struct {
	Conn net.Conn
	io.ReadWriteCloser
}

func NewCipherConn(conn net.Conn, cryptMethod string, password string) (*CipherConn, error) {
	var rwc io.ReadWriteCloser
	var err error
	switch strings.ToLower(cryptMethod) {
	default:
		rwc = conn
	case "des":
		rwc, err = NewDESCFBCipher(conn, password)
	}
	if err != nil {
		return nil, err
	}
	return &CipherConn{
		Conn:            conn,
		ReadWriteCloser: rwc,
	}, nil
}
