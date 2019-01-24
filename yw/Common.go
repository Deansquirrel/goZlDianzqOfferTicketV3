package yw

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/Deansquirrel/go-tool"
	"github.com/Deansquirrel/goZl"
	"github.com/Deansquirrel/goZlDianzqOfferTicketV3/common"
	"github.com/Deansquirrel/goZlDianzqOfferTicketV3/global"
	"github.com/Deansquirrel/goZlDianzqOfferTicketV3/object"
	"github.com/Deansquirrel/goZlDianzqOfferTicketV3/repository"
	"github.com/kataras/iris"
	"strconv"
	"strings"
	"time"
)

func RefreshConfig(config *object.SysConfig) error {
	global.SysConfig = config
	//============================================================
	conn, err := getPeiZhDbConn()
	if err != nil {
		return err
	}
	global.PeiZhDbConn = conn
	//============================================================
	pZhR := repository.PeiZhRepository{}
	//============================================================
	global.SnoServer, err = pZhR.GetXtWxAppIdJoinInfo(global.PeiZhDbConn, global.SysConfig.Total.JPeiZh, "SnoServer", 0)
	if err != nil {
		return err
	}
	snoWorkIdStr, err := pZhR.GetXtWxAppIdJoinInfo(global.PeiZhDbConn, global.SysConfig.Total.JPeiZh, "WorkerId", 0)
	if err != nil {
		return err
	}
	global.SnoWorkerId, err = strconv.Atoi(snoWorkIdStr)
	if err != nil {
		return err
	}
	//============================================================
	redisConfigStr, err := pZhR.GetXtWxAppIdJoinInfo(global.PeiZhDbConn, global.SysConfig.Total.JPeiZh, "SERedis", 0)
	if err != nil {
		return err
	}
	redisConfig := strings.Split(redisConfigStr, "|")
	if len(redisConfig) != 2 {
		return errors.New("redis配置参数异常.expected 2 , got " + strconv.Itoa(len(redisConfig)))
	}

	global.Redis = go_tool.NewRedis(redisConfig[0], redisConfig[1], 5000, 5000, 5)
	if err != nil {
		return err
	}

	redisDbId1Str, err := pZhR.GetXtWxAppIdJoinInfo(global.PeiZhDbConn, global.SysConfig.Total.JPeiZh, "RedisDbId1", 0)
	if err != nil {
		return err
	}
	global.RedisDbId1, err = strconv.Atoi(redisDbId1Str)
	if err != nil {
		return err
	}

	redisDbId2Str, err := pZhR.GetXtWxAppIdJoinInfo(global.PeiZhDbConn, global.SysConfig.Total.JPeiZh, "RedisDbId2", 0)
	if err != nil {
		return err
	}
	global.RedisDbId2, err = strconv.Atoi(redisDbId2Str)
	if err != nil {
		return err
	}
	//============================================================
	rabbitMQConfigStr, err := pZhR.GetXtWxAppIdJoinInfo(global.PeiZhDbConn, global.SysConfig.Total.JPeiZh, "RabbitConnection", 0)
	if err != nil {
		return err
	}
	rabbitMQConfig := strings.Split(rabbitMQConfigStr, "|")
	if len(rabbitMQConfig) != 5 {
		return errors.New("rabbitMQ配置参数异常.expected 5 , got " + strconv.Itoa(len(rabbitMQConfig)))
	}
	rabbitMQPort, err := strconv.Atoi(rabbitMQConfig[1])
	if err != nil {
		return err
	}
	global.RabbitMQ = go_tool.NewRabbitMQ(rabbitMQConfig[3], rabbitMQConfig[4], rabbitMQConfig[0], rabbitMQPort, rabbitMQConfig[2], time.Second*60, time.Millisecond*500, 3, time.Second*5)

	err = rabbitMqInit()
	if err != nil {
		return err
	}
	//============================================================

	return nil
}

func getPeiZhDbConn() (*sql.DB, error) {
	return repository.GetDbConn(global.SysConfig.PeiZhDb.Server,
		global.SysConfig.PeiZhDb.Port,
		global.SysConfig.PeiZhDb.DbName,
		global.SysConfig.PeiZhDb.User,
		global.SysConfig.PeiZhDb.Password)
}

func rabbitMqInit() error {
	conn, err := global.RabbitMQ.GetConn()
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()

	err = global.RabbitMQ.QueueDeclareSimple(conn, "TktCreateYwdetail")
	if err != nil {
		return err
	}

	err = global.RabbitMQ.QueueBind(conn, "TktCreateYwdetail", "", "amq.fanout", true)
	if err != nil {
		return err
	}

	err = global.RabbitMQ.AddProducer("")
	if err != nil {
		return err
	}

	//err = global.RabbitMQ.AddConsumer("","TktCreateYwdetail",lsHandler)

	return nil
}

