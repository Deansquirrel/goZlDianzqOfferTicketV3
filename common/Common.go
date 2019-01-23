package common

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/Deansquirrel/go-tool"
	"github.com/Deansquirrel/goZlDianzqOfferTicketV3/global"
	"github.com/Deansquirrel/goZlDianzqOfferTicketV3/object"
)

func PrintAndLog(msg string) {
	fmt.Println(msg)
	if global.SysConfig.Total.IsDebug {
		err := go_tool.Log(msg)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func PrintOrLog(msg string) {
	if global.SysConfig.Total.IsDebug {
		err := go_tool.Log(msg)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println(msg)
	}
}

func GetSysConfig(fileName string) (*object.SysConfig, error) {
	path, err := go_tool.GetCurrPath()
	if err != nil {
		return nil, err
	}
	var config object.SysConfig
	_, err = toml.DecodeFile(path+"\\"+fileName, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
