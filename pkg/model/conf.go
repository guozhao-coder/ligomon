package model

type LiGoMoniConf struct {
	HTTPPort        int        `json:"httpPort"`
	PprofPort       int        `json:"pprofPort"`
	UseDB           bool       `json:"useDB"`
	DBConf          *DBConfig  `json:"dbConf"`
	TopFlushTime    int32      `json:"topFlushTime"`
	DockerFlushTime int32      `json:"dockerFlushTime"`
	LogPath         string     `json:"logPath"`
	MailConf        *MailParam `json:"mailConf"`
}

type DBConfig struct {
	DBType        string   `json:"dbType"`
	DBParams      *DBParam `json:"dbParams"`
	DBTopFlush    int32    `json:"dbTopFlush"`
	DBDockerFlush int32    `json:"dbDockerFlush"`
}

type DBParam struct {
	DBIP   string `json:"dbIP"`
	DBPort string `json:"dbPort"`
	DBUser string `json:"dbUser"`
	DBName string `json:"dbName"`
	DBPwd  string `json:"dbPwd"`
}

type MailParam struct {
	User string `json:"user"`
	Pass string `json:"pass"`
	Host string `json:"host"`
	Port string `json:"port"`
}
