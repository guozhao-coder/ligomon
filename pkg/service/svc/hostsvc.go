package svc

import (
	"github.com/cihub/seelog"
	"ligomonitor/pkg/model"
	"ligomonitor/pkg/service/host"
)

func GetHostResourceSvc() (*model.HostResourceUsed, error) {
	cpuUsage, err := host.GetCPUUsage()
	if err != nil {
		seelog.Error("get cpuusage error :", err.Error())
		return nil, err
	}
	memUsage, err := host.GetMemAndSwapUsed()
	if err != nil {
		seelog.Error("get memusage error :", err.Error())
		return nil, err
	}
	return &model.HostResourceUsed{
		CPUUsed: cpuUsage,
		MemUsed: memUsage,
	}, nil
}

//
func GetHostResourceStreamSvc(msgchan chan model.HostResourceMsg) {
	defer func() {
		//恢复程序
		if err := recover(); err != nil {
			seelog.Info(err)
		}
	}()
	cpuUsage, err := host.GetCPUUsage()
	if err != nil {
		msgchan <- model.HostResourceMsg{
			Resource: model.HostResourceUsed{},
			Err:      err,
		}
		return
	}
	memUsage, err := host.GetMemAndSwapUsed()
	if err != nil {
		msgchan <- model.HostResourceMsg{
			Resource: model.HostResourceUsed{},
			Err:      err,
		}
		return
	}
	msgchan <- model.HostResourceMsg{
		Resource: model.HostResourceUsed{
			CPUUsed: cpuUsage,
			MemUsed: memUsage,
		},
		Err: nil,
	}
	return
}
