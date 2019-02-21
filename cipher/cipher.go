package cipher

import (
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"io"
)

type DESCFBCipher struct {
	block cipher.Block
	rwc   io.ReadWriteCloser
	*cipher.StreamReader
	*cipher.StreamWriter
}

func NewDESCFBCipher(rwc io.ReadWriteCloser, password string) (*DESCFBCipher, error) {
	block, err := des.NewCipher([]byte(password))
	if err != nil {
		return nil, err
	}
	return &DESCFBCipher{
		block: block,
		rwc:   rwc,
	}, nil
}

func (des *DESCFBCipher) Read(b []byte) (int, error) {
	if des.StreamReader == nil {
		iv := make([]byte, des.block.BlockSize())
		n, err := io.ReadFull(des.rwc, iv)
		if err != nil {
			return n, err
		}
		stream := cipher.NewCFBDecrypter(des.block, iv)
		des.StreamReader = &cipher.StreamReader{
			S: stream,
			R: des.rwc,
		}
	}
	return des.StreamReader.Read(b)
}

func (des *DESCFBCipher) Write(b []byte) (int, error) {
	if des.StreamWriter == nil {
		iv := make([]byte, des.block.BlockSize())
		_, err := rand.Read(iv)
		if err != nil {
			return 0, err
		}
		stream := cipher.NewCFBEncrypter(des.block, iv)
		des.StreamWriter = &cipher.StreamWriter{
			S: stream,
			W: des.rwc,
		}
		n, err := des.rwc.Write(iv)
		if err != nil {
			return n, err
		}
	}
	return des.StreamWriter.Write(b)
}

func (des *DESCFBCipher) Close() error {
	if des.StreamWriter != nil {
		des.StreamWriter.Close()
	}
	if des.rwc != nil {
		des.rwc.Close()
	}
	return nil
}