func GetResponseCreateLittleTkt(ctx iris.Context, request *object.RequestCreateLittleTkt) (response object.ResponseCreateLittleTkt) {
	response.BaseFunc(request, &response)
	response.Refresh()
	//==================================================================================================================
	//生成券号码
	tktNos, err := GetTktNoMulti(len(request.Body.CrmCardInfo))
	if err != nil {
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}
	//==================================================================================================================
	//记录redis,防止重复提交
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
	//==================================================================================================================
	//生成电子券系统流水号
	//common.PrintOrLog("准备生成电子券系统流水号")
	sno, err := GetSno("CT")
	if err != nil {
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}
	//==================================================================================================================
	t := conftools.ConfTools{}
	tktInfos := make([]object.TktInfo, 0)
	tktReturnInfos := make([]object.TktReturnInfo, 0)
	tktModels := make([]object.TktModel, 0)
	j := 0

	for _, val := range request.Body.CrmCardInfo {
		accIdInput, err := t.DecryptFromBase64Format(val.CardNo, "accid")
		if err != nil {
			response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
			return
		}
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

	if err != nil {
		common.PrintOrLog(err.Error())
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}
	response = createLittleTktCreate(request.ReturnTktNo, crTktInfo, tkTs, ctx, request)
	return
	//==================================================================================================================
}

func createLittleTktCreate(returnTktNo int, crTktInfo object.TktCreateInfo, tktModels object.TktModels, ctx iris.Context, request *object.RequestCreateLittleTkt) (response object.ResponseCreateLittleTkt) {
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
	//==================================================================================================================
	pzR := repository.PeiZhRepository{}
	dbConnList, err := pzR.GetXtMappingDbConnInfo(request.AppId)
	if err != nil {
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}
	if len(dbConnList) < 1 {
		err = errors.New("传入的APPID无效（APPID错误或配置库缺少配置）")
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}
	//==================================================================================================================
	//执行Hx库存储过程
	hxDb, err := getHxConn(dbConnList[0])
	if err != nil {
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}
	hxR := repository.HeXRepository{}
	err = hxR.CreateLittleTktCreate(hxDb, crTktInfo.TktInfo)
	if err != nil {
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}
	//==================================================================================================================
	response = object.ResponseCreateLittleTkt{TktReturn: make([]object.TktReturnInfo, 0)}
	err = redisSetTktModel(tktModels.TktModels, global.RedisDbId1)
	if err != nil {
		common.PrintAndLog(err.Error())
	}
	//==================================================================================================================
	if returnTktNo == 1 && len(crTktInfo.TktReturnInfo) > 0 {
		for _, val := range crTktInfo.TktReturnInfo {
			response.TktReturn = append(response.TktReturn, val)
		}
	}
	//==================================================================================================================
	err = rabbitMqSetTktInfo(crTktInfo)
	if err != nil {
		common.PrintAndLog(err.Error())
	}
	//==================================================================================================================
	return
}

func getHxConn(configStr string) (*sql.DB, error) {
	if global.HxDbConnMap == nil {
		global.HxDbConnMap = make(map[string]*sql.DB)
	}
	if _, ok := global.HxDbConnMap[configStr]; ok {
		return global.HxDbConnMap[configStr], nil
	}

	config := strings.Split(configStr, "|")
	if len(config) != 5 {
		common.PrintAndLog("数据库配置串解析失败 - " + configStr)
		return nil, errors.New("数据库配置串解析失败")
	}
	port, err := strconv.Atoi(config[1])
	if err != nil {
		common.PrintAndLog("数据库配置串端口解析失败 - " + configStr)
		return nil, errors.New("数据库配置串端口解析失败")
	}
	conn, err := repository.GetDbConn(config[0], port, config[2], config[3], config[4])
	if err != nil {
		return nil, err
	}
	global.HxDbConnMap[configStr] = conn
	return conn, nil
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

func GetResponseCreateLittleTktError(request *object.RequestCreateLittleTkt, err error, httpCode int) (response object.ResponseCreateLittleTkt) {
	response.BaseFunc(request, &response)
	response.HttpCode = httpCode
	response.ErrorModel = object.ErrorModel{
		ErrType: "1",
		Desc:    err.Error(),
	}
	response.Refresh()
	return
}
