package svc

import (
	"github.com/cihub/seelog"
	"ligomonitor/pkg/model"
	"ligomonitor/pkg/service/dbsync"
	"ligomonitor/pkg/service/host"
	"sort"
)

func GetCurrentProcInfoSvc(pid int) ([]model.Process, error) {
	processes, err := host.GetProcessTotalInfo(pid, true)
	if err != nil {
		return []model.Process{}, err
	}
	if len(processes) > 1 {
		sort.Sort(processWapper{proc: processes, by: func(p, q *model.Process) bool {
			return p.CPUUsage > q.CPUUsage
		}})
	}
	return processes, nil
}

func GetProcInfoStreamSvc(msgchan chan model.ProcessMsg) {
	defer func() {
		//恢复程序
		if err := recover(); err != nil {
			seelog.Info(err)
		}
	}()
	processes, err := host.GetProcessTotalInfo(0, true)
	if err != nil {
		msgchan <- model.ProcessMsg{
			Procs: nil,
			Err:   err,
		}
		return
	}
	//长度大于1就排序
	if len(processes) > 1 {
		sort.Sort(processWapper{proc: processes, by: func(p, q *model.Process) bool {
			return p.CPUUsage > q.CPUUsage
		}})
	}
	msgchan <- model.ProcessMsg{
		Procs: processes,
		Err:   nil,
	}
	return
}

func KillProcSvc(pid int) error {
	return host.KillProcess(pid)
}

func GetProcRectDataSvc(pid int) ([]model.Process, error) {
	cli := dbsync.NewDBCli()
	return cli.GetProcessRectData(pid)
}

type processWapper struct {
	proc []model.Process
	by   func(p, q *model.Process) bool
}

func (pw processWapper) Len() int {
	return len(pw.proc)
}

func (pw processWapper) Swap(i, j int) {
	pw.proc[i], pw.proc[j] = pw.proc[j], pw.proc[i]
}

func (pw processWapper) Less(i, j int) bool {
	return pw.by(&pw.proc[i], &pw.proc[j])
}
