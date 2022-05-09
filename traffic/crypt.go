package traffic

import (
	log "github.com/cihub/seelog"
	"github.com/xtaci/kcp-go"
	"tunn-hub/common/config"
)

//
// GetEncryptFP
// @Description:
//
func GetEncryptFP(cfg config.DataProcess, key []byte) FlowProcessor {
	if len(key) < 32 {
		_ = log.Warn("key too short")
		return nil
	}
	switch cfg.CipherType {
	case config.CipherTypeNone:
		return nil
	case config.CipherTypeAES256, config.CipherTypeAES192, config.CipherTypeAES128:
		//32位秘钥->AES-CBC-256
		var k []byte
		if cfg.CipherType == config.CipherTypeAES192 {
			//24位秘钥->AES-CBC-192
			k = key[:24]
		} else if cfg.CipherType == config.CipherTypeAES128 {
			//16位秘钥->AES-CBC-128
			k = key[:16]
		} else {
			k = key[:32]
		}
		crypt, err := kcp.NewAESBlockCrypt(k)
		if err != nil {
			_ = log.Warn("initialize crypt[", cfg.CipherType, "] failed ", err)
			return nil
		}
		log.Info("set encrypt cipher [", cfg.CipherType, "]")
		return NewBlockEncryptFP(crypt)
	case config.CipherTypeXOR:
		crypt, err := kcp.NewSimpleXORBlockCrypt(key)
		if err != nil {
			_ = log.Warn("initialize crypt[", cfg.CipherType, "] failed ", err)
			return nil
		}
		log.Info("set encrypt cipher [", cfg.CipherType, "]")
		return NewBlockEncryptFP(crypt)
	case config.CipherTypeTEA:
		crypt, err := kcp.NewTEABlockCrypt(key[:16])
		if err != nil {
			_ = log.Warn("initialize crypt[", cfg.CipherType, "] failed ", err)
			return nil
		}
		log.Info("set decrypt cipher [", cfg.CipherType, "]")
		return NewBlockEncryptFP(crypt)
	case config.CipherTypeXTEA:
		crypt, err := kcp.NewXTEABlockCrypt(key[:16])
		if err != nil {
			_ = log.Warn("initialize crypt[", cfg.CipherType, "] failed ", err)
			return nil
		}
		log.Info("set encrypt cipher [", cfg.CipherType, "]")
		return NewBlockEncryptFP(crypt)
	case config.CipherTypeSM4:
		crypt, err := kcp.NewSM4BlockCrypt(key[:16])
		if err != nil {
			_ = log.Warn("initialize crypt[", cfg.CipherType, "] failed ", err)
			return nil
		}
		log.Info("set encrypt cipher [", cfg.CipherType, "]")
		return NewBlockEncryptFP(crypt)
	case config.CipherTypeSalsa20:
		fp, err := NewSalsa20EncryptFP(key[:32])
		if err != nil {
			_ = log.Warn("initialize crypt[", cfg.CipherType, "] failed ", err)
			return nil
		}
		log.Info("set encrypt cipher [", cfg.CipherType, "]")
		return fp
	case config.CipherTypeBlowfish:
		//1-56
		crypt, err := kcp.NewBlowfishBlockCrypt(key[:32])
		if err != nil {
			_ = log.Warn("initialize crypt[", cfg.CipherType, "] failed ", err)
			return nil
		}
		log.Info("set encrypt cipher [", cfg.CipherType, "]")
		return NewBlockEncryptFP(crypt)
	}
	log.Info("unknown encrypt cipher type [", cfg.CipherType, "]")
	return nil
}

//
// GetDecryptFP
// @Description:
//
func GetDecryptFP(cfg config.DataProcess, key []byte) FlowProcessor {
	if len(key) < 32 {
		_ = log.Warn("key too short")
		return nil
	}
	switch cfg.CipherType {
	case config.CipherTypeNone:
		return nil
	case config.CipherTypeAES256, config.CipherTypeAES192, config.CipherTypeAES128:
		//32位秘钥->AES-CBC-256
		var k []byte
		if cfg.CipherType == config.CipherTypeAES192 {
			//24位秘钥->AES-CBC-192
			k = key[:24]
		} else if cfg.CipherType == config.CipherTypeAES128 {
			//16位秘钥->AES-CBC-128
			k = key[:16]
		} else {
			k = key[:32]
		}
		crypt, err := kcp.NewAESBlockCrypt(k)
		if err != nil {
			_ = log.Warn("initialize crypt[", cfg.CipherType, "] failed ", err)
			return nil
		}
		log.Info("set decrypt cipher [", cfg.CipherType, "]")
		return NewBlockDecryptFP(crypt)
	case config.CipherTypeXOR:
		crypt, err := kcp.NewSimpleXORBlockCrypt(key)
		if err != nil {
			_ = log.Warn("initialize crypt[", cfg.CipherType, "] failed ", err)
			return nil
		}
		log.Info("set decrypt cipher [", cfg.CipherType, "]")
		return NewBlockDecryptFP(crypt)
	case config.CipherTypeTEA:
		crypt, err := kcp.NewTEABlockCrypt(key[:16])
		if err != nil {
			_ = log.Warn("initialize crypt[", cfg.CipherType, "] failed ", err)
			return nil
		}
		log.Info("set decrypt cipher [", cfg.CipherType, "]")
		return NewBlockDecryptFP(crypt)
	case config.CipherTypeXTEA:
		crypt, err := kcp.NewXTEABlockCrypt(key[:16])
		if err != nil {
			_ = log.Warn("initialize crypt[", cfg.CipherType, "] failed ", err)
			return nil
		}
		log.Info("set decrypt cipher [", cfg.CipherType, "]")
		return NewBlockDecryptFP(crypt)
	case config.CipherTypeSM4:
		crypt, err := kcp.NewSM4BlockCrypt(key[:16])
		if err != nil {
			_ = log.Warn("initialize crypt[", cfg.CipherType, "] failed ", err)
			return nil
		}
		log.Info("set decrypt cipher [", cfg.CipherType, "]")
		return NewBlockDecryptFP(crypt)
	case config.CipherTypeSalsa20:
		fp, err := NewSalsa20DecryptFP(key[:32])
		if err != nil {
			_ = log.Warn("initialize crypt[", cfg.CipherType, "] failed ", err)
			return nil
		}
		log.Info("set decrypt cipher [", cfg.CipherType, "]")
		return fp
	case config.CipherTypeBlowfish:
		//1-56
		crypt, err := kcp.NewBlowfishBlockCrypt(key[:32])
		if err != nil {
			_ = log.Warn("initialize crypt[", cfg.CipherType, "] failed ", err)
			return nil
		}
		log.Info("set decrypt cipher [", cfg.CipherType, "]")
		return NewBlockDecryptFP(crypt)
	}
	log.Info("unknown decrypt cipher type [", cfg.CipherType, "]")
	return nil
}
