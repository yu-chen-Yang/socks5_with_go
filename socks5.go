package socks555

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

type Server interface {
	Run() error
}

type SOCKS5server struct {
	IP   string
	Port int
}

func (s *SOCKS5server) Run() error {
	//localhost:8160
	fmt.Printf("listening on :%s:%d\n", s.IP, s.Port)
	lis, err := net.Listen("tcp", ":1080")
	if err != nil {
		log.Fatalf(err.Error())
		return err
	}
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Fatalf("connection failed from %s:%s\n", conn.RemoteAddr(), err.Error())
			continue
		}
		//create a new go routine
		go func() {
			defer conn.Close()
			if err := handleConnection(conn); err != nil {
				log.Fatalf("handle connection failed from %s:%s\n", conn.RemoteAddr(), err.Error())
			}
		}()
	}
}

func handleConnection(conn net.Conn) error {

	//协商过程
	if err := Auth(conn); err != nil {
		return err
	}
	//请求过程
	targetConn, err := Request(conn)
	if err != nil {
		return err
	}
	//转发过程
	if err := forward(conn, targetConn); err != nil {
		return err
	}

	return nil
}

// Auth negotiate
func Auth(conn net.Conn) error {
	clientMessage, err := NewClientAuthMessage(conn)
	if err != nil {
		return err
	}
	log.Println(clientMessage.Version, clientMessage.NMethods, clientMessage.Methods)

	//no-auth
	acceptable := false
	for _, method := range clientMessage.Methods {
		if method == MethodNoAccept {
			acceptable = false
			break
		}
		if method == MethodNoAuth {
			acceptable = true
		}
	}
	if !acceptable {
		NewServerAuthMessage(conn, MethodNoAccept)
		return errors.New("method not supported")

	}
	return NewServerAuthMessage(conn, MethodNoAuth)
}

func forward(conn io.ReadWriter, targetConn io.ReadWriteCloser) error {
	defer targetConn.Close()
	go io.Copy(targetConn, conn)
	_, err := io.Copy(conn, targetConn)
	return err
}

// Request  to be continued, should add ipv4, udp and more failure handler
func Request(conn io.ReadWriter) (io.ReadWriteCloser, error) {

	message, err := NewClientRequestMessage(conn)
	if err != nil {
		return nil, err
	}
	//check if the command is supported
	if message.Cmd != CmdConnect {
		return nil, WriteRequestFailureMessage(conn, ReplyCommandNotSupported)
	}

	//check if the address type is supported
	if message.AddrType != TypeIPv4 {
		return nil, WriteRequestFailureMessage(conn, ReplyAddressTypeNotSupported)
	}

	//send request to the target server
	address := fmt.Sprintf("%s:%d", message.Address, message.Port)
	targetConn, err := net.Dial("tcp", address)

	//send connection failure reply
	if err != nil {
		log.Printf("can`t reach the target server! error: %s\n", err)
		return nil, WriteRequestFailureMessage(conn, ReplyConnectionRefused)
	}

	//send success reply
	addrValue := targetConn.LocalAddr()
	addr := addrValue.(*net.TCPAddr)
	return targetConn, WriteRequestSuccessMessage(conn, addr.IP, uint16(addr.Port))
}

//转发
