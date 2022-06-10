package administration

//
// setupSchedules
// @Description:
// @param service
//
func setupSchedules(service *scheduleService) {
	//err := service.Register("@every 30s", "flow_recorder", func() {
	//	if ServerServiceInstance() == nil {
	//		return
	//	}
	//	rx := uint64(0)
	//	tx := uint64(0)
	//	if ServerServiceInstance().rxFlowCounter != nil {
	//		rx = ServerServiceInstance().rxFlowCounter.FlowSpeed
	//	}
	//	if ServerServiceInstance().txFlowCounter != nil {
	//		tx = ServerServiceInstance().txFlowCounter.FlowSpeed
	//	}
	//	fmt.Println()
	//	fmt.Println("-------------------------")
	//	fmt.Println("rx <-- ", rx)
	//	fmt.Println("tx <-- ", tx)
	//	fmt.Println("-------------------------")
	//	fmt.Println()
	//}, true)
	//if err != nil {
	//	_ = log.Warn("flow_recorder register failed : ", err)
	//	return
	//}
}

////
//// recordAndResetFlow
//// @Description:
////
//func recordAndResetFlow() {
//	infos, err := UserServiceInstance().ListUserInfo()
//	if err != nil {
//		return
//	}
//	for i := range infos {
//		info := model.UserInfo{
//			Id: infos[i].Id,
//		}
//		err := UserServiceInstance().ResetFlowCounter(&info)
//		if err != nil {
//			return
//		}
//	}
//}
