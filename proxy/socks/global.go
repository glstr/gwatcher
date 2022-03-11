package socks

import "errors"

var (
	ErrVersionFailed  = errors.New("socks5 version err")
	ErrMessageInvalid = errors.New("socks5 message invalid")
	ErrNotSupportCmd  = errors.New("not support cmd")
)
