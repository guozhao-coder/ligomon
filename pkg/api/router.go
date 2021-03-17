package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"ligomonitor/cmd/app"
	"ligomonitor/pkg/cons"
	"ligomonitor/pkg/model"
	"ligomonitor/pkg/service/ctl"
	"net/http"
	"os"
	"strconv"
)

func StartRouter() error {
	//start pprof
	go startPprof()
	//start http serve
	r := Router(app.LigoConf)
	s := &http.Server{
		Addr:    ":" + strconv.Itoa(app.LigoConf.HTTPPort),
		Handler: r,
	}
	//start listen
	if err := s.ListenAndServe(); err != nil {
		fmt.Println("http listenAndServer error : ", err.Error())
		return err
	}
	return nil
}

func startPprof() {
	if err := http.ListenAndServe(":"+strconv.Itoa(app.LigoConf.PprofPort), nil); err != nil {
		fmt.Println("pprof listenAndServer error : ", err.Error())
		os.Exit(cons.PPROFSERVEERR)
	}
}

func Router(conf *model.LiGoMoniConf) *gin.Engine {
	r := gin.Default()
	setGinLog(conf.LogPath, r)
	groupMonitor := r.Group("/monitor")
	{
		groupProcess := groupMonitor.Group("/process")
		{
			//websocket,push the process info to client
			groupProcess.GET("/stream", ctl.GetProcInfoStreamCtl)
			//return the process infomation snapshot
			groupProcess.GET("/snapshot/:pid", ctl.GetCurrentProcInfoCtl)
			//kill the target process
			groupProcess.GET("/kill/:pid", ctl.KillProcCtl)
			//get the process recently
			groupProcess.GET("/status/:pid", ctl.GetProcRectDataCtl)
			//websocket,get the process syscall
			groupProcess.GET("/ptrace/:pid", ctl.GetProcSyscallStreamCtl)
			//register alarm event
			groupProcess.POST("/alarm/register", ctl.RegisterAlarmEventCtl)
		}
		groupHost := groupMonitor.Group("/host")
		{
			//get the host resource snapshot
			groupHost.GET("/snapshot", ctl.GetHostResourceCtl)
			//websocket,push the host resource status
			groupHost.GET("/stream", ctl.GetHostResourceStreamCtl)

		}
		/*
			groupDocker := groupMonitor.Group("/docker")
			{
				groupDocker.GET("/stream")
			}

		*/
	}

	return r
}

func setGinLog(path string, engine *gin.Engine) {
	logfd, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("open gin log file error, check the logPath :", err.Error())
		os.Exit(cons.LOGPATHERR)
	}
	gin.DefaultWriter = io.MultiWriter(logfd, os.Stdout)
	engine.Use(gin.Recovery())
}
