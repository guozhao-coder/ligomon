package host

import (
	"bufio"
	"fmt"
	"github.com/cihub/seelog"
	"ligomonitor/pkg/model"
	"ligomonitor/utils"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var cpuNum int
var cpuNumOnce sync.Once

func GetCPUUsage() (*model.CPU, error) {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error(err)
		}
	}()
	var us1, ni1, sy1, id1, io1, hi1, si1, st1 int32
	var us2, ni2, sy2, id2, io2, hi2, si2, st2 int32
	var t1, t2 int32

	file, err := os.Open("/proc/stat")
	defer file.Close()
	if err != nil {
		err = utils.ErrJoint("open /proc/stat err :", err)
		seelog.Error(err.Error())
		return nil, err
	}
	var extraStr string
	scanStatFile := bufio.NewScanner(file)
	fmt.Println(scanStatFile.Text())
	for scanStatFile.Scan() {
		if scanStatFile.Text() != "" {
			fmt.Sscan(scanStatFile.Text(), &extraStr, &us1, &ni1, &sy1, &id1, &io1, &hi1, &si1, &st1)
			break
		}
	}
	t1 = us1 + ni1 + sy1 + id1 + io1 + hi1 + si1 + st1

	time.Sleep(time.Second * 2)

	file2, err := os.Open("/proc/stat")
	defer file2.Close()
	if err != nil {
		err = utils.ErrJoint("open /proc/stat err :", err)
		seelog.Error(err.Error())
		return nil, err
	}
	scanStatFile2 := bufio.NewScanner(file2)
	for scanStatFile2.Scan() {
		if scanStatFile2.Text() != "" {
			fmt.Sscan(scanStatFile2.Text(), &extraStr, &us2, &ni2, &sy2, &id2, &io2, &hi2, &si2, &st2)
			break
		}
	}
	t2 = us2 + ni2 + sy2 + id2 + io2 + hi2 + si2 + st2

	return &model.CPU{
		UserUsage:   float32(us2-us1) / float32(t2-t1),
		NiceUsage:   float32(ni2-ni1) / float32(t2-t1),
		SystemUsage: float32(sy2-sy1) / float32(t2-t1),
		IdleUsage:   float32(id2-id1) / float32(t2-t1),
		IowaitUsage: float32(io2-io1) / float32(t2-t1),
		HIrqUsage:   float32(hi2-hi1) / float32(t2-t1),
		SIrqUsage:   float32(si2-si1) / float32(t2-t1),
	}, nil
}

//获取cpu总时长
func GetTotalCPUTime() (int32, error) {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error(err)
		}
	}()
	cpuf, err := os.Open("/proc/stat")
	defer cpuf.Close()
	if err != nil {
		err = utils.ErrJoint("open /proc/stat err :", err)
		seelog.Error(err.Error())
		return 0, err
	}
	fbuf := make([]byte, 5000)
	n, err := cpuf.Read(fbuf)
	if err != nil {
		err = utils.ErrJoint("read to buf err :", err)
		seelog.Error(err.Error())
		return 0, err
	}
	splitCpuTime := strings.Split(strings.Split(string(fbuf[:n]), "\n")[0], " ")
	var cpuTime int32
	for i := 2; i < len(splitCpuTime); i++ {
		cpuPTime, err := strconv.Atoi(splitCpuTime[i])
		if err != nil {
			err = utils.ErrJoint("strconv err :", err)
			seelog.Error(err.Error())
			return 0, err
		}
		cpuTime += int32(cpuPTime)
	}
	return cpuTime, nil
}

//获取某个进程的cpu时长
func GetProcessCPUTime(pid int) (int64, error) {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error(err)
		}
	}()
	pCpuf, err := os.Open("/proc/" + strconv.Itoa(pid) + "/stat")
	defer pCpuf.Close()
	if err != nil {
		seelog.Error(err.Error())
		return 0, err
	}
	pBuf := make([]byte, 5000)
	n, err := pCpuf.Read(pBuf)
	if err != nil {
		err = utils.ErrJoint("read to buf err :", err)
		seelog.Error(err.Error())
		return 0, err
	}
	pstat := strings.Split(string(pBuf[:n]), " ")
	var pCpuTime int64
	//第14到17个元素是cpu时间
	for i := 13; i < 17; i++ {
		pCpuPTime, err := strconv.ParseInt(pstat[i], 10, 64)
		if err != nil {
			err = utils.ErrJoint("strconv err :", err)
			seelog.Error(err.Error())
			return 0, err
		}
		pCpuTime += pCpuPTime

	}
	return pCpuTime, nil
}

func GetCPUNum() int {
	cpuNumOnce.Do(func() {
		cpuNum = runtime.NumCPU()
	})
	return cpuNum
}
