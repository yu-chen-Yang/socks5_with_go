package socks555

import (
	"fmt"
	"log"
	"testing"
)

func TestRequest(t *testing.T) {
	var tests = []testType{
		{
			Cmd:      CmdConnect,
			AddrType: TypeIPv4,
			Address:  []byte{61, 173, 85, 238},
			Port:     []byte{0x7e, 0x90},
			OutPut: ClientRequestMessage{
				Cmd:      CmdConnect,
				Address:  "61.173.85.238",
				Port:     32400,
				Reserved: ReservedField,
				Version:  SOCKS5VERSION,
			},
		},
	}
	t.Run("should pass!", func(t *testing.T) {
		for _, test := range tests {
			message, err := Request(createBytes(test))
			if err != nil {
				log.Fatalf("error! test not pass!")
			}
			fmt.Printf("the return is :%v\n", message)
		}
	})
}
