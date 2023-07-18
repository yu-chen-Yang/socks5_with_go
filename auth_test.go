package socks555

import (
	"bytes"
	"reflect"
	"testing"
)

func TestNewClientAuthMessage(t *testing.T) {
	t.Run("should generate a message", func(t *testing.T) {
		b := []byte{SOCKS5VERSION, 2, 0x00, 0x01}
		r := bytes.NewReader(b)
		mes, err := NewClientAuthMessage(r)
		if err != nil {
			t.Fatalf("failed %s", err)
		}

		if mes.Version != SOCKS5VERSION {
			t.Fatalf("version error %s", err)
		}
		if !reflect.DeepEqual(mes.Methods, []byte{0x00, 0x01}) {
			t.Fatalf("want %v, but got %v", []byte{0x00, 0x01}, mes.Methods)
		}
	})
}

func TestNewServerAuthMessage(t *testing.T) {
	t.Run("pass", func(t *testing.T) {
		var buf bytes.Buffer
		err := NewServerAuthMessage(&buf, MethodNoAuth)
		if err != nil {
			t.Fatalf("failed %s", err)
		}
		got := buf.Bytes()
		if !reflect.DeepEqual(got, []byte{SOCKS5VERSION, MethodNoAuth}) {
			t.Fatal("failed, not equal")
		}
	})
}
