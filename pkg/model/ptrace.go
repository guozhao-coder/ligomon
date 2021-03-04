package model

type PtraceMsg struct {
	SyscallMsg string `json:"syscallMsg"`
	Err        error  `json:"err"`
}
