package model

import "sync"

type Process struct {
	Pid           int      `json:"pid"`
	PPid          int      `json:"pPid"`
	Name          string   `json:"name"`
	Tgid          int      `json:"tgid"`
	State         string   `json:"state"`
	Uid           int      `json:"uid"`
	Gid           int      `json:"gid"`
	Threads       int      `json:"threads"`
	VmPeak        int64    `json:"vmPeak"`        //虚拟内存峰值
	VmSize        int64    `json:"vmSize"`        //当前虚拟内存使用
	VmLck         int64    `json:"vmLck"`         //进程已经锁住的物理内存的大小.锁住的物理内存不能交换到硬盘
	VmPin         int64    `json:"vmPin"`         //不可被移动的内存大小
	VmHWM         int64    `json:"vmHWM"`         //物理内存峰值
	VmRss         int64    `json:"vmRss"`         //当前物理内存
	VmData        int64    `json:"vmData"`        //进程数据段大小
	VmStk         int64    `json:"vmStk"`         //进程堆栈段大小
	VmExe         int64    `json:"vmExe"`         //进程代码段大小
	VmLib         int64    `json:"vmLib"`         //进程Lib库大小
	VmPTE         int64    `json:"vmPTE"`         //进程页表大小
	VmSwap        int64    `json:"vmSwap"`        //进程进入swap空间大小
	VoluntaryCS   int64    `json:"voluntaryCS"`   //进程主动切换次数
	NoVoluntaryCS int64    `json:"noVoluntaryCS"` //进程被动切换次数
	CPUUsage      float32  `json:"cpuUsage"`      //cpu使用率
	CPUUsed       int64    `json:"cpuUsed"`       //cpu使用情况
	Time          int      `json:"time"`
	IFAlarm       bool     `json:"ifAlarm"`                //是否报警
	AlarmMessage  AlarmMsg `json:"alarmMessage,omitempty"` //报警信息
}

type ProcessMsg struct {
	Procs []Process `json:"procs"`
	Err   error     `json:"err"`
}

type AlarmMsg struct {
	CPUMsg string `json:"cpuMsg"`
	VMMsg  string `json:"vmMsg"`
}

type AlarmLimit struct {
	VMLimit  int64        //内存最大限制
	CPULimit float32      //cpu最大限制
	Operate  AlarmOperate //告警后的操作
}

type AlarmOperate struct {
	Fnc func(int) error
}

type AlarmLimitData struct {
	Lck           sync.RWMutex
	AlarmLimitMap map[int]*AlarmLimit
}
