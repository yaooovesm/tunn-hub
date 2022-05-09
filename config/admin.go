package config

//
// Admin
// @Description:
//
type Admin struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
	Https   bool   `json:"https"`
	DBFile  string `json:"db_file"`
}
