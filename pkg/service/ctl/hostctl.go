package ctl

import (
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"ligomonitor/cmd/app"
	"ligomonitor/pkg/cons"
	"ligomonitor/pkg/model"
	"ligomonitor/pkg/service/svc"
	"ligomonitor/utils"
	"time"
)

func GetHostResourceCtl(c *gin.Context) {
	hr, err := svc.GetHostResourceSvc()
	if err != nil {
		c.JSON(cons.HTTPHEADERCODE, model.NormalResponse{
			Code:    cons.SERVERERR,
			Message: err.Error(),
		})
		return
	}
	c.JSON(cons.SUCCESS, model.NormalResponse{
		Code:    cons.SUCCESS,
		Message: "success",
		Data:    hr,
	})
}

func GetHostResourceStreamCtl(c *gin.Context) {
	wsconn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		err = utils.ErrJoint("websocket conn error : ", err)
		seelog.Error(err.Error())
		c.JSON(cons.SUCCESS, model.NormalResponse{
			Code:    cons.WSCONNERR,
			Message: err.Error(),
		})
		return
	}
	defer wsconn.Close()
	topTicker := time.NewTicker(time.Second * time.Duration(app.LigoConf.TopFlushTime))
	msgChan := make(chan model.HostResourceMsg, 100)
	defer close(msgChan)
	go func() {
		for {
			hostMsg, ifClose := <-msgChan
			if !ifClose {
				seelog.Info("the host hostResourceMsg channel is close")
				return
			}
			if hostMsg.Err != nil {
				wsconn.WriteJSON(model.NormalResponse{
					Code:    cons.SERVERERR,
					Message: hostMsg.Err.Error(),
				})
			}
			wsconn.WriteJSON(model.NormalResponse{
				Code:    cons.SUCCESS,
				Message: "success",
				Data:    hostMsg.Resource,
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
			go svc.GetHostResourceStreamSvc(msgChan)
		}
	}
}
