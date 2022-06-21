package administration

import (
	"github.com/gin-gonic/gin"
	"path"
	"time"
	"tunn-hub/config"
	"tunn-hub/security"
	"tunn-hub/version"
)

//
// ApiGetServerFlowStatus
// @Description:
// @param ctx
//
func ApiGetServerFlowStatus(ctx *gin.Context) {
	status := ServerServiceInstance().GetFlowStatus()
	responseSuccess(ctx, status, "")
}

//
// ApiGetServerConfigs
// @Description:
// @param ctx
//
func ApiGetServerConfigs(ctx *gin.Context) {
	cfg := ServerServiceInstance().GetServerConfigs()
	responseSuccess(ctx, cfg, "")
}

//
// ApiGetServerVersion
// @Description:
// @param ctx
//
func ApiGetServerVersion(ctx *gin.Context) {
	responseSuccess(ctx, map[string]interface{}{
		"version": version.Version,
		"develop": version.Develop,
	}, "")
}

//
// ApiKickById
// @Description:
// @param ctx
//
func ApiKickById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		response400(ctx)
		return
	}
	err := ServerServiceInstance().KickById(id)
	if err != nil {
		responseError(ctx, err, "")
		return
	}
	responseSuccess(ctx, "", "已断开与用户["+id+"]的连接")
}

//
// ApiReconnectById
// @Description:
// @param ctx
//
func ApiReconnectById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		response400(ctx)
		return
	}
	err := ServerServiceInstance().ReconnectById(id)
	if err != nil {
		responseError(ctx, err, "")
		return
	}
	responseSuccess(ctx, "", "已重置与用户["+id+"]的连接")
}

//
// ApiIPPoolGeneral
// @Description:
// @param ctx
//
func ApiIPPoolGeneral(ctx *gin.Context) {
	responseSuccess(ctx, ServerServiceInstance().GetIPPoolGeneral(), "")
}

//
// ApiIPPoolInfoList
// @Description:
// @param ctx
//
func ApiIPPoolInfoList(ctx *gin.Context) {
	responseSuccess(ctx, ServerServiceInstance().GetIPPoolAllocInfo(), "")
}

//
// ApiGetServerSystemData
// @Description:
// @param ctx
//
func ApiGetServerSystemData(ctx *gin.Context) {
	responseSuccess(ctx, ServerServiceInstance().monitorService.GetSystemData(), "")
}

//
// ApiCreateTLSCert
// @Description:
// @param ctx
//
func ApiCreateTLSCert(ctx *gin.Context) {
	req := struct {
		Overwrite bool     `json:"overwrite"`
		Addresses []string `json:"addresses"`
		Names     []string `json:"names"`
		Before    int64    `json:"before"`
	}{}
	err := ctx.BindJSON(&req)
	if err != nil {
		response400(ctx)
		return
	}
	dir := "./cert/"
	certification := security.NewTunnX509Certification(req.Addresses, req.Names, time.UnixMilli(req.Before))
	name, err := certification.CreateAndWriteWithRandomName(dir)
	if err != nil {
		responseError(ctx, err, "")
		return
	}
	if req.Overwrite {
		oldCert := config.Current.Security.CertPem
		oldKey := config.Current.Security.KeyPem
		config.Current.Security.CertPem = dir + name + ".cert"
		config.Current.Security.KeyPem = dir + name + ".key"
		err := config.Current.Storage()
		if err != nil {
			responseError(ctx, err, "")
			return
		}
		config.Current.Security.CertPem = oldCert
		config.Current.Security.KeyPem = oldKey
	}
	responseSuccess(ctx, name, "")
}

//
// ApiDownloadCurrentCert
// @Description:
// @param ctx
//
func ApiDownloadCurrentCert(ctx *gin.Context) {
	filename := path.Base(config.Current.Security.CertPem)
	fileContentDisposition := "attachment;filename=\"" + filename + "\""
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", fileContentDisposition)
	ctx.Header("FileName", filename)
	ctx.File(config.Current.Security.CertPem)
}

//
// ApiRecentHourTrafficData
// @Description:
// @param ctx
//
func ApiRecentHourTrafficData(ctx *gin.Context) {
	responseSuccess(ctx, ServerServiceInstance().GetRecentHourTrafficData(), "")
}

//
// ApiRecentDayTrafficData
// @Description:
// @param ctx
//
func ApiRecentDayTrafficData(ctx *gin.Context) {
	responseSuccess(ctx, ServerServiceInstance().GetRecentDayTrafficData(), "")
}
