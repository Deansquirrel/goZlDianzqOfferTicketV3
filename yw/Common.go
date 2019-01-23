package yw

import (
	"database/sql"
	"errors"
	"github.com/Deansquirrel/go-tool"
	"github.com/Deansquirrel/goZlDianzqOfferTicketV3/global"
	"github.com/Deansquirrel/goZlDianzqOfferTicketV3/object"
	"github.com/Deansquirrel/goZlDianzqOfferTicketV3/repository"
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
