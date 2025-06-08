package errorutils

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// blocking! wait for wg to finish with debug logging every 5 seconds. If maxCycles is -1, will wait indefinitely otherwise will terminate after maxCycles. Optional identifier for logging.
func MonitorWaitGroup(wg *sync.WaitGroup, maxCycles int, wgCompleted chan struct{}, wgName, id string) {
	logrus.Debugf("WG:%s waiting for wg", id)
	idStr := ""
	if id != "" {
		idStr = " " + id
	}
	wgNameStr := ""
	if wgName != "" {
		wgNameStr = " " + wgName
	}
	var cycles int
	done := make(chan bool)
	go func() {
		wg.Wait()
		close(done)
	}()
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
waitg:
	for {
		select {
		case <-ticker.C:
			logrus.Debugf("WG:%s waiting for wg%s", idStr, wgNameStr)
			if maxCycles != -1 && cycles > maxCycles {
				break waitg
			}
		case <-done:
			break waitg
		}
		cycles++
	}
	logrus.Debugf("WG:waitgroup%sfinished%s", wgNameStr, idStr)
	if wgCompleted != nil {
		close(wgCompleted)
	}
}

// tell the world that the block of code is still running
func ActiveBarker(name string, details string, finish_ch chan struct{}) {
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			logrus.Debugf("BARKER:Running %s %s", name, details)
		case <-finish_ch:
			logrus.Debugf("BARKER:execution ended for %s with %s", name, details)
			ticker.Stop()
			return
		}
	}
}
