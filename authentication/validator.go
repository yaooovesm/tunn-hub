package authentication

import (
	"errors"
	"tunn-hub/administration"
	"tunn-hub/config"
)

//
// IValidator
// @Description:
//
type IValidator interface {
	//
	// ValidateUser
	// @Description:
	// @param user
	// @return bool
	//
	ValidateUser(user config.User) error

	//
	// ValidateConfig
	// @Description:
	// @param config
	// @return error
	//
	ValidateConfig(cfg config.Config) error
}

//
// DefaultValidator
// @Description:
//
type DefaultValidator struct {
}

//
// ValidateUser
// @Description:
// @receiver DefaultValidator
// @param user
// @return error
//
func (DefaultValidator) ValidateUser(user config.User) error {
	return administration.UserServiceInstance().TunnelLogin(user.Account, user.Password)
}

//
// ValidateConfig
// @Description:
// @receiver DefaultValidator
// @param cfg
// @return error
//
func (DefaultValidator) ValidateConfig(cfg config.Config) error {
	//check protocol
	if cfg.Global.Protocol != config.Current.Global.Protocol {
		return errors.New("different from server's protocol")
	}
	//check data process method
	if cfg.DataProcess.CipherType != config.Current.DataProcess.CipherType {
		return errors.New("different from server's data process")
	}
	return nil
}
