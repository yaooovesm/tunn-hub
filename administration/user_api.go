package administration

import (
	"errors"
	"github.com/gin-gonic/gin"
	"tunn-hub/administration/model"
	"tunn-hub/config"
)

//
// ApiUserLogin
// @Description:
// @param ctx
//
func ApiUserLogin(ctx *gin.Context) {
	info := model.UserInfo{}
	err := ctx.BindJSON(&info)
	if info.Account == "" || info.Password == "" {
		responseError(ctx, err, "用户名或密码不能为空")
		return
	}
	if err != nil {
		response400(ctx)
		return
	}
	token, err := UserServiceInstance().Login(ctx, &info)
	if err != nil {
		responseError(ctx, err, "登录失败")
		return
	}
	//去除敏感信息
	info.RemoveSensitive()
	res := map[string]interface{}{
		"token":    token,
		"info":     info,
		"reporter": config.Current.Admin.ReporterPort,
	}
	responseSuccess(ctx, res, "登录成功")
}

//
// ApiUserLogout
// @Description:
// @param ctx
//
func ApiUserLogout(ctx *gin.Context) {
	token := ctx.GetHeader("token")
	if token == "" {
		responseError(ctx, errors.New(""), "invalid token")
		return
	}
	TokenServiceInstance().RemoveByCode(token)
	responseSuccess(ctx, "", "退出成功")
}

//
// ApiGetUserStatusById
// @Description:
// @param ctx
//
func ApiGetUserStatusById(ctx *gin.Context) {
	status := UserServiceInstance().statusService.GetStatus(ctx.Param("id"))
	responseSuccess(ctx, status, "")
}

//
// ApiGetUserInfoById
// @Description:
// @param ctx
//
func ApiGetUserInfoById(ctx *gin.Context) {
	info, err := UserServiceInstance().GetUserInfoById(ctx.Param("id"))
	if err != nil {
		responseError(ctx, err, "获取用户失败")
		return
	}
	info.RemoveSensitive()
	responseSuccess(ctx, info, "")
}

//
// ApiGetUserByAccount
// @Description:
// @param ctx
//
func ApiGetUserByAccount(ctx *gin.Context) {
	user, err := UserServiceInstance().GetUserByAccount(ctx.Param("account"))
	if err != nil {
		responseError(ctx, err, "获取用户失败")
		return
	}
	//去除敏感信息
	user.UserInfo.RemoveSensitive()
	responseSuccess(ctx, user, "")
}

//
// ApiGetUserById
// @Description:
// @param ctx
//
func ApiGetUserById(ctx *gin.Context) {
	user, err := UserServiceInstance().GetUserById(ctx.Param("id"))
	if err != nil {
		responseError(ctx, err, "获取用户失败")
		return
	}
	//去除敏感信息
	user.Password = ""
	responseSuccess(ctx, user, "")
}

//
// ApiListUser
// @Description:
// @param ctx
//
func ApiListUser(ctx *gin.Context) {
	users, err := UserServiceInstance().ListUsers()
	if err != nil {
		responseError(ctx, err, "获取用户列表失败")
		return
	}
	responseSuccess(ctx, users, "")
}

//
// ApiListUserInfo
// @Description:
// @param ctx
//
func ApiListUserInfo(ctx *gin.Context) {
	users, err := UserServiceInstance().ListUserInfo()
	if err != nil {
		responseError(ctx, err, "获取用户列表失败")
		return
	}
	responseSuccess(ctx, users, "")
}

//
// ApiCreateUser
// @Description:
//
func ApiCreateUser(ctx *gin.Context) {
	info := model.UserInfo{}
	err := ctx.BindJSON(&info)
	if err != nil {
		response400(ctx)
		return
	}
	info.FormatCreate()
	err = UserServiceInstance().CreateUser(&info)
	if err != nil {
		responseError(ctx, err, "创建失败")
		return
	}
	responseSuccess(ctx, info, "创建["+info.Account+"]成功")
}

//
// ApiUpdateUser
// @Description:
// @param ctx
//
func ApiUpdateUser(ctx *gin.Context) {
	info := model.UserInfo{}
	err := ctx.BindJSON(&info)
	if err != nil {
		response400(ctx)
		return
	}
	if info.Id == "" {
		response400(ctx)
		return
	}
	err = UserServiceInstance().UpdateUser(&info)
	if err != nil {
		responseError(ctx, err, "修改失败")
		return
	}
	//去除敏感信息
	info.RemoveSensitive()
	responseSuccess(ctx, info, "修改["+info.Account+"]成功")
}

//
// ApiSetUserDisable
// @Description:
// @param ctx
//
func ApiSetUserDisable(ctx *gin.Context) {
	info := model.UserInfo{}
	err := ctx.BindJSON(&info)
	if err != nil {
		response400(ctx)
		return
	}
	if info.Id == "" {
		response400(ctx)
		return
	}
	err = UserServiceInstance().SetUserDisable(&info)
	if err != nil {
		responseError(ctx, err, "设置失败")
		return
	}
	msg := "禁用[" + info.Account + "]成功"
	if info.Disabled == 0 {
		msg = "解禁[" + info.Account + "]成功"
	}
	responseSuccess(ctx, "", msg)
}

//
// ApiDeleteUser
// @Description:
// @param ctx
//
func ApiDeleteUser(ctx *gin.Context) {
	info := model.UserInfo{}
	err := ctx.BindJSON(&info)
	if info.Id == "" {
		response400(ctx)
		return
	}
	if err != nil {
		response400(ctx)
		return
	}
	err = UserServiceInstance().DeleteUser(&info)
	if err != nil {
		responseError(ctx, err, "删除失败")
		return
	}
	//去除敏感信息
	info.RemoveSensitive()
	responseSuccess(ctx, info, "删除["+info.Account+"]成功")
}

//
// ApiGetUserGeneral
// @Description:
// @param ctx
//
func ApiGetUserGeneral(ctx *gin.Context) {
	general, err := UserServiceInstance().GetGeneral()
	if err != nil {
		responseError(ctx, err, "获取用户概况失败")
		return
	}
	responseSuccess(ctx, general, "")
}

//
// ApiGenerateUserConfig
// @Description:
// @param ctx
//
func ApiGenerateUserConfig(ctx *gin.Context) {
	pushedConfig, err := UserServiceInstance().GenerateUserConfig(ctx.Param("id"))
	if err != nil {
		responseError(ctx, err, "获取配置失败")
		return
	}
	responseSuccess(ctx, pushedConfig, "")
}
