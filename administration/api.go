package administration

import (
	"github.com/gin-gonic/gin"
)

//
// ServerAdminApi
// @Description:
//
type ServerAdminApi struct {
	r *gin.Engine
}

//
// NewServerAdminApi
// @Description:
// @param engine
// @return *ServerAdminApi
//
func NewServerAdminApi(r *gin.Engine) *ServerAdminApi {
	//注册中间件
	r.Use(AuthCheckMiddleWare())
	return &ServerAdminApi{r: r}
}

//
// Serv
// @Description:
// @receiver api
//
func (api *ServerAdminApi) Serv() {
	group := api.r.Group("/api")
	api.v1(group)
}

//
// v1
// @Description:
// @receiver api
// @param r
//
func (api *ServerAdminApi) v1(r *gin.RouterGroup) {
	group := r.Group("/v1")
	//user
	{
		userGroup := group.Group("/user")
		userGroup.GET("/info/:id", ApiGetUserInfoById)
		userGroup.GET("/status/:id", ApiGetUserStatusById)
		userGroup.GET("/:account", ApiGetUserByAccount)
		userGroup.GET("/id/:id", ApiGetUserById)
		userGroup.GET("/list", ApiListUser)
		userGroup.GET("/list/info", ApiListUserInfo)
		userGroup.PUT("/create", ApiCreateUser)
		userGroup.POST("/update/:id", ApiUpdateUser)
		userGroup.POST("/disable/:id", ApiSetUserDisable)
		userGroup.DELETE("/delete", ApiDeleteUser)
		userGroup.POST("/login", ApiUserLogin)
		userGroup.GET("/logout/:account", ApiUserLogout)
		userGroup.GET("/general", ApiGetUserGeneral)
		userGroup.GET("/config/:id", ApiGenerateUserConfig)
	}
	//token
	{
		tokenGroup := group.Group("/token")
		tokenGroup.POST("/", ApiCreateToken)
		tokenGroup.POST("/check", ApiCheckToken)
	}
	//server
	{
		serverGroup := group.Group("/server")
		serverGroup.GET("/flow", ApiGetServerFlowStatus)
		serverGroup.GET("/config", ApiGetServerConfigs)
		serverGroup.GET("/version", ApiGetServerVersion)
		serverGroup.GET("/disconnect/id/:id", ApiKickById)
		serverGroup.GET("/reconnect/id/:id", ApiReconnectById)
		serverGroup.GET("/ippool", ApiIPPoolGeneral)
		serverGroup.GET("/ippool/list", ApiIPPoolInfoList)
	}
	//config
	{
		cfgGroup := group.Group("/cfg")
		cfgGroup.GET("/:id", ApiGetConfigById)
		cfgGroup.GET("/list", ApiListConfig)
		cfgGroup.DELETE("/delete/:id", ApiDeleteConfigById)
		cfgGroup.POST("/update", ApiUpdateConfigById)
		cfgGroup.PUT("/create", ApiCreateConfig)
	}
}
