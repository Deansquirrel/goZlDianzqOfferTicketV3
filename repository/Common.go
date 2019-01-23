package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/alexbrainman/odbc"
)

//根据配置获取数据库连接
func GetDbConn(server string, port int, dbName string, user string, pwd string) (*sql.DB, error) {
	connString := fmt.Sprintf("Driver={SQL Server};Server=%s,%d;Database=%s;Uid=%s;Pwd=%s;Network=dbmssocn;", server, port, dbName, user, pwd)
	conn, err := sql.Open("odbc", connString)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	//conn.SetMaxIdleConns(30)
	//conn.SetMaxOpenConns(30)
	//conn.SetConnMaxLifetime(time.Second * 60 * 10)
	return conn, nil
}

//检查数据库连接是否有效
func CheckV(db *sql.DB) bool {
	if db != nil {
		err := db.Ping()
		if err == nil {
			return true
		}
	}
	return false
}
