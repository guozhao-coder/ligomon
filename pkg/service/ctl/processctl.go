package ctl

import (
	"encoding/json"
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"ligomonitor/cmd/app"
	"ligomonitor/pkg/cons"
	"ligomonitor/pkg/model"
	"ligomonitor/pkg/service/host"
	"ligomonitor/pkg/service/svc"
	"net/http"
	"strconv"
	"time"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func GetCurrentProcInfoCtl(c *gin.Context) {
	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil {
		c.JSON(cons.HTTPHEADERCODE, model.NormalResponse{
			Code:    cons.URLPARAMERR,
			Message: "url analysis error,please check pid",
		})
		return
	}
	processes, err := svc.GetCurrentProcInfoSvc(pid)
	if err != nil {
		c.JSON(cons.HTTPHEADERCODE, model.NormalResponse{
			Code:    cons.SERVERERR,
			Message: err.Error(),
		})
		return
	}

	c.JSON(cons.HTTPHEADERCODE, model.NormalResponse{
		Code:    cons.SUCCESS,
		Message: "success",
		Data:    processes,
	})

}

func GetProcInfoStreamCtl(c *gin.Context) {
	wsconn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		seelog.Error("websocket conn error : ", err.Error())
		c.JSON(cons.SUCCESS, model.NormalResponse{
			Code:    cons.WSCONNERR,
			Message: err.Error(),
		})
		return
	}
	defer wsconn.Close()
	topTicker := time.NewTicker(time.Second * time.Duration(app.LigoConf.TopFlushTime))
	msgChan := make(chan model.ProcessMsg, 100)
	defer close(msgChan)
	go func() {
		for {
			procsMsg, ifClose := <-msgChan
			//首先判断下管道是否关闭，如果关闭了就退出协程，避免资源泄露
			if !ifClose {
				seelog.Info("the procsMsg channel is closed .")
				return
			}
			if procsMsg.Err != nil {
				wsconn.WriteJSON(model.NormalResponse{
					Code:    cons.SERVERERR,
					Message: procsMsg.Err.Error(),
				})
			}
			wsconn.WriteJSON(model.NormalResponse{
				Code:    cons.SUCCESS,
				Message: "success",
				Data:    procsMsg.Procs,
			})

		}

	}()
	for {
		//发送心跳探测客户端是否存在
		if err := wsconn.WriteJSON(model.NormalResponse{Code: cons.HEARTBEAT, Message: "heartbeat"}); err != nil {
			seelog.Info("client heartbeat failed :", err.Error())
			return
		}
		select {
		case <-topTicker.C:
			//启动一个goroutine去执行任务，获取完放入管道
			go svc.GetProcInfoStreamSvc(msgChan)
		}
	}
}

func KillProcCtl(c *gin.Context) {
	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil {
		c.JSON(cons.HTTPHEADERCODE, model.NormalResponse{
			Code:    cons.URLPARAMERR,
			Message: "url analysis error,please check pid",
		})
		return
	}
	err = svc.KillProcSvc(pid)
	if err != nil {
		c.JSON(cons.HTTPHEADERCODE, model.NormalResponse{
			Code:    cons.SERVERERR,
			Message: err.Error(),
		})
		return
	}
	c.JSON(cons.HTTPHEADERCODE, model.NormalResponse{
		Code:    cons.SUCCESS,
		Message: "success",
	})
}

func GetProcRectDataCtl(c *gin.Context) {
	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil {
		c.JSON(cons.HTTPHEADERCODE, model.NormalResponse{
			Code:    cons.URLPARAMERR,
			Message: "url analysis error,please check pid",
		})
		return
	}
	processes, err := svc.GetProcRectDataSvc(pid)
	if err != nil {
		c.JSON(cons.HTTPHEADERCODE, model.NormalResponse{
			Code:    cons.SERVERERR,
			Message: err.Error(),
		})
		return
	}
	c.JSON(cons.HTTPHEADERCODE, model.NormalResponse{
		Code:    cons.SUCCESS,
		Message: "success",
		Data:    processes,
	})
}

//对于trace信息，生产者消费者模型，由消费者关闭管道
//todo bug:当被追踪进程退出时候，有些内存泄露,但是当socket断掉，就恢复了
func GetProcSyscallStreamCtl(c *gin.Context) {
	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil {
		c.JSON(cons.HTTPHEADERCODE, model.NormalResponse{
			Code:    cons.URLPARAMERR,
			Message: "url analysis error,please check pid",
		})
		return
	}
	wsconn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		seelog.Error("websocket conn error : ", err.Error())
		c.JSON(cons.SUCCESS, model.NormalResponse{
			Code:    cons.WSCONNERR,
			Message: err.Error(),
		})
		return
	}
	defer wsconn.Close()
	msgchan := make(chan model.PtraceMsg, 1000)
	defer close(msgchan)
	//开启消费者
	go func() {
		for {
			msg, ifClose := <-msgchan
			if !ifClose {
				seelog.Info("the ptraceMsg channel is closed .")
				return
			}
			if msg.Err != nil {
				wsconn.WriteJSON(model.NormalResponse{
					Code:    cons.SERVERERR,
					Message: msg.Err.Error(),
				})
				return
			}
			wsconn.WriteJSON(model.NormalResponse{
				Code: cons.SUCCESS,
				Data: msg.SyscallMsg,
			})

		}
	}()
	//开启生产者
	go host.PtraceProvider(msgchan, pid)

	//判断客户端是否断开，或者发出信号，如果有，关闭管道
	resp := model.NormalResponse{}
	for {
		_, p, err := wsconn.ReadMessage()
		if err != nil {
			seelog.Warn("readmsg err", err.Error())
			return
		}
		if err = json.Unmarshal(p, &resp); err != nil {
			seelog.Warn("client input error")
			continue
		}
		if resp.Code == cons.INTERRUPT {
			seelog.Info("client interrupt the websocket")
			return
		}
	}

}

func RegisterAlarmEventCtl(c *gin.Context) {
	alarmEvtReq := model.AlarmEvtRegisterRequest{}
	err := c.ShouldBindJSON(&alarmEvtReq)
	if err != nil {
		c.JSON(cons.HTTPHEADERCODE, model.NormalResponse{
			Code:    cons.UNMARSHALERR,
			Message: "request body unmarshal error !",
		})
		return
	}
	err = svc.RegisterAlarmEventSvc(alarmEvtReq)
	if err != nil {
		c.JSON(cons.HTTPHEADERCODE, model.NormalResponse{
			Code:    cons.SERVERERR,
			Message: err.Error(),
		})
		return
	}
	c.JSON(cons.HTTPHEADERCODE, model.NormalResponse{
		Code:    cons.SUCCESS,
		Message: "success",
	})
}
