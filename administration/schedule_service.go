package administration

import (
	log "github.com/cihub/seelog"
	"tunn-hub/schedule"
)

//
// scheduleService
// @Description:
//
type scheduleService struct {
	Schedules map[string]*schedule.Schedule
}

//
// newScheduleService
// @Description:
// @return *scheduleService
//
func newScheduleService() *scheduleService {
	return &scheduleService{Schedules: make(map[string]*schedule.Schedule)}
}

//
// Register
// @Description:
// @receiver s
// @param spec
// @param name
// @param f
// @return error
//
func (s *scheduleService) Register(spec string, name string, f func(), autorun bool) error {
	sc, err := schedule.NewSchedule(spec, name, f)
	if err != nil {
		return err
	}
	if autorun {
		sc.Start()
	}
	s.Schedules[name] = sc
	return nil
}

//
// StopByName
// @Description:
// @receiver s
// @param name
//
func (s *scheduleService) StopByName(name string) {
	if sc, ok := s.Schedules[name]; ok && sc != nil {
		log.Info("[schedule] ", name, " stopped...")
		sc.Stop()
	}
}

//
// StartByName
// @Description:
// @receiver s
// @param name
//
func (s *scheduleService) StartByName(name string) {
	if sc, ok := s.Schedules[name]; ok && sc != nil {
		log.Info("[schedule] ", name, " started...")
		sc.Start()
	}
}
