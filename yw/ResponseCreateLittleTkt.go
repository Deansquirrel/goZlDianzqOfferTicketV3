package yw

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/Deansquirrel/goZl"
	"github.com/Deansquirrel/goZlDianzqOfferTicketV3/common"
	"github.com/Deansquirrel/goZlDianzqOfferTicketV3/global"
	"github.com/Deansquirrel/goZlDianzqOfferTicketV3/object"
	"github.com/Deansquirrel/goZlDianzqOfferTicketV3/repository"
	"github.com/kataras/iris"
	"strconv"
	"strings"
)

type ResponseCreateLittleTkt struct {
	TmpCol    int                    `json:"tmpcol"`
	TktReturn []object.TktReturnInfo `json:"tktReturn"`

	ErrorModel  object.ErrorModel `json:"errormodel"`
	Description string            `json:"description"`

	HttpCode    int  `josn:"httpcode"`
	DBCommitted bool `json:"dbcommitted"`

	IsSuccess bool   `json:"issuccess"`
	FindAccid bool   `json:"findaccid"`
	HasCommit bool   `json:"hascommit"`
	Guid      string `json:"requestGuid"`
}

func GetResponseCreateLittleTkt(ctx iris.Context, request *RequestCreateLittleTkt) (response ResponseCreateLittleTkt) {
	response.BaseFunc(request, &response)
	response.refresh()

	//生成券号码
	tktNos, err := GetTktNoMulti(len(request.Body.CrmCardInfo))
	if err != nil {
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}
	//common.PrintOrLog("生成的券号码")
	//for index := range tktNos {
	//	common.PrintOrLog(tktNos[index])
	//}

	//记录redis,防止重复提交
	//common.PrintOrLog("准备记录Reids,防止重复提交")
	jsonNoList, err := json.Marshal(tktNos)
	if err != nil {
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}

	_, err = global.Redis.Set(strconv.Itoa(global.RedisDbId1), request.AppId+request.Body.YwInfo.OprYwSno, string(jsonNoList))
	if err != nil {
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}
	defer func() {
		err = global.Redis.Del(strconv.Itoa(global.RedisDbId1), request.AppId+request.Body.YwInfo.OprYwSno)
		if err != nil {
			common.PrintOrLog(err.Error())
		}
	}()

	//common.PrintOrLog(string(jsonNoList))

	//生成电子券系统流水号
	//common.PrintOrLog("准备生成电子券系统流水号")
	sno, err := GetSno("CT")
	if err != nil {
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}
	//common.PrintOrLog(sno)

	pzR := repository.PeiZhRepository{}
	dbConnList, err := pzR.GetXtMappingDbConnInfo(request.AppId, "DB_TicketHx", "TicketHx")
	if err != nil {
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}

	if len(dbConnList) < 1 {
		err = errors.New("传入的APPID无效（APPID错误或配置库缺少配置）")
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}

	t := conftools.ConfTools{}
	tktInfos := make([]object.TktInfo, 0)
	tktReturnInfos := make([]object.TktReturnInfo, 0)
	tktModels := make([]object.TktModel, 0)
	j := 0

	for _, val := range request.Body.CrmCardInfo {
		//common.PrintOrLog("CardNo : " + val.CardNo)
		accIdInput, err := t.DecryptFromBase64Format(val.CardNo, "accid")
		if err != nil {
			response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
			return
		}
		//common.PrintOrLog(accIdInput)
		accIdInputLong, err := strconv.Atoi(accIdInput)
		if err != nil {
			response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
			return
		}
		var tktInfo object.TktInfo
		if request.Body.TktInfo == nil || len(request.Body.TktInfo) < 1 {
			break
		} else {
			tktInfo = request.Body.TktInfo[0]
		}
		tktItem := object.TktInfo{
			AppId:    request.AppId,
			AccId:    accIdInputLong,
			TktKind:  tktInfo.TktKind,
			EffDate:  tktInfo.EffDate,
			Deadline: tktInfo.Deadline,
			CrYwLsh:  request.Body.YwInfo.OprYwSno,
			CrBr:     request.Body.YwInfo.OprBrid,
			CashMy:   tktInfo.CashMy,
			AddMy:    tktInfo.AddMy,
			TktName:  tktInfo.TktName,
			PCno:     tktInfo.PCno,
			TktNo:    tktNos[j],
		}
		tktInfos = append(tktInfos, tktItem)

		tktRetItem := object.TktReturnInfo{
			Sn:     val.Sn,
			TktNo:  tktNos[j],
			TktSno: sno,
		}
		tktReturnInfos = append(tktReturnInfos, tktRetItem)

		tktModel := object.TktModel{
			AccId:    accIdInputLong,
			AddMy:    tktInfo.AddMy,
			AppId:    request.AppId,
			CashMy:   tktInfo.CashMy,
			DeadLine: tktInfo.Deadline,
			PcNo:     tktInfo.PCno,
			EffTime:  tktInfo.EffDate,
			TktKind:  tktInfo.TktKind,
			TktName:  tktInfo.TktName,
			TktNo:    tktNos[j],
		}
		tktModels = append(tktModels, tktModel)

		j++
	}

	crTktInfo := object.TktCreateInfo{
		TktInfo:       make([]object.TktInfo, 0),
		TktYwInfo:     object.YwInfo{},
		TktReturnInfo: make([]object.TktReturnInfo, 0),
		CzLx:          2,
		CzLxSm:        "",
	}

	for _, val := range tktInfos {
		crTktInfo.TktInfo = append(crTktInfo.TktInfo, val)
	}
	crTktInfo.TktYwInfo = request.Body.YwInfo
	for _, val := range tktReturnInfos {
		crTktInfo.TktReturnInfo = append(crTktInfo.TktReturnInfo, val)
	}

	tkTs := object.TktModels{
		TktModels: make([]object.TktModel, 0),
	}
	for _, val := range tktModels {
		tkTs.TktModels = append(tkTs.TktModels, val)
	}

	//for _, db := range hxDbList {
	//	verInfo, err := hxR.GetVerInfo(db)
	//	if err != nil {
	//		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
	//		return
	//	}
	//	fmt.Println(verInfo)
	//}

	if err != nil {
		common.PrintOrLog(err.Error())
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}
	response = createLittleTktCreate(request.ReturnTktNo, crTktInfo, tkTs, ctx, request)
	return
}

