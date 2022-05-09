package model

//
// ServerStatus
// @Description:
//
type ServerStatus struct {
	RX LinkStatus `json:"rx"`
	TX LinkStatus `json:"tx"`
}
