package administration

import "encoding/json"

//
// FetchResourceInfo
// @Description:
//
type FetchResourceInfo struct {
	Name   string            `json:"name"`
	Params map[string]string `json:"params"`
}

//
// FetchResourceResult
// @Description:
//
type FetchResourceResult struct {
	Data  interface{}
	Error string
}

//
// Byte
// @Description:
// @receiver res
// @return []byte
//
func (res FetchResourceResult) Byte() []byte {
	marshal, err := json.Marshal(res)
	if err != nil {
		return nil
	}
	return marshal
}

//
// Fetch
// @Description:
// @receiver i
// @return interface{}
//
func (i *FetchResourceInfo) Fetch() FetchResourceResult {
	switch i.Name {
	case "/api/v1/user/status/:id":
		return FetchResourceResult{
			UserServiceInstance().statusService.GetStatus(i.Params["id"]),
			"",
		}
	case "/api/v1/user/list":
		users, err := UserServiceInstance().ListUsers()
		return FetchResourceResult{
			users,
			err.Error(),
		}
	case "/api/v1/server/flow":
		return FetchResourceResult{
			ServerServiceInstance().GetFlowStatus(),
			"",
		}
	case "/api/v1/server/ippool":
		return FetchResourceResult{
			ServerServiceInstance().GetIPPoolGeneral(),
			"",
		}
	case "/api/v1/server/monitor":
		return FetchResourceResult{
			ServerServiceInstance().monitorService.GetSystemData(),
			"",
		}
	}
	return FetchResourceResult{
		Data:  nil,
		Error: "resource not found",
	}
}

//
// FetchResources
// @Description:
// @param resources
//
func FetchResources(resources map[string]FetchResourceInfo) map[string]interface{} {
	result := make(map[string]interface{})
	for cusKey := range resources {
		res := resources[cusKey]
		result[cusKey] = res.Fetch()
	}
	return result
}
