package yw

import (
	"encoding/json"
	"errors"
	"github.com/Deansquirrel/goZlDianzqOfferTicketV3/common"
	"github.com/Deansquirrel/goZlDianzqOfferTicketV3/global"
	"github.com/Deansquirrel/goZlDianzqOfferTicketV3/object"
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

func getResponse(ctx iris.Context) (response object.ResponseCreateLittleTkt) {
	request, err := getRequestCreateLittleTktByContext(ctx)
	if err != nil {
		return getErrorResponse(request, ctx, err)
	}
	ok, err := checkReSubmit(&request)
	if err != nil {
		return getErrorResponse(request, ctx, err)
	}
	if ok {
		err = errors.New("券发放请求重复提交")
		return getErrorResponse(request, ctx, err)
	}
	err = request.CheckRequest()
	if err != nil {
		return getErrorResponse(request, ctx, err)
	}
	return GetResponseCreateLittleTkt(ctx, &request)
}

func getResponseData(response object.ResponseCreateLittleTkt) []byte {
	data, err := json.Marshal(response)
	if err != nil {
		common.PrintAndLog(err.Error())
		return []byte(err.Error())
	} else {
		return data
	}
}

func getErrorResponse(request object.RequestCreateLittleTkt, ctx iris.Context, err error) (response object.ResponseCreateLittleTkt) {
	response = GetResponseCreateLittleTktError(&request, err, ctx.GetStatusCode())
	return response
}

func checkReSubmit(request *object.RequestCreateLittleTkt) (bool, error) {
	val, err := global.Redis.Get(strconv.Itoa(global.RedisDbId1), request.AppId+request.Body.YwInfo.OprYwSno)
	if err != nil {
		if err.Error() != "redigo: nil returned" {
			return false, err
		}
	}
	if val != "" {
		return true, nil
	}
	return false, nil
}
func getRequestCreateLittleTktByContext(ctx iris.Context) (request object.RequestCreateLittleTkt, err error) {
	if ctx.URLParamExists("returntktno") {
		returnTktNo, err2 := ctx.URLParamInt("returntktno")
		if err2 != nil {
			err = errors.New("获取returntktno时发生错误," + err2.Error())
			return
		}
		request.ReturnTktNo = returnTktNo
	}
	request.Guid = ctx.GetHeader("requestGuid")
	request.AppId = ctx.GetHeader("appid")
	if request.AppId == "" {
		err = errors.New("appid不允许为空")
		return
	}
	err = ctx.ReadJSON(&request.Body)
	return
}
