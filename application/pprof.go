package application

import (
	log "github.com/cihub/seelog"
	"net/http"
	_ "net/http/pprof"
	"strconv"
)

//
// PProf
// @Description: start pprof
// @receiver app
//
func (app *Application) PProf() {
	if app.Config.Global.Pprof == 0 {
		return
	}
	if app.Config.Global.Pprof > 0 && app.Config.Global.Pprof <= 65535 {
		go func() {
			log.Info("pprof serve on 0.0.0.0:", app.Config.Global.Pprof)
			if err := http.ListenAndServe("0.0.0.0:"+strconv.Itoa(app.Config.Global.Pprof), nil); err != nil {
				_ = log.Warn("pprof start failed : ", err)
			}
		}()
	}
}
