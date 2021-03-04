package host

import (
	"bufio"
	"fmt"
	"github.com/cihub/seelog"
	"ligomonitor/pkg/model"
	"ligomonitor/utils"
	"os"
)

func GetMemAndSwapUsed() (*model.MemAndSwap, error) {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error(err)
		}
	}()
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		err = utils.ErrJoint("open /proc/meminfo error :", err)
		seelog.Error(err.Error())
		return nil, err
	}
	var mk string
	var mv int32
	memInfoMap := make(map[string]int32)
	scanMemFile := bufio.NewScanner(file)
	for scanMemFile.Scan() {
		fmt.Sscan(scanMemFile.Text(), &mk, &mv)
		memInfoMap[mk] = mv
	}
	return &model.MemAndSwap{
		Mem: &model.Memory{
			Total:     memInfoMap["MemTotal:"],
			Free:      memInfoMap["MemFree:"],
			Used:      memInfoMap["MemTotal:"] - memInfoMap["MemFree:"],
			BuffCache: memInfoMap["Buffers:"] + memInfoMap["Cached:"],
		},
		Swa: &model.Swap{
			Total: memInfoMap["SwapTotal:"],
			Used:  memInfoMap["SwapTotal:"] - memInfoMap["SwapFree:"],
			Free:  memInfoMap["SwapFree:"],
		},
	}, nil

}
