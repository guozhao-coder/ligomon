package dbsync

import (
	"fmt"
	"ligomonitor/pkg/model"
	"time"
)

var DBCheck *model.DBConfig

func SyncDBTicker(dbconfig *model.DBConfig) {
	topTicket := time.NewTicker(time.Minute * time.Duration(dbconfig.DBTopFlush))
	dockTicket := time.NewTicker(time.Minute * time.Duration(dbconfig.DBDockerFlush))
	DBCheck = dbconfig
	for {
		select {
		case <-topTicket.C:
			go syncProcess()
		case <-dockTicket.C:
			fmt.Println("docker ticker")
		}
	}

}
