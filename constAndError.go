package socks555

import "errors"

// here are errors
var (
	unsupportedCMD          = errors.New("the command not supported")
	unsupportedVersion      = errors.New("protocol version not supported")
	incorrectProtocolFormat = errors.New("incorrect protocol format")
	addressTypeError        = errors.New("unsupported addressType")
)

const SOCKS5VERSION = 0x05

const ReservedField = 0x00

const (
	MethodNoAuth   Method = 0x00
	MethodGSSAPI   Method = 0x01
	MethodPassword Method = 0x02
	MethodNoAccept Method = 0xff
)

const (
	CmdConnect Command = 0x01
	CmdBind    Command = 0x02
	CmdUDP     Command = 0x03
)
const (
	TypeIPv4   AddressType = 0x01
	TypeDomain AddressType = 0x03
	TypeIPv6   AddressType = 0x04
)

const (
	IPv4Length byte = 4
	IPv6Length byte = 16
	PortLength byte = 2
)

type ReplyType = byte

const (
	ReplySuccess ReplyType = iota
	ReplyServerFailure
	ReplyConnectionNotAllowed
	ReplyNetworkUnreachable
	ReplyHostUnreachable
	ReplyConnectionRefused
	ReplyTTLExpired
	ReplyCommandNotSupported
	ReplyAddressTypeNotSupported
)
