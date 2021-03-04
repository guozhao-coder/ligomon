package db

import "ligomonitor/pkg/model"

type DBOperate interface {
	InsertProcessData(processData []model.Process) error
	GetProcessRectData(pid int) ([]model.Process, error)
}
