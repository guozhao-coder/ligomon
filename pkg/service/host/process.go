/*
提供宿主机数据的pkg
*/

package host

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/cihub/seelog"
	"io"
	"io/ioutil"
	"ligomonitor/pkg/model"
	"ligomonitor/utils"
	"os"
	"strconv"
	"sync"
	"syscall"
	"time"
)

type ifProcess bool
type withCPU bool

var AlarmLimitMap = model.AlarmLimitData{
	Lck:           sync.RWMutex{},
	AlarmLimitMap: make(map[int]*model.AlarmLimit),
}

func GetProcessTotalInfo(pid int, withCpu withCPU) ([]model.Process, error) {
	if !withCpu {
		return GetProcessParam(pid)
	}
	//snapshot1
	//cpu snapshot 1
	cpuSS1, err := GetTotalCPUTime()
	if err != nil {
		err = utils.ErrJoint("get cpu time err :", err)
		seelog.Error(err.Error())
		return nil, err
	}
	//process cpu snapshot
	procSS1, err := GetProcessParam(pid)
	if err != nil {
		err = utils.ErrJoint("get process param err :", err)
		seelog.Error(err.Error())
		return nil, err
	}
	procMap1 := make(map[int]model.Process)
	for i := 0; i < len(procSS1); i++ {
		procMap1[procSS1[i].Pid] = procSS1[i]
	}

	time.Sleep(time.Second * 2)

	//snapshot 2
	//cpu snapshot 2
	cpuSS2, err := GetTotalCPUTime()
	if err != nil {
		err = utils.ErrJoint("get cpu time err :", err)
		seelog.Error(err.Error())
		return nil, err
	}
	//process snapshot 2
	procSS2, err := GetProcessParam(pid)
	if err != nil {
		err = utils.ErrJoint("get process time err :", err)
		seelog.Error(err.Error())
		return nil, err
	}
	procMap2 := make(map[int]model.Process)
	for i := 0; i < len(procSS2); i++ {
		procMap2[procSS2[i].Pid] = procSS2[i]
	}

	//the array for return
	procs := []model.Process{}

	//calculate the cpu usage
	var cpuusg float32
	for m2k, m2v := range procMap2 {
		if m1v, ok := procMap1[m2k]; ok {
			cpuusg = float32(m2v.CPUUsed-m1v.CPUUsed) / float32(cpuSS2-cpuSS1) * float32(GetCPUNum())
			m2v.CPUUsage = cpuusg
			checkResourceUpperLimit(&m2v)
			procs = append(procs, m2v)
		}
	}
	return procs, nil

}

//provide the process info, if pid is 0, return the total process info
func GetProcessParam(pid int) ([]model.Process, error) {
	processes := []model.Process{}
	if pid == 0 {
		//iterator the /proc
		folders, err := ioutil.ReadDir("/proc")
		if err != nil {
			seelog.Error("open dir err : ", err.Error())
			return nil, err
		}

		for _, f := range folders {
			pid, err := strconv.Atoi(f.Name())
			if f.IsDir() && err == nil {
				processI, isProcess, err := getProcessInfo(pid)
				if err != nil {
					seelog.Error(err.Error())
					return nil, err
				}
				if isProcess {
					//如果出错，则返回孔结构体默认值0，就append
					if processI.Pid != 0 {
						processes = append(processes, processI)
					}
				}
			}
		}
		return processes, nil
	}
	//提供某个进程信息
	process, isProcess, err := getProcessInfo(pid)
	if err != nil {
		seelog.Error(err.Error())
		return nil, err
	}
	if isProcess {
		//如果不是默认值0，就append
		if process.Pid != 0 {
			processes = append(processes, process)
			return processes, nil
		}
		return []model.Process{}, nil
	}
	return []model.Process{}, errors.New("this Pid's process is a Thread !")
}

