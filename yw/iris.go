package yw

import (
	"encoding/json"
	"github.com/Deansquirrel/goZlDianzqOfferTicketV3/common"
	"github.com/Deansquirrel/goZlDianzqOfferTicketV3/global"
	"github.com/kataras/iris"
	"strconv"
)

func StartWebServer() {
	app := iris.New()
	app.Post("/CreateLittleTkt", handler)
	addr := ":" + strconv.Itoa(global.SysConfig.Total.Port)
	err := app.Run(iris.Addr(addr), iris.WithoutServerError(iris.ErrServerClosed), iris.WithOptimizations)
	if err != nil {
		common.PrintAndLog(err.Error())
	}
}

func handler(ctx iris.Context) {
	response := getResponse(ctx)
	_, err := ctx.Write(getResponseData(response))
	if err != nil {
		common.PrintAndLog(err.Error())
	}
	return
}

func getResponse(ctx iris.Context) (response ResponseCreateLittleTkt) {
	request, err := GetRequestCreateLittleTktByContext(ctx)
	if err != nil {
		return getErrorResponse(request, ctx, err)
	}
	err = request.CheckRequest()
	if err != nil {
		return getErrorResponse(request, ctx, err)
	}
	return GetResponseCreateLittleTkt(ctx, &request)
}

func getResponseData(response ResponseCreateLittleTkt) []byte {
	data, err := json.Marshal(response)
	if err != nil {
		common.PrintAndLog(err.Error())
		return []byte(err.Error())
	} else {
		return data
	}
}

func getErrorResponse(request RequestCreateLittleTkt, ctx iris.Context, err error) (response ResponseCreateLittleTkt) {
	response = GetResponseCreateLittleTktError(&request, err, ctx.GetStatusCode())
	return response
}
