package dbsync

import (
	"github.com/cihub/seelog"
	"ligomonitor/pkg/conn"
	"ligomonitor/pkg/service/db"
	"ligomonitor/pkg/service/host"
	"ligomonitor/utils"
	"sync"
	"time"
)

var dbCli db.DBOperate
var dbCliOnce sync.Once

func NewDBCli() db.DBOperate {
	dbCliOnce.Do(func() {
		switch DBCheck.DBType {
		case "mysql":
			dbCli = &db.MysqlCliStruct{DBClient: conn.MysqlClient}
		case "mongo":
			dbCli = &db.MongoCliStruct{DBClient: conn.MgoClient}
		}
	})

	return dbCli
}

func syncProcess() {
	//获取进程信息
	seelog.Info("-----------------同步开始", time.Now())
	defer seelog.Info("-----------------同步结束", time.Now())
	dbCli := NewDBCli()
	processes, err := host.GetProcessTotalInfo(0, true)
	if err != nil {
		err = utils.ErrJoint("sync operation,get process info err ", err)
		seelog.Error(err.Error())
		return
	}
	//同步进程信息
	err = dbCli.InsertProcessData(processes)
	if err != nil {
		err = utils.ErrJoint("sync operation, insert into table err", err)
		seelog.Error(err.Error())
		return
	}
}

func syncDocker() {

}
