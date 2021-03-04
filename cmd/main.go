package main

import (
	"ligomonitor/cmd/app"
	"ligomonitor/pkg/api"
	"ligomonitor/utils"
	_ "net/http/pprof"
)

func main() {
	//init the seelog config
	utils.GetLogConfig("/app/GoWork/ligomonitor/configs/logcfg.xml")
	//init the gloable config
	app.NewLigoConf("/app/GoWork/ligomonitor/configs/conf.json")
	//sync data to the database
	go app.SyncToDB()

	//start the router
	if err := api.StartRouter(); err != nil {
		return
	}

}
