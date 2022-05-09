package authentication

import "errors"

var (
	ErrCertFileNotFound = errors.New("cert file not found")
	ErrConnectFailed    = errors.New("connect failed")
)

var (
	ErrAuthTimeout   = errors.New("timeout")
	ErrAuthConnect   = errors.New("connect to authentication server failed")
	ErrAuthBadPacket = errors.New("invalid authentication packet")
	ErrAuthFailed    = errors.New("authentication failed")
)
