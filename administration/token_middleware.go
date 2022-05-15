package administration

import (
	"errors"
	"github.com/gin-gonic/gin"
)

var permissionAssignMap = map[string]UserAuthLevel{
	//api
	//user
	"/api/v1/user/:account":        User,
	"/api/v1/user/info/:id":        User,
	"/api/v1/user/status/:id":      User,
	"/api/v1/user/id/:id":          User,
	"/api/v1/user/list":            Administrator,
	"/api/v1/user/list/info":       Administrator,
	"/api/v1/user/create":          Administrator,
	"/api/v1/user/disable/:id":     Administrator,
	"/api/v1/user/update/:id":      User,
	"/api/v1/user/delete":          Administrator,
	"/api/v1/user/login":           None,
	"/api/v1/user/logout/:account": User,
	"/api/v1/user/general":         Administrator,
	"/api/v1/user/config/:id":      User,
	//token
	"/api/v1/token/":      None,
	"/api/v1/token/check": None,
	//server
	"/api/v1/server/version":           None,
	"/api/v1/server/flow":              Administrator,
	"/api/v1/server/config":            Administrator,
	"/api/v1/server/disconnect/id/:id": User,
	"/api/v1/server/reconnect/id/:id":  User,
	"/api/v1/server/ippool":            Administrator,
	"/api/v1/server/ippool/list":       Administrator,
	"/api/v1/server/monitor":           Administrator,
	//config
	"/api/v1/cfg/:id":             User,
	"/api/v1/cfg/list":            Administrator,
	"/api/v1/cfg/delete/:id":      Administrator,
	"/api/v1/cfg/create":          Administrator,
	"/api/v1/cfg/update":          User,
	"/api/v1/cfg/route/available": User,
}

//通过账号验证是否为对应用户操作
var selfPermissionAccountMap = map[string]int{
	"/api/v1/user/:account":        1,
	"/api/v1/user/logout/:account": 1,
}

//通过ID验证是否为对应用户操作
var selfPermissionIdMap = map[string]int{
	"/api/v1/user/info/:id":            1,
	"/api/v1/user/status/:id":          1,
	"/api/v1/user/id/:id":              1,
	"/api/v1/user/update/:id":          1,
	"/api/v1/user/config/:id":          1,
	"/api/v1/server/disconnect/id/:id": 1,
	"/api/v1/server/reconnect/id/:id":  1,
}

//
// AuthCheckMiddleWare
// @Description:
// @return gin.HandlerFunc
//
func AuthCheckMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//token检查
		if level, ok := permissionAssignMap[ctx.FullPath()]; ok {
			if level != None {
				err := TokenServiceInstance().CheckHeader(ctx, level)
				if err != nil {
					responseError(ctx, err, "没有操作权限")
					ctx.Abort()
					return
				}
			}
		} else {
			responseError(ctx, errors.New("permission denied"), "禁止访问")
			ctx.Abort()
		}
		//ID对象检查
		if _, ok := selfPermissionIdMap[ctx.FullPath()]; ok {
			err := TokenServiceInstance().CheckSelfById(ctx, ctx.Param("id"))
			if err != nil {
				responseError(ctx, err, "没有操作权限")
				ctx.Abort()
				return
			}
		}
		//Account对象检查
		if _, ok := selfPermissionAccountMap[ctx.FullPath()]; ok {
			err := TokenServiceInstance().CheckSelfByAccount(ctx, ctx.Param("account"))
			if err != nil {
				responseError(ctx, err, "没有操作权限")
				ctx.Abort()
				return
			}
		}
	}
}
