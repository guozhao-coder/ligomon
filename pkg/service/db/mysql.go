package db

import (
	"database/sql"
	"github.com/cihub/seelog"
	"ligomonitor/pkg/model"
	"ligomonitor/utils"
	"time"
)

type MysqlCliStruct struct {
	DBClient *sql.DB
}

func (m *MysqlCliStruct) InsertProcessData(processData []model.Process) error {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error(err)
		}
	}()
	begin, err := m.DBClient.Begin()
	if err != nil {
		err = utils.ErrJoint("begin error : ", err)
		seelog.Error(err.Error())
		begin.Rollback()
		return err
	}
	for i := 0; i < len(processData); i++ {
		timeNow := time.Now().Unix()
		sqlStr := "insert into process_info values(default,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
		result, err := begin.Exec(sqlStr, processData[i].Pid, processData[i].PPid, processData[i].Name, processData[i].Tgid, processData[i].State, processData[i].Uid, processData[i].Gid, processData[i].Threads, processData[i].VmPeak, processData[i].VmSize, processData[i].VmHWM, processData[i].VmRss, processData[i].VmSwap, processData[i].VoluntaryCS, processData[i].NoVoluntaryCS, processData[i].CPUUsage, timeNow)
		if err != nil {
			err = utils.ErrJoint("mysql exec error : ", err)
			seelog.Error(err)
			begin.Rollback()
			return err
		}
		if affe, err := result.RowsAffected(); err != nil || affe != 1 {
			err = utils.ErrJoint("affected error : ", err)
			seelog.Error(err)
			begin.Rollback()
			return err
		}
	}
	begin.Commit()
	return nil
}

func (m *MysqlCliStruct) GetProcessRectData(pid int) ([]model.Process, error) {
	defer func() {
		if err := recover(); err != nil {
			seelog.Error(err)
		}
	}()

	begin, err := m.DBClient.Begin()
	if err != nil {
		err = utils.ErrJoint("begin error : ", err)
		seelog.Error(err.Error())
		begin.Rollback()
		return nil, err
	}
	sqlStr := "select ppid,name,state,vm_rss,cpu_usage,time from process_info where pid = ? order by time desc limit 20"
	rows, err := begin.Query(sqlStr, pid)
	if err != nil {
		err = utils.ErrJoint("query error : ", err)
		seelog.Error(err.Error())
		begin.Rollback()
		return nil, err
	}
	var ppid sql.NullInt32
	var name sql.NullString
	var state sql.NullString
	var vmRss sql.NullInt64
	var cpuUsage sql.NullFloat64
	var time sql.NullInt32
	procs := []model.Process{}
	proc := model.Process{}
	for rows.Next() {
		err := rows.Scan(&ppid, &name, &state, &vmRss, &cpuUsage, &time)
		if err != nil {
			err = utils.ErrJoint("rows scan err ", err)
			seelog.Error(err.Error())
			begin.Rollback()
			return nil, nil
		}
		//添加到数组
		proc.PPid = int(ppid.Int32)
		proc.Name = name.String
		proc.State = state.String
		proc.VmRss = vmRss.Int64
		proc.CPUUsage = float32(cpuUsage.Float64)
		proc.Time = int(time.Int32)
		procs = append(procs, proc)
	}
	return procs, nil
}
