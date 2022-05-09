package tunnel

import (
	"errors"
)

var (
	ErrDisconnect             = errors.New("disconnect")
	ErrDisconnectAccidentally = errors.New("disconnect accidentally")
	ErrStoppedByServer        = errors.New("client stopped")
)
