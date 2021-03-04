package db

import (
	"github.com/globalsign/mgo"
	"ligomonitor/pkg/model"
)

type MongoCliStruct struct {
	DBClient *mgo.Session
}

func (m *MongoCliStruct) InsertProcessData(processData []model.Process) error {
	return nil
}

func (m *MongoCliStruct) GetProcessRectData(pid int) ([]model.Process, error) {
	return nil, nil
}
