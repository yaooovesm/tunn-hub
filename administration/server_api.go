package administration

import (
	"github.com/gin-gonic/gin"
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
