package socks555

import (
	"io"
)

type clientAuthMessage struct {
	Version  byte
	NMethods byte
	Methods  []Method
}

type Method = byte

func NewClientAuthMessage(conn io.Reader) (*clientAuthMessage, error) {
	// Read version, nMethods
	buffer := make([]byte, 2)
	_, err := io.ReadFull(conn, buffer)
	if err != nil {
		return nil, err
	}

	// validate version
	if buffer[0] != SOCKS5VERSION {
		return nil, unsupportedVersion
	}

	// read methods
	nmethods := buffer[1]
	buffer = make([]byte, nmethods)
	_, err = io.ReadFull(conn, buffer)
	if err != nil {
		return nil, err
	}

	return &clientAuthMessage{
		Version:  SOCKS5VERSION,
		NMethods: nmethods,
		Methods:  buffer,
	}, nil
}
func NewServerAuthMessage(conn io.Writer, method Method) error {
	buf := []byte{SOCKS5VERSION, method}
	_, err := conn.Write(buf)
	return err
}
