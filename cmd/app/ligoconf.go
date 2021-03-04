/*
this package is to init the config
*/
package app

import (
	"fmt"
	"ligomonitor/pkg/conn"
	"ligomonitor/pkg/cons"
	"ligomonitor/pkg/model"
	"ligomonitor/pkg/service/dbsync"
	"ligomonitor/utils"
	_ "net/http/pprof"
	"os"
)

//定义全局变量
var LigoConf *model.LiGoMoniConf

func NewLigoConf(path string) {
	conf := model.LiGoMoniConf{}
	config := utils.ReadJsonFile(path, &conf).(*model.LiGoMoniConf)
	//check
	if config.DockerFlushTime <= 0 || config.DockerFlushTime <= 0 {
		fmt.Println(`please check param with "Flush",it could not be 0 or smaller than 0.`)
		os.Exit(cons.FLUSHPARAMERR)
	}
	if config.UseDB {
		if config.DBConf.DBTopFlush <= 0 || config.DBConf.DBDockerFlush <= 0 {
			fmt.Println(`please check db param with "dbConf.Flush",it could not be 0 or smaller than 0.`)
			os.Exit(cons.FLUSHPARAMERR)
		}
		//连接数据库,这里有mongo或mysql两种,TODO
		switch config.DBConf.DBType {
		case "mysql":
			conn.NewMysqlClient(config.DBConf.DBParams)
		case "mongo":
			conn.NewMongoClient(config.DBConf.DBParams)
		default:
			fmt.Println("not useable database ! please choose mysql or mongo")
			os.Exit(cons.DBCHOOSEERR)
		}
	}
	LigoConf = config
}

func SyncToDB() {
	if LigoConf.UseDB {
		//TODO
		dbsync.SyncDBTicker(LigoConf.DBConf)
	}
	return
}
