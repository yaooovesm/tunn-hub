package config

//
// Schedule
// @Description:
//
type Schedule struct {
	DailyFlowRecord          bool   `json:"daily_flow_record"`           //每日记录流量
	WeeklyFlowRecord         bool   `json:"weekly_flow_record"`          //每周记录流量
	MonthlyFlowRecord        bool   `json:"monthly_flow_record"`         //每月记录流量
	MonthlyFlowReset         bool   `json:"monthly_flow_reset"`          //每月重置流量计数
	HourlyNetworkTrafficDump string `json:"hourly_network_traffic_dump"` //每小时保存流量速率记录路径
}
