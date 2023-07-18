package socks555

import (
	"bytes"
	"log"
	"net"
	"reflect"
	"testing"
)

type testType struct {
	Cmd      Command
	AddrType AddressType
	Address  []byte
	Port     []byte
	OutPut   ClientRequestMessage
}

func createBytes(t testType) *bytes.Buffer {
	var buf bytes.Buffer
	buf.Write([]byte{SOCKS5VERSION, t.Cmd, ReservedField, t.AddrType})
	buf.Write(t.Address)
	buf.Write(t.Port)
	return &buf
}

func TestNewClientRequestMessage(t *testing.T) {
	var tests = []testType{
		{
			Cmd:      CmdConnect,
			AddrType: TypeIPv4,
			Address:  []byte{192, 168, 0, 105},
			Port:     []byte{0, 80},
			OutPut: ClientRequestMessage{
				Cmd:      CmdConnect,
				Address:  "192.168.0.105",
				Port:     80,
				Reserved: ReservedField,
				Version:  SOCKS5VERSION,
			},
		},
	}

	t.Run("pass", func(t *testing.T) {
		for _, test := range tests {
			message, err := NewClientRequestMessage(createBytes(test))
			if err != nil {
				log.Fatalf("error! test not pass!")
			}
			if reflect.DeepEqual(message, test.OutPut) {
				log.Fatalf("error, should get %v, got %v", test.OutPut, message)
			}
		}
	})
}

type successMessageValidation struct {
	AddrType AddressType
	Addr     []byte
	Port     []byte
	Version  byte
	Cmd      byte
	Reserved byte
}

func TestWriteRequestSuccessMessage(t *testing.T) {
	tests := []struct {
		Ip    net.IP
		Port  uint16
		Valid successMessageValidation
	}{
		{
			Ip:   net.IP{192, 168, 0, 106},
			Port: 1080,
			Valid: successMessageValidation{
				Addr:     []byte{192, 168, 0, 106},
				AddrType: TypeIPv4,
				Port:     []byte{4, 56},
				Cmd:      ReplySuccess,
				Reserved: ReservedField,
				Version:  SOCKS5VERSION,
			},
		},
	}

	t.Run("pass", func(t *testing.T) {
		for _, test := range tests {
			var data bytes.Buffer
			valid := []byte{test.Valid.Version, test.Valid.Cmd, test.Valid.Reserved, test.Valid.AddrType, test.Valid.Addr[0],
				test.Valid.Addr[1], test.Valid.Addr[2], test.Valid.Addr[3], test.Valid.Port[0], test.Valid.Port[1]}

			if err := WriteRequestSuccessMessage(&data, test.Ip, test.Port); err != nil {
				log.Fatalf("test not passed! %s\n", err)
			}
			if !reflect.DeepEqual(valid, data.Bytes()) {
				log.Fatalf("test no tpassed! The result does not match! should be %v, got %v\n", valid, data.Bytes())
			}

		}
	})
}
