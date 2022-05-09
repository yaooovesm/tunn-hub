package config

/*RouteOption
对于客户端：
	RouteOptionExport : ①服务端本地添加系统路由表 ②服务器添加路由表(dst-->tunnel)将dst解析到tunnel
	RouteOptionImport : 客户端本地添加系统路由表(如/sbin/ip route add...)
对于服务端：
	RouteOptionExport : ①发送给客户端 ②客户端本地添加系统路由表
	RouteOptionImport : 服务端不执行import操作
*/
type RouteOption string

const (
	RouteOptionImport RouteOption = "import"
	RouteOptionExport RouteOption = "export"
)

//
// Route
// @Description:
//
type Route struct {
	Network string      `json:"network"`
	Option  RouteOption `json:"option"`
}
