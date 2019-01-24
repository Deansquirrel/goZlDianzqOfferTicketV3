package global

import (
	"database/sql"
	"github.com/Deansquirrel/go-tool"
	"github.com/Deansquirrel/goZlDianzqOfferTicketV3/object"
)

var SysConfig *object.SysConfig
var PeiZhDbConn *sql.DB
var Redis *go_tool.MyRedis
var RabbitMQ *go_tool.MyRabbitMQ

var HxDbConnConfigMap map[string][]string
var HxDbConnMap map[string]*sql.DB

var RedisDbId1 int
var RedisDbId2 int

var SnoServer string
var SnoWorkerId int
