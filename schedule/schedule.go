package schedule

import (
	log "github.com/cihub/seelog"
	"github.com/robfig/cron/v3"
)

//
// Schedule
// @Description:
//
type Schedule struct {
	Name    string     `json:"name"`
	Cron    *cron.Cron `json:"-"`
	Running bool       `json:"running"`
}

//
// NewSchedule
// @Description:
// @param spec
// @param name
// @param f
// @return *Schedule
// @return error
//
func NewSchedule(spec string, name string, f func()) (*Schedule, error) {
	c := cron.New()
	_, err := c.AddFunc(spec, f)
	if err != nil {
		return nil, err
	}
	return &Schedule{
		Name:    name,
		Cron:    c,
		Running: false,
	}, nil
}

//
// Start
// @Description:
// @receiver s
//
func (s *Schedule) Start() {
	if s.Running {
		return
	}
	s.Cron.Start()
	log.Info("[schedule:", s.Name, "] started...")
	s.Running = true
}

//
// Stop
// @Description:
// @receiver s
//
func (s *Schedule) Stop() {
	if !s.Running {
		return
	}
	s.Cron.Stop()
	log.Info("[schedule:", s.Name, "] stopped...")
	s.Running = false
}
