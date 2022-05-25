package authenticationv2

import "errors"

var (
	ErrCertFileNotFound = errors.New("cert file not found")
	ErrConnectFailed    = errors.New("connect failed")
	ErrDisconnect       = errors.New("disconnect")
	ErrKick             = errors.New("kick")
	ErrRestart          = errors.New("restart")
)

var (
	ErrAuthTimeout   = errors.New("timeout")
	ErrAuthConnect   = errors.New("connect to authentication server failed")
	ErrAuthBadPacket = errors.New("invalid authentication packet")
	ErrAuthFailed    = errors.New("authentication failed")
)
