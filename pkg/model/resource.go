package model

type CPU struct {
	UserUsage   float32 `json:"userUsage"`   //用户空间占比
	NiceUsage   float32 `json:"niceUsage"`   //改变过优先级的进程的占比
	SystemUsage float32 `json:"systemUsage"` //内核空间占比
	IdleUsage   float32 `json:"idleUsage"`   //空闲占比
	IowaitUsage float32 `json:"iowaitUsage"` //等待占比
	HIrqUsage   float32 `json:"hIrqUsage"`   //硬中断占比
	SIrqUsage   float32 `json:"sIrqUsage"`   //软中断占比
}

type HostResourceUsed struct {
	CPUUsed *CPU        `json:"cpuUsed"`
	MemUsed *MemAndSwap `json:"memUsed"`
}

type HostResourceMsg struct {
	Resource HostResourceUsed
	Err      error
}

type Memory struct {
	Total     int32 `json:"total"`     //总共的内存
	Free      int32 `json:"free"`      //剩余的内存
	Used      int32 `json:"used"`      //使用的内存(包括缓存)
	BuffCache int32 `json:"buffCache"` //buff/cache
}

type Swap struct {
	Total int32 `json:"total"`
	Used  int32 `json:"used"`
	Free  int32 `json:"free"`
}

type MemAndSwap struct {
	Mem *Memory `json:"mem"`
	Swa *Swap   `json:"swa"`
}
