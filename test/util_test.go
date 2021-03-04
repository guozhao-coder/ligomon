package test

import (
	"fmt"
	"ligomonitor/utils"
	"os"
	"testing"
)

func Test_err(t *testing.T) {
	_, err := os.Open("/proc/aaa")
	if err != nil {
		err = utils.ErrJoint("open aaa err :", err)
		fmt.Println(err.Error())
	}
}
