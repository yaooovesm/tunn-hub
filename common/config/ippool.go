package config

//
// IPPool
// @Description:
//
type IPPool struct {
	Start   string `json:"start"`
	End     string `json:"end"`
	Network string `json:"network"`
}
