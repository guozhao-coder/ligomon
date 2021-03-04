/*
这个包主要作用是提供连接对象，包括数据库连接，以及其他客户端连接如docker
*/

package conn

import (
	"database/sql"
	"fmt"
	"github.com/globalsign/mgo"
	_ "github.com/go-sql-driver/mysql"
	"ligomonitor/pkg/cons"
	"ligomonitor/pkg/model"
	"os"
	"time"
)

var MysqlClient *sql.DB

var MgoClient *mgo.Session

func NewMysqlClient(conf *model.DBParam) {
	dsName := fmt.Sprintf("%s:%s@(%s:%s)/%s", conf.DBUser, conf.DBPwd, conf.DBIP, conf.DBPort, conf.DBName)
	db, err := sql.Open("mysql", dsName)
	if err != nil {
		fmt.Println("mysql open error : ", err.Error())
		os.Exit(cons.MYSQLCONNERR)
	}
	if err := db.Ping(); err != nil {
		fmt.Println("mysql ping error : ", err.Error())
		os.Exit(cons.MYSQLCONNERR)
	}
	fmt.Println("mysql conn success")
	MysqlClient = db
	MysqlClient.SetMaxIdleConns(10)
	MysqlClient.SetMaxOpenConns(50)
	//init the tables
	//todo
	_, err = MysqlClient.Exec(MySQLDBCreate)
	if err != nil{
		fmt.Println("mysql init table error:",err.Error())
	}
}

func NewMongoClient(conf *model.DBParam) error {
	mgoAddr := conf.DBIP + ":" + conf.DBPort
	mgoClient, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:       []string{mgoAddr},
		Database:    conf.DBName,
		Username:    conf.DBUser,
		Password:    conf.DBPwd,
		MinPoolSize: 2048,
		PoolLimit:   2048,
		Timeout:     10 * time.Second,
	})
	if err != nil {
		fmt.Println("mongo open error : ", err.Error())
		os.Exit(cons.MONGOCONNERR)
	}
	if err := mgoClient.Ping(); err != nil {
		fmt.Println("mongo ping error : ", err.Error())
		os.Exit(cons.MONGOCONNERR)
	}
	fmt.Println("mongo conn success")
	MgoClient = mgoClient
	return nil
}

var MySQLDBCreate = `create table if not exists process_info(
id int not null auto_increment primary key,
pid int not null,
ppid int not null,
name text,
tgid int,
state varchar(10),
uid int,
gid int,
threads int,
vm_peak int,
vm_size int,
vm_hwm int,
vm_rss int,
vm_swap int,
voluntary_cs int,
no_voluntary_cs int,
cpu_usage double,
time int
)`
