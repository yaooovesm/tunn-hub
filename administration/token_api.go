package administration

import (
	"github.com/gin-gonic/gin"
	"tunn-hub/administration/model"
)

//
// ApiCreateToken
// @Description:
// @param ctx
//
func ApiCreateToken(ctx *gin.Context) {
	info := model.UserInfo{}
	err := ctx.BindJSON(&info)
	if err != nil {
		response400(ctx)
		return
	}
	remote, _ := ctx.RemoteIP()
	token, err := TokenServiceInstance().CreateAndCheck(remote, info.Account, info.Password)
	if err != nil {
		responseError(ctx, err, "获取token失败")
		return
	}
	responseSuccess(ctx, token, "")
}

//
// ApiCheckToken
// @Description:
// @param ctx
//
func ApiCheckToken(ctx *gin.Context) {
	info := model.UserInfo{}
	err := ctx.BindJSON(&info)
	if err != nil {
		response400(ctx)
		return
	}
	code := ctx.GetHeader("token")
	account := info.Account
	if code == "" || account == "" {
		response400(ctx)
		return
	}
	err = TokenServiceInstance().IsExist(account, code)
	if err != nil {
		responseError(ctx, err, "验证失败")
		return
	}
	responseSuccess(ctx, "", "验证成功")
}
