package cons

const (
	HTTPHEADERCODE = 200
)

const (
	SUCCESS = 200 //成功标志

	/*websocket  code*/
	WSCONNERR = 300
	HEARTBEAT = 301 //服务器发送的心跳信号
	INTERRUPT = 302 //来自客户端发来的中断，用于中断trace

	/*normal code*/
	SERVERERR   = 500          //服务端出错
	URLPARAMERR = 5001 << iota //url解析出错
	UNMARSHALERR
)

type AlarmSignal int

const (
	KILLSIG AlarmSignal = 1
	MailSIG AlarmSignal = 2
)
