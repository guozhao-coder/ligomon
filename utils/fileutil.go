package utils

import (
	"encoding/json"
	"fmt"
	"ligomonitor/pkg/cons"
	"os"
)

func ReadJsonFile(fileDir string, data interface{}) interface{} {
	fmt.Println("ReadJsonFile open json file:" + fileDir)
	r, err := os.Open(fileDir)
	if err != nil {
		fmt.Println("open jsonfile error：" + err.Error())
		os.Exit(cons.OPENFILEERR)
	}
	decoder := json.NewDecoder(r)
	err = decoder.Decode(data)
	if err != nil {
		fmt.Println("decode json file error：" + err.Error())
		os.Exit(cons.DECODEJSONERR)
	}
	return data
}
