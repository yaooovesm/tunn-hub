package application

import (
	"fmt"
	log "github.com/cihub/seelog"
	"os"
	"os/signal"
	"syscall"
)

//
// StopWatcher
// @Description:
// @param after
//
func StopWatcher(after func()) {
	go func() {
		osc := make(chan os.Signal, 1)
		signal.Notify(osc, syscall.SIGTERM, syscall.SIGINT)
		s := <-osc
		fmt.Println()
		log.Info("receive stop signal : ", s)
		after()
	}()
}
