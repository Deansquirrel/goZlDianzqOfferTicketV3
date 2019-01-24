package main

import (
	"fmt"
	"github.com/Deansquirrel/goZlDianzqOfferTicketV3/common"
	"github.com/Deansquirrel/goZlDianzqOfferTicketV3/global"
	"github.com/Deansquirrel/goZlDianzqOfferTicketV3/yw"
)

func main() {
	//==================================================================================================================
	config, err := common.GetSysConfig("config.toml")
	if err != nil {
		fmt.Println("加载配置文件时遇到错误：" + err.Error())
		return
	}
	global.SysConfig = config
	configStr, err := global.SysConfig.GetConfigStr()
	if err != nil {
		common.PrintOrLog(err.Error())
	} else {
		common.PrintOrLog(configStr)
	}
	err = yw.RefreshConfig(global.SysConfig)
	if err != nil {
		common.PrintAndLog("加载配置时遇到错误：" + err.Error())
		return
	}
	//==================================================================================================================
	common.PrintOrLog("程序启动")
	defer common.PrintOrLog("程序退出")
	//==================================================================================================================
	yw.StartWebServer()
}