func createLittleTktCreate(returnTktNo int, crTktInfo object.TktCreateInfo, tktModels object.TktModels, ctx iris.Context, request *RequestCreateLittleTkt) (response ResponseCreateLittleTkt) {
	if crTktInfo.TktInfo == nil || len(crTktInfo.TktInfo) < 1 {
		response = GetResponseCreateLittleTktError(request, errors.New("TktInfo列表不能为空"), ctx.GetStatusCode())
		return
	}
	appId := crTktInfo.TktInfo[0].AppId
	for _, val := range crTktInfo.TktInfo {
		if val.AppId != appId {
			response = GetResponseCreateLittleTktError(request, errors.New("一次请求APPID必须相同"), ctx.GetStatusCode())
			return
		}
	}

	pzR := repository.PeiZhRepository{}
	dbConnList, err := pzR.GetXtMappingDbConnInfo(request.AppId, "DB_TicketHx", "TicketHx")
	if err != nil {
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}

	if len(dbConnList) < 1 {
		err = errors.New("传入的APPID无效（APPID错误或配置库缺少配置）")
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}

	hxR := repository.HeXRepository{}
	dbConn := make([]*sql.DB, 0)
	for _, val := range dbConnList {
		db, err := hxR.GetDbConnByString(val.MConnStr)
		if err != nil {
			response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
			return
		}
		dbConn = append(dbConn, db)
	}
	if dbConn == nil || len(dbConn) < 1 {
		response = GetResponseCreateLittleTktError(request, errors.New("未获取到有效的hx库连接"), ctx.GetStatusCode())
		return
	}

	//===============================================================================================
	//执行Hx库存储过程
	hxDb := dbConn[0]
	err = hxR.CreateLittleTktCreate(hxDb, crTktInfo.TktInfo)
	if err != nil {
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}

	//===============================================================================================
	response = ResponseCreateLittleTkt{TktReturn: make([]object.TktReturnInfo, 0)}
	err = redisSetTktModel(tktModels.TktModels, global.RedisDbId1)
	if err != nil {
		common.PrintAndLog(err.Error())
	}
	//===============================================================================================
	if returnTktNo == 1 && len(crTktInfo.TktReturnInfo) > 0 {
		for _, val := range crTktInfo.TktReturnInfo {
			response.TktReturn = append(response.TktReturn, val)
		}
	}
	//===============================================================================================
	err = rabbitMqSetTktInfo(crTktInfo)
	if err != nil {
		common.PrintAndLog(err.Error())
	}
	//===============================================================================================
	//response = GetResponseCreateLittleTktError(request, errors.New("Test End"), ctx.GetStatusCode())
	return
}

