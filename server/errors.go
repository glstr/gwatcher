package server

import "errors"

var (
	ErrNoHandler = errors.New("no invalid handler")
	ErrNoServer  = errors.New("no invalid server")
)
