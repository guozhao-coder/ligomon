package utils

import (
	"fmt"
	logs "github.com/cihub/seelog"
	"ligomonitor/pkg/cons"
	"os"
)

func GetLogConfig(dir string) {
	fmt.Println("Read log file:", dir)
	log, err := logs.LoggerFromConfigAsFile(dir)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(cons.LOGREADERR)
	}
	logs.ReplaceLogger(log)
}
