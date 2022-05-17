package config

//
// Admin
// @Description:
//
type Admin struct {
	Address      string `json:"address"`
	Port         int    `json:"port"`
	ReporterPort int    `json:"reporter"`
	Https        bool   `json:"https"`
	DBFile       string `json:"db_file"`
}
