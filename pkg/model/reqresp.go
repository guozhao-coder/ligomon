package model

type NormalResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type AlarmEvtRegisterRequest struct {
	Pid      int     `json:"pid"`      //进程
	VMLimit  int64   `json:"vmLimit"`  //内存最大限制
	CPULimit float32 `json:"cpuLimit"` //cpu最大限制
	Operate  int     `json:"operate"`  //告警后的操作
}
