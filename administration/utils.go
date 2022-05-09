package administration

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"net/http"
	"strings"
)

//
// response400
// @Description:
// @param ctx
//
func response400(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, map[string]string{
		"msg": "bad request",
	})
}

//
// responseError
// @Description:
// @param ctx
// @param err
//
func responseError(ctx *gin.Context, err error, msg string) {
	ctx.JSON(http.StatusInternalServerError, map[string]string{
		"error": err.Error(),
		"msg":   msg,
	})
}

//
// responseError
// @Description:
// @param ctx
// @param err
//
func responseSuccess(ctx *gin.Context, data interface{}, msg string) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"error": "",
		"data":  data,
		"msg":   msg,
	})
}

//
// UUID
// @Description:
// @return string
//
func UUID() string {
	v4, err := uuid.NewV4()
	if err != nil {
		return ""
	}
	return strings.ReplaceAll(v4.String(), "-", "")
}
