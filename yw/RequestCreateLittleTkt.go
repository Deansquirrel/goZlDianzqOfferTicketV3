package yw

import (
	"errors"
	"github.com/Deansquirrel/goZlDianzqOfferTicketV3/common"
	"github.com/Deansquirrel/goZlDianzqOfferTicketV3/global"
	"github.com/Deansquirrel/goZlDianzqOfferTicketV3/object"
	"github.com/kataras/iris"
	"strconv"
)

type RequestCreateLittleTkt struct {
	Guid        string
	AppId       string
	ReturnTktNo int
	Body        requestBody
}

type requestBody struct {
	TmpCol      int
	TktInfo     []object.TktInfo
	CrmFqYwInfo object.CrmFqYwInfo
	CrmCardInfo []object.CrmCardInfo
	MdFqYwInfo  []object.MdFqYwInfo
	YwInfo      object.YwInfo
}

func GetRequestCreateLittleTktByContext(ctx iris.Context) (request RequestCreateLittleTkt, err error) {
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

func (request *RequestCreateLittleTkt) CheckRequest() error {
	if request.Body.TktInfo == nil || request.Body.CrmFqYwInfo.CrmJhh == "" || request.Body.CrmCardInfo == nil || request.Body.YwInfo.OprYwSno == "" {
		common.PrintOrLog("传入的记录集为空")
		return errors.New("传入的记录集为空")
	}
	for index := range request.Body.CrmCardInfo {
		if request.Body.CrmCardInfo[index].CardType != 5 {
			common.PrintOrLog("会员券目前仅支持按账户id加密值操作")
			return errors.New("会员券目前仅支持按账户id加密值操作")
		}
	}
	if len(request.Body.CrmCardInfo) > 100 {
		common.PrintOrLog("券发放（立即生效）禁止超过100张")
		return errors.New("券发放（立即生效）禁止超过100张")
	}

	val, err := global.Redis.Get(strconv.Itoa(global.RedisDbId1), request.AppId+request.Body.YwInfo.OprYwSno)
	if err != nil {
		if err.Error() != "redigo: nil returned" {
			common.PrintOrLog("保存Redis时发生错误")
			return err
		}
	}
	if val != "" {
		return errors.New("券发放请求重复提交")
	}
	return nil
}