func rabbitMqSetTktInfo(tktInfo object.TktCreateInfo) error {
	valStr, err := json.Marshal(tktInfo)
	if err != nil {
		common.PrintOrLog(string(valStr))
	}
	err = global.RabbitMQ.Publish("", "amq.fanout", tktInfo.MessageRoute, string(valStr))
	if err != nil {
		common.PrintOrLog(string(valStr))
	}
	return nil
}

func redisSetTktModel(tktModels []object.TktModel, redisDbid int) error {
	tktInputModels := make(map[string]object.TktModel)
	for _, val := range tktModels {
		tktInputModels[val.AppId+"|"+val.TktNo] = val
	}

	for key, val := range tktInputModels {
		str, err := json.Marshal(val)
		if err != nil {
			common.PrintOrLog(err.Error())
		}
		_, err = global.Redis.Set(strconv.Itoa(redisDbid), key, string(str))
		if err != nil {
			common.PrintOrLog(err.Error())
		}
	}

	for _, val := range tktModels {
		if val.AccId != 0 {
			str, err := json.Marshal(val)
			if err != nil {
				common.PrintOrLog(err.Error())
				continue
			}
			sKey := "AA|" + val.AppId + "|" + strconv.Itoa(val.AccId)
			_, err = global.Redis.Set(strconv.Itoa(redisDbid), sKey, string(str))
			if err != nil {
				common.PrintOrLog(err.Error())
				continue
			}
		}
	}

	return nil
}

func GetResponseCreateLittleTktError(request *RequestCreateLittleTkt, err error, httpCode int) (response ResponseCreateLittleTkt) {
	response.BaseFunc(request, &response)
	response.HttpCode = httpCode
	response.ErrorModel = object.ErrorModel{
		ErrType: "1",
		Desc:    err.Error(),
	}
	response.refresh()
	return
}

func (response *ResponseCreateLittleTkt) BaseFunc(req *RequestCreateLittleTkt, resp *ResponseCreateLittleTkt) {
	resp.Guid = req.Guid
}

func (response *ResponseCreateLittleTkt) refresh() {
	//if (response.HttpCode == 503 || response.HttpCode == 403) && response.Description != "" {
	//	if strings.Contains(response.Description,"账号无法找到") {
	//		panic(response.Description)
	//	}
	//}

	if response.ErrorModel.ErrType == "" && response.ErrorModel.Desc == "" {
		response.DBCommitted = true
	} else {
		response.DBCommitted = false
	}

	if response.ErrorModel.ErrType != "" || response.ErrorModel.Desc != "" {
		response.IsSuccess = false
	} else if response.HttpCode != 200 {
		response.IsSuccess = false
	} else {
		response.IsSuccess = true
	}

	if response.HttpCode == 503 && strings.Contains(response.Description, "账号无法找到") {
		response.FindAccid = false
	} else {
		response.FindAccid = true
	}
}
