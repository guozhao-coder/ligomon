package host

import (
	"github.com/cihub/seelog"
	"github.com/patrickmn/go-cache"
	"ligomonitor/utils"
	"strconv"
	"sync"
	"syscall"
	"time"
)

const (
	//Cache expiration time
	CACHETIMEOUT = time.Minute * 10
	//Clear time, here should be less than the expiration time, less than means not clear
	CACHECLEANUP = time.Second
)

var cacheAlarm *cache.Cache
var cacheOnce sync.Once

func alarmProcFunc(respFunc func(int), pid int) {
	cacheOnce.Do(func() {
		cacheAlarm = cache.New(CACHETIMEOUT, CACHECLEANUP)
	})
	//Determine if there is the event, add it if it is not, and trigger the function,
	// if it is, ignore it to prevent it from triggering every time it is traversed
	if _, ifExist := cacheAlarm.Get(strconv.Itoa(pid)); !ifExist {
		cacheAlarm.Set(strconv.Itoa(pid), respFunc, cache.DefaultExpiration)
		respFunc(pid)
	}
}

//kill the target process
func KillProcFunc(pid int) {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error(err)
		}
	}()
	err := syscall.Kill(pid, syscall.SIGKILL)
	if err != nil {
		err = utils.ErrJoint("kill process err :", err)
		seelog.Error(err.Error())
		return
	}
	return
}

//send mail to admin user
func MailAlertFunc(pid int) {

}
