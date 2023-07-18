package socks555

import (
	"fmt"
	"io"
	"log"
	"net"
)

type ClientRequestMessage struct {
	Version  byte
	Cmd      Command
	Reserved byte
	Address  string
	AddrType AddressType
	Port     uint16
}

type Command = byte

type AddressType = byte

func NewClientRequestMessage(conn io.ReadWriter) (*ClientRequestMessage, error) {

	buf := make([]byte, 4)
	if _, err := io.ReadFull(conn, buf); err != nil {
		log.Fatalf("can not read the first four format data, error: %s", err)
	}
	version, command, reserved, addressType := buf[0], buf[1], buf[2], buf[3]

	if version != SOCKS5VERSION {
		log.Printf("not socks5, illegal protocal")
		return nil, unsupportedVersion
	}
	if command < CmdConnect || command > CmdUDP {
		return nil, unsupportedCMD
	}
	if reserved != ReservedField {
		return nil, incorrectProtocolFormat
	}
	if addressType != TypeIPv4 && addressType != TypeIPv6 && addressType != TypeDomain {
		return nil, incorrectProtocolFormat
	}

	message := ClientRequestMessage{
		Cmd:      command,
		AddrType: addressType,
	}

	buf = make([]byte, IPv4Length)
	switch addressType {
	case TypeIPv6:
		buf = make([]byte, IPv6Length)
		fallthrough
	case TypeIPv4:
		if _, err := io.ReadFull(conn, buf); err != nil {
			log.Fatalf("error when reading address! %s\n", err)
		}
		ip := net.IP(buf)
		message.Address = ip.String()
	case TypeDomain:
		if _, err := io.ReadFull(conn, buf[:1]); err != nil {
			log.Fatalf("error when reading address! %s\n", err)
		}
		domainLength := buf[0]
		buf = make([]byte, domainLength)
		if _, err := io.ReadFull(conn, buf); err != nil {
			log.Fatalf("error when reading address! %s\n", err)
		}
		message.Address = string(buf)
	}
	if len(buf) < 2 {
		buf = make([]byte, PortLength)
	}
	if _, err := io.ReadFull(conn, buf[:PortLength]); err != nil {
		log.Printf("error when reading port")
	}
	message.Port = (uint16(buf[0]) << 8) + uint16(buf[1])
	return &message, nil
}

func WriteRequestSuccessMessage(conn io.Writer, ip net.IP, port uint16) error {
	addressType := TypeIPv4
	if byte(len(ip)) == IPv6Length {
		addressType = TypeIPv6
	}

	//write in success information
	_, err := conn.Write([]byte{SOCKS5VERSION, ReplySuccess, ReservedField, addressType})
	if err != nil {
		log.Printf("write in successful message failed! %s\n", err)
		return err
	}

	//write ip
	if _, err := conn.Write(ip); err != nil {
		log.Printf("write in ip address in the success message failed! %s\n", err)
		return err
	}

	// write port
	buf := make([]byte, PortLength)
	buf[0] = byte(port >> 8)
	buf[1] = byte(port & 0xff)
	fmt.Printf("the port is : %d and  %d\n", buf[0], buf[1])
	if _, err := conn.Write(buf); err != nil {
		log.Printf("write in port information in success messaged failed! %s\n", err)
		return err
	}
	return nil
}

func WriteRequestFailureMessage(conn io.Writer, replyType ReplyType) error {
	addressType := TypeIPv4

	//write in failure information
	_, err := conn.Write([]byte{SOCKS5VERSION, replyType, ReservedField, addressType})
	if err != nil {
		log.Printf("write in failure message failed! %s\n", err)
		return err
	}
	return nil
}
