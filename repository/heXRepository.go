package repository

import (
	"bytes"
	"database/sql"
	"errors"
	"github.com/Deansquirrel/goZlDianzqOfferTicketV3/common"
	"github.com/Deansquirrel/goZlDianzqOfferTicketV3/object"
)

type HeXRepository struct {
}

func (hx *HeXRepository) CreateLittleTktCreate(conn *sql.DB, tktInfo []object.TktInfo) error {

	if conn == nil {
		return errors.New("数据库连接不能为空")
	}
	if tktInfo == nil || len(tktInfo) < 1 {
		return errors.New("传入列表不能为空")
	}

	stmt, err := conn.Prepare(getQueryString(len(tktInfo)))
	if err != nil {
		return err
	}
	defer func() {
		errLs := stmt.Close()
		if errLs != nil {
			common.PrintOrLog(errLs.Error())
		}
	}()

	var c = make([]interface{}, 0)
	var val object.TktInfo
	for i := 0; i < len(tktInfo); i++ {
		val = tktInfo[i]
		c = append(c, val.AppId)
		c = append(c, val.AccId)
		c = append(c, val.TktNo)
		c = append(c, val.CashMy)
		c = append(c, val.AddMy)
		c = append(c, val.TktName)
		c = append(c, val.TktKind)
		c = append(c, val.PCno)
		c = append(c, val.EffDate)
		c = append(c, val.Deadline)
		c = append(c, val.CrYwLsh)
		c = append(c, val.CrBr)
	}

	_, err = stmt.Exec(c...)
	if err != nil {
		return err
	}

	return nil
}

////解析配置字符串,并获取连接
//func (hx *HeXRepository) GetDbConnByString(s string) (*sql.DB, error) {
//	if global.HxDbConnMap == nil {
//		global.HxDbConnMap = make(map[string]*sql.DB)
//	}
//	var conn *sql.DB
//	if _, ok := global.HxDbConnMap[s]; ok {
//		conn = global.HxDbConnMap[s]
//		err := conn.Ping()
//		if err != nil {
//			delete(global.HxDbConnMap, s)
//			return hx.getNewConn(s)
//		} else {
//			return conn, nil
//		}
//	}
//	return hx.getNewConn(s)
//}
//
//func (hx *HeXRepository) getNewConn(s string) (*sql.DB, error) {
//	conn, err := hx.getDbConnByString(s)
//	if err != nil {
//		return nil, err
//	} else {
//		global.HxDbConnMap[s] = conn
//		return conn, nil
//	}
//}
//
//func (hx *HeXRepository) getDbConnByString(s string) (*sql.DB, error) {
//	config := strings.Split(s, "|")
//	if len(config) != 5 {
//		common.PrintAndLog("数据库配置串解析失败 - " + s)
//		return nil, errors.New("数据库配置串解析失败")
//	}
//	port, err := strconv.Atoi(config[1])
//	if err != nil {
//		common.PrintAndLog("数据库配置串端口解析失败 - " + s)
//		return nil, errors.New("数据库配置串端口解析失败")
//	}
//	return GetDbConn(config[0], port, config[2], config[3], config[4])
//}

func getQueryString(n int) string {
	if n > 0 {
		var buffer bytes.Buffer
		buffer.WriteString(getCreateTempTableTktInfoSqlStr())
		for i := 0; i < n; i++ {
			buffer.WriteString(getInsertTempTableTktInfoSqlStr())
		}
		buffer.WriteString(getExecProc())
		buffer.WriteString(getDropTempTableTktInfoSqlStr())
		return buffer.String()
	} else {
		return ""
	}
}

func getCreateTempTableTktInfoSqlStr() string {

	var buffer bytes.Buffer
	buffer.WriteString("CREATE TABLE #TktInfo")
	buffer.WriteString("(")
	buffer.WriteString("    Appid varchar(30),")
	buffer.WriteString("    Accid bigint,")
	buffer.WriteString("    Tktno varchar(30),")
	buffer.WriteString("    Cashmy decimal(18,2),")
	buffer.WriteString("    Addmy decimal(18,2),")
	buffer.WriteString("    Tktname nvarchar(30),")
	buffer.WriteString("    TktKind	varchar(30),")
	buffer.WriteString("    Pcno varchar(30),")
	buffer.WriteString("    EffDate smalldatetime,")
	buffer.WriteString("    Deadline smalldatetime,")
	buffer.WriteString("    CrYwlsh varchar(12),")
	buffer.WriteString("    CrBr varchar(30)")
	buffer.WriteString(") ")
	return buffer.String()
}

func getInsertTempTableTktInfoSqlStr() string {
	var buffer bytes.Buffer
	buffer.WriteString("insert into #TktInfo(Appid,Accid,Tktno,Cashmy,Addmy,Tktname,TktKind,Pcno,EffDate,Deadline,CrYwlsh,CrBr) ")
	buffer.WriteString("select ?,?,?,?,?,?,?,?,?,?,?,? ")
	return buffer.String()
}

func getExecProc() string {
	sqlStr := "exec pr_CreateLittleTkt_Create "
	return sqlStr
}

func getDropTempTableTktInfoSqlStr() string {
	sqlStr := "Drop table #TktInfo "
	return sqlStr
}