//Provide the information of a process,
//if it is a thread, return false,
//If there is an error, an empty struct is returned
func getProcessInfo(pid int) (model.Process, ifProcess, error) {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error(err)
		}
	}()

	file, err := os.Open("/proc/" + strconv.Itoa(pid) + "/status")
	defer file.Close()
	if err != nil {
		seelog.Error(err.Error())
		return model.Process{}, true, nil
	}
	procMap := make(map[string]string)
	scanProcFile := bufio.NewScanner(file)
	var sk, sv string
	for scanProcFile.Scan() {
		fmt.Sscan(scanProcFile.Text(), &sk, &sv)
		procMap[sk] = sv
	}
	var vmPeak1, vmSize1, vmLck1, vmPin1, vmHWM1, vmRSS1, vmData1, vmStk1, vmExe1, vmLib1, vmPTE1, vmSwap1 int64

	pName, err := getPidName(pid)
	if err != nil {
		return model.Process{}, true, nil
	}

	pid1, err := strconv.Atoi(procMap["Pid:"])
	if err != nil {
		seelog.Error(err.Error())
		return model.Process{}, true, nil
	}
	pPid1, err := strconv.Atoi(procMap["PPid:"])
	if err != nil {
		seelog.Error(err.Error())
		return model.Process{}, true, nil
	}
	tgid1, err := strconv.Atoi(procMap["Tgid:"])
	if err != nil {
		seelog.Error(err.Error())
		return model.Process{}, true, nil
	}

	//如果相等，说明是线程,需要返回false
	if pid1 != tgid1 {
		return model.Process{}, false, nil
	}

	uid1, err := strconv.Atoi(procMap["Uid:"])
	if err != nil {
		seelog.Error(err.Error())
		return model.Process{}, true, nil
	}
	gid1, err := strconv.Atoi(procMap["Gid:"])
	if err != nil {
		seelog.Error(err.Error())
		return model.Process{}, true, nil
	}
	threads1, err := strconv.Atoi(procMap["Threads:"])
	if err != nil {
		seelog.Error(err.Error())
		return model.Process{}, true, nil
	}

	if vmPeak, ok := procMap["VmPeak:"]; ok {
		vmPeak1, err = strconv.ParseInt(vmPeak, 10, 64)
		if err != nil {
			seelog.Error(err.Error())
			return model.Process{}, true, nil
		}
	}
	if vmSize, ok := procMap["VmSize:"]; ok {
		vmSize1, err = strconv.ParseInt(vmSize, 10, 64)
		if err != nil {
			seelog.Error(err.Error())
			return model.Process{}, true, nil
		}
	}

	if vmLck, ok := procMap["VmLck:"]; ok {
		vmLck1, err = strconv.ParseInt(vmLck, 10, 64)
		if err != nil {
			seelog.Error(err.Error())
			return model.Process{}, true, nil
		}
	}
	if vmPin, ok := procMap["VmPin:"]; ok {
		vmPin1, err = strconv.ParseInt(vmPin, 10, 64)
		if err != nil {
			seelog.Error(err.Error())
			return model.Process{}, true, nil
		}
	}

	if vmHWM, ok := procMap["VmHWM:"]; ok {
		vmHWM1, err = strconv.ParseInt(vmHWM, 10, 64)
		if err != nil {
			seelog.Error(err.Error())
			return model.Process{}, true, nil
		}
	}

	if vmRSS, ok := procMap["VmRSS:"]; ok {
		vmRSS1, err = strconv.ParseInt(vmRSS, 10, 64)
		if err != nil {
			seelog.Error(err.Error())
			return model.Process{}, true, nil
		}
	}
	if vmData, ok := procMap["VmData:"]; ok {
		vmData1, err = strconv.ParseInt(vmData, 10, 64)
		if err != nil {
			seelog.Error(err.Error())
			return model.Process{}, true, nil
		}
	}
	if vmStk, ok := procMap["VmStk:"]; ok {
		vmStk1, err = strconv.ParseInt(vmStk, 10, 64)
		if err != nil {
			seelog.Error(err.Error())
			return model.Process{}, true, nil
		}
	}
	if vmExe, ok := procMap["VmExe:"]; ok {
		vmExe1, err = strconv.ParseInt(vmExe, 10, 64)
		if err != nil {
			seelog.Error(err.Error())
			return model.Process{}, true, nil
		}
	}
	if vmLib, ok := procMap["VmLib:"]; ok {
		vmLib1, err = strconv.ParseInt(vmLib, 10, 64)
		//由于没有大于int64的类型，这里通过错误类型判断
		if err != nil {
			errStr := err.Error()[len(err.Error())-5 : len(err.Error())]
			if errStr != "range" {
				seelog.Error(err.Error())
				return model.Process{}, true, nil
			}
			vmLib1 = 0
		}
	}
	if vmPTE, ok := procMap["VmPTE:"]; ok {
		vmPTE1, err = strconv.ParseInt(vmPTE, 10, 64)
		if err != nil {
			seelog.Error(err.Error())
			return model.Process{}, true, nil
		}
	}

	if vmSwap, ok := procMap["VmSwap:"]; ok {
		vmSwap1, err = strconv.ParseInt(vmSwap, 10, 64)
		if err != nil {
			seelog.Error(err.Error())
			return model.Process{}, true, nil
		}
	}

	voluntaryCS1, err := strconv.ParseInt((procMap["voluntary_ctxt_switches:"]), 10, 64)
	if err != nil {
		seelog.Error(err.Error())
		return model.Process{}, true, nil
	}
	noVoluntaryCS1, err := strconv.ParseInt((procMap["nonvoluntary_ctxt_switches:"]), 10, 64)
	if err != nil {
		seelog.Error(err.Error())
		return model.Process{}, true, nil
	}
	var pCpuT1 int64
	if pCpuT1, err = GetProcessCPUTime(pid); err != nil {
		return model.Process{}, true, nil
	}

	return model.Process{
		Pid:           pid1,
		PPid:          pPid1,
		Name:          pName,
		Tgid:          tgid1,
		State:         procMap["State:"],
		Uid:           uid1,
		Gid:           gid1,
		Threads:       threads1,
		VmPeak:        vmPeak1,
		VmSize:        vmSize1,
		VmLck:         vmLck1,
		VmPin:         vmPin1,
		VmHWM:         vmHWM1,
		VmRss:         vmRSS1,
		VmData:        vmData1,
		VmStk:         vmStk1,
		VmExe:         vmExe1,
		VmLib:         vmLib1,
		VmPTE:         vmPTE1,
		VmSwap:        vmSwap1,
		VoluntaryCS:   voluntaryCS1,
		NoVoluntaryCS: noVoluntaryCS1,
		CPUUsed:       pCpuT1,
	}, true, nil
}

