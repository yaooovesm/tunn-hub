package administration

import (
	"errors"
	"github.com/gin-gonic/gin"
	"tunn-hub/administration/model"
)

//
// ApiGetConfigById
// @Description:
// @param ctx
//
func ApiGetConfigById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		response400(ctx)
		return
	}
	cfg, err := UserServiceInstance().configService.GetById(id)
	if err != nil {
		responseError(ctx, err, "")
		return
	}
	responseSuccess(ctx, cfg, "")
}

//
// ApiListConfig
// @Description:
// @param ctx
//
func ApiListConfig(ctx *gin.Context) {
	list, err := UserServiceInstance().configService.List()
	if err != nil {
		responseError(ctx, err, "")
		return
	}
	responseSuccess(ctx, list, "")
}

//
// ApiDeleteConfigById
// @Description:
// @param ctx
//
func ApiDeleteConfigById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		response400(ctx)
		return
	}
	err := UserServiceInstance().configService.DeleteById(id)
	if err != nil {
		responseError(ctx, err, "删除设置失败")
		return
	}
	responseSuccess(ctx, "", "删除["+id+"]成功")
}

//
// ApiUpdateConfigById
// @Description:
// @param ctx
//
func ApiUpdateConfigById(ctx *gin.Context) {
	cfg := model.ClientConfig{}
	err := ctx.BindJSON(&cfg)
	if err != nil || cfg.Id == "" {
		response400(ctx)
		return
	}
	cfg, err = UserServiceInstance().configService.UpdateById(cfg, true)
	if err != nil {
		responseError(ctx, err, "更新设置失败")
		return
	}
	responseSuccess(ctx, cfg, "更新["+cfg.Id+"]成功")
}

//
// ApiResetConfigById
// @Description:
// @param ctx
//
func ApiResetConfigById(ctx *gin.Context) {
	cfg := model.ClientConfig{
		Id: ctx.Param("id"),
	}
	err := UserServiceInstance().configService.ResetById(cfg)
	if err != nil {
		responseError(ctx, err, "重置失败")
		return
	}
	responseSuccess(ctx, cfg, "重置成功")
}

//
// ApiCreateConfig
// @Description:
// @param ctx
//
func ApiCreateConfig(ctx *gin.Context) {
	cfg := model.ClientConfig{}
	err := ctx.BindJSON(&cfg)
	if err != nil {
		response400(ctx)
		return
	}
	create, err := UserServiceInstance().configService.Create(cfg)
	if err != nil {
		responseError(ctx, err, "")
		return
	}
	responseSuccess(ctx, create, "创建["+create.Id+"]成功")
}

//
// ApiAvailableExports
// @Description:
// @param ctx
//
func ApiAvailableExports(ctx *gin.Context) {
	account := ""
	if acc, exists := ctx.Get("account"); exists && acc != nil {
		account = acc.(string)
	} else {
		responseError(ctx, errors.New("no account"), "获取可导入网络列表失败")
		return
	}
	routes, err := UserServiceInstance().configService.AvailableExports(account)
	if err != nil {
		responseError(ctx, err, "获取可导入网络列表失败")
		return
	}
	responseSuccess(ctx, routes, "")
}

//
// ApiAvailableExportsByAccount
// @Description:
// @param ctx
//
func ApiAvailableExportsByAccount(ctx *gin.Context) {
	account := ctx.Param("account")
	if account == "" {
		responseError(ctx, errors.New("no account"), "获取可导入网络列表失败")
		return
	}
	routes, err := UserServiceInstance().configService.AvailableExports(account)
	if err != nil {
		responseError(ctx, err, "获取可导入网络列表失败")
		return
	}
	responseSuccess(ctx, routes, "")
}
