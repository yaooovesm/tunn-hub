package config

//
// Security
// @Description:
//
type Security struct {
	CertPem string `json:"cert"`
	KeyPem  string `json:"key"`
}