func getPidName(pid int) (string, error) {
	nameFile, err := os.Open("/proc/" + strconv.Itoa(pid) + "/cmdline")
	defer nameFile.Close()
	if err != nil {
		seelog.Error(err.Error())
		return "", err
	}
	nameBuf := make([]byte, 1024)
	n, err := nameFile.Read(nameBuf)
	if err != nil && err != io.EOF {
		seelog.Error(err.Error())
		return "", err
	}
	return string(nameBuf[:n]), nil
}

//kill the target process
func KillProcess(pid int) error {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error(err)
		}
	}()
	err := syscall.Kill(pid, syscall.SIGKILL)
	if err != nil {
		err = utils.ErrJoint("kill process err :", err)
		seelog.Error(err.Error())
		return err
	}
	return nil
}

//check the resource ,it will be alarmed if it is upper to the limit
func checkResourceUpperLimit(proc *model.Process) {
	AlarmLimitMap.Lck.RLock()
	defer AlarmLimitMap.Lck.RUnlock()
	if alarmLimitData, ok := AlarmLimitMap.AlarmLimitMap[proc.Pid]; ok {
		//compare
		timeNow := time.Now().Format("2006-01-02 15:04:05")
		var errS string
		if alarmLimitData.CPULimit < proc.CPUUsage {
			errS = fmt.Sprintf("%s : %d process's cpu usage is upper the limit , upper limit is %f , current is %f ,please check it !", timeNow, proc.Pid, alarmLimitData.CPULimit, proc.CPUUsage)
			seelog.Error(errS)
			proc.AlarmMessage.CPUMsg = errS
		}
		if alarmLimitData.VMLimit < proc.VmRss {
			errS = fmt.Sprintf("%s : %d process's memory is upper the limit,  upper limit is %d , current is %d  ,please check it !", timeNow, proc.Pid, alarmLimitData.VMLimit, proc.VmRss)
			seelog.Error(errS)
			proc.AlarmMessage.VMMsg = errS
		}
		if errS != "" {
			proc.IFAlarm = true
			//execute the register func
			alarmLimitData.Operate.Fnc(proc.Pid)
		}
	}
}
