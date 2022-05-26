package authenticationv2

import (
	"encoding/json"
	"tunn-hub/administration"
	"tunn-hub/config"
	"tunn-hub/transmitter"
)

type OperationName string

const (
	OperationGetConfig           OperationName = "OperationGetConfig"
	OperationUpdateRoutes        OperationName = "OperationUpdateRoutes"
	OperationResetRoutes         OperationName = "ResetRoutes"
	OperationGetAvailableExports OperationName = "GetAvailableExports"
)

//
// OperationResult
// @Description:
//
type OperationResult struct {
	UUID   string      //连接uuid
	Result interface{} //结果
	Error  string      //错误
}

//
// Operation
// @Description:
//
type Operation struct {
	UUID      string                 //连接uuid
	User      string                 //发起人
	Operation OperationName          //操作
	Params    map[string]interface{} //参数
}

//
// reply
// @Description:
// @receiver o
// @param res
// @param err
// @return OperationResult
//
func (o *Operation) reply(res interface{}, err error) AuthReply {
	var errStr = ""
	var ok = true
	if err != nil {
		errStr = err.Error()
		ok = false
	}
	result := OperationResult{
		UUID:   o.UUID,
		Result: res,
		Error:  errStr,
	}
	marshal, _ := json.Marshal(result)
	return AuthReply{
		Ok:      ok,
		Error:   errStr,
		Message: string(marshal),
	}
}

//
// GetParams
// @Description:
// @receiver o
// @param key
//
func (o *Operation) GetParams(key string) interface{} {
	if o.Params != nil {
		return o.Params[key]
	}
	return nil
}

//
// Process
// @Description:
// @receiver o
// @param tunn
//
func (o *Operation) Process() AuthReply {
	switch o.Operation {
	case OperationGetAvailableExports:
		exports, err := administration.UserServiceInstance().AvailableExports()
		return o.reply(exports, err)
	case OperationUpdateRoutes:
		params := o.GetParams("routes")
		var routes []config.Route
		b, _ := json.Marshal(params)
		_ = json.Unmarshal(b, &routes)
		err := administration.UserServiceInstance().UpdateRoutesByConnectUUID(o.UUID, routes)
		return o.reply("", err)
	case OperationResetRoutes:
		err := administration.UserServiceInstance().ResetRoutesByConnectUUID(o.UUID)
		return o.reply("", err)
	case OperationGetConfig:
		cfg, err := administration.UserServiceInstance().GetConfigByConnectUUID(o.UUID)
		return o.reply(cfg, err)
	}
	return AuthReply{
		Ok:      false,
		Error:   "unknown operation",
		Message: "未知的操作",
	}
}

//
// HandleOperation
// @Description:
// @param conn
// @param packet
//
func (s *Server) HandleOperation(tunn *transmitter.Tunnel, packet *TransportPacket) {
	//1.检查是否登录
	if !s.CheckByUUID(packet.UUID) {
		s.reply(AuthReply{
			Ok:      false,
			Error:   "not login",
			Message: "操作失败，用户未登录",
		}, PacketTypeOperation, packet.UUID, tunn)
		return
	}
	//2.解包
	operation := Operation{}
	err := json.Unmarshal(packet.Payload, &operation)
	if err != nil {
		s.reply(AuthReply{
			Ok:      false,
			Error:   "unable to unpack",
			Message: "操作失败",
		}, PacketTypeOperation, packet.UUID, tunn)
		return
	}
	//3.执行并返回
	s.reply(operation.Process(), PacketTypeOperation, packet.UUID, tunn)
}
