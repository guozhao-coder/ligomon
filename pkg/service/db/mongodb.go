package db

import (
	"github.com/cihub/seelog"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"ligomonitor/pkg/model"
	"ligomonitor/utils"
)

type MongoCliStruct struct {
	DBClient *mgo.Session
}

func (m *MongoCliStruct) InsertProcessData(processData []model.Process) error {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error(err)
		}
	}()
	coll := m.DBClient.Copy().DB("").C("process_info")
	for i := 0; i < len(processData); i++ {
		err := coll.Insert(processData[i])
		if err != nil {
			err = utils.ErrJoint("mongo insert error : ", err)
			seelog.Error(err.Error())
		}
	}
	return nil
}

func (m *MongoCliStruct) GetProcessRectData(pid int) ([]model.Process, error) {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error(err)
		}
	}()
	procs := []model.Process{}
	coll := m.DBClient.Copy().DB("").C("process_info")
	err := coll.Find(bson.M{"pid": pid}).All(&procs)
	if err != nil {
		err = utils.ErrJoint("mongo query error : ", err)
		seelog.Error(err.Error())
		return []model.Process{}, nil
	}
	return procs, nil
}
