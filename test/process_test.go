package test

import (
	"fmt"
	"ligomonitor/pkg/service/host"
	"testing"
)

func TestProcessInfo(t *testing.T) {
	//for{
	processes, e := host.GetProcessParam(1)
	if e != nil {
		fmt.Println(e)
		return
	}
	for i := 0; i < len(processes); i++ {
		fmt.Println(processes[i].Pid, "   ", processes[i].Name, "   ", processes[i].CPUUsage)
		//	}
	}
}

func TestCpuUsage(t *testing.T) {
	cpu, e := host.GetCPUUsage()
	if e != nil {
		fmt.Println(e)
		return
	}
	fmt.Println(cpu)
}

func TestGetTotalCPUTime(t *testing.T) {
	//cpuTime, e := host.GetTotalCPUTime()
	//if e != nil{
	//	fmt.Println(e)
	//	return
	//}
	//fmt.Println(“”“”"cpu总时长",cpuTime)

	ptime, e := host.GetProcessCPUTime(7742)
	if e != nil {
		fmt.Println(e)
		return
	}
	fmt.Println("进程cpu时长", ptime)

}

func TestMemory(t *testing.T) {
	swap, e := host.GetMemAndSwapUsed()
	if e != nil {
		fmt.Println(e)
		return
	}
	fmt.Println(swap.Mem, swap.Swa)
}

func TestGetCpuUsage(t *testing.T) {
	var a int
	var b int
	var c float32

	a = 10
	b = 3
	c = float32(a) / float32(b)
	fmt.Println(c)
}

func BenchmarkProc(b *testing.B) {
	//b.ResetTimer()
	for i := 0; i < b.N; i++ {
		host.GetProcessTotalInfo(0)
	}
}
