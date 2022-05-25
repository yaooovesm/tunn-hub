package authenticationv2

import "encoding/json"

//
// AuthReply
// @Description:
//
type AuthReply struct {
	Ok      bool
	Error   string
	Message string
	Bytes   []byte
}

//
// Encode
// @Description:
// @receiver r
//
func (r *AuthReply) Encode() []byte {
	reply := r
	marshal, _ := json.Marshal(reply)
	r.Bytes = marshal
	return r.Bytes
}

//
// Decode
// @Description:
// @receiver r
//
func (r *AuthReply) Decode() error {
	rep := &AuthReply{}
	err := json.Unmarshal(r.Bytes, rep)
	if err != nil {
		return err
	}
	r.Ok = rep.Ok
	r.Error = rep.Error
	r.Message = rep.Message
	return nil
}
