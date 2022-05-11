package authentication

import (
	"net"
	"tunn-hub/config"
)

//
// getExportRoutes
// @Description:
// @return string
//
func getExportRoutes() []config.Route {
	var routes []config.Route
	for i := range config.Current.Routes {
		if config.Current.Routes[i].Option == config.RouteOptionExport {
			route := config.Current.Routes[i]
			route.Option = config.RouteOptionImport
			routes = append(routes, route)
		}
	}
	return routes
}

//
// GetRemoteAddr
// @Description:
// @param conn
// @return string
//
func GetRemoteAddr(conn net.Conn) string {
	remoteAddr := conn.RemoteAddr().String()
	//remoteAddr = remoteAddr[0:strings.Index(remoteAddr, ":")]
	return remoteAddr
}
