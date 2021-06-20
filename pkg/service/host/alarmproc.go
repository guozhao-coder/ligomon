package host

import (
	"fmt"
	"github.com/cihub/seelog"
	"github.com/patrickmn/go-cache"
	"gopkg.in/gomail.v2"
	"ligomonitor/pkg/model"
	"ligomonitor/utils"
	"strconv"
	"sync"
	"syscall"
	"time"
)

const (
	//Cache expiration time
	CACHETIMEOUT = time.Minute * 20
	//Clear time, here should be less than the expiration time, less than means not clear
	CACHECLEANUP = time.Second
)

var ligoConf *model.LiGoMoniConf

var cacheAlarm *cache.Cache
var cacheOnce sync.Once

func alarmProcFunc(respFunc func(int, string), pid int, excp string) {
	cacheOnce.Do(func() {
		cacheAlarm = cache.New(CACHETIMEOUT, CACHECLEANUP)
	})
	//Determine if there is the event, add it if it is not, and trigger the function,
	// if it is, ignore it to prevent it from triggering every time it is traversed
	if _, ifExist := cacheAlarm.Get(strconv.Itoa(pid)); !ifExist {
		cacheAlarm.Set(strconv.Itoa(pid), respFunc, cache.DefaultExpiration)
		respFunc(pid, excp)
	}
}

//kill the target process
func KillProcFunc(pid int, excpt string) {
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
func MailAlertFunc(pid int, except string) {
	go func() {
		mailConn := map[string]string{
			"user": ligoConf.MailConf.User,
			"pass": ligoConf.MailConf.Pass,
			"host": ligoConf.MailConf.Host,
			"port": ligoConf.MailConf.Port,
		}
		port, _ := strconv.Atoi(mailConn["port"])
		m := gomail.NewMessage()
		m.SetHeader("From", "ligomon"+"<"+mailConn["user"]+">")
		mailTo := []string{ligoConf.MailConf.User}
		m.SetHeader("To", mailTo...)
		m.SetHeader("Subject", "ligomon告警通知")
		m.SetBody("text/html", fmt.Sprintf("您好，您的机器有进程告警，进程号为%d，告警消息为：%v", pid, except))
		d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])

		err := d.DialAndSend(m)
		if err != nil {
			seelog.Error(err.Error())
		}
	}()
}

//the global config in package host
func InitHostConf(conf *model.LiGoMoniConf) {
	ligoConf = conf
}
