package svc

import (
	"github.com/cihub/seelog"
	"ligomonitor/pkg/model"
	"ligomonitor/pkg/service/host"
)

func GetHostResourceSvc() (*model.HostResourceUsed, error) {
	host.GlobalHostInfoMap.Lck.RLock()
	cpuUsage := host.GlobalHostInfoMap.HostResource.CPUUsed
	memUsage := host.GlobalHostInfoMap.HostResource.MemUsed
	host.GlobalHostInfoMap.Lck.RUnlock()
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
	host.GlobalHostInfoMap.Lck.RLock()
	cpuUsage := host.GlobalHostInfoMap.HostResource.CPUUsed
	memUsage := host.GlobalHostInfoMap.HostResource.MemUsed
	host.GlobalHostInfoMap.Lck.RUnlock()
	msgchan <- model.HostResourceMsg{
		Resource: model.HostResourceUsed{
			CPUUsed: cpuUsage,
			MemUsed: memUsage,
		},
		Err: nil,
	}
	return
}
