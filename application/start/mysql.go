package start

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
	"time"
)

type MysqlConfig struct {
	DbId int `json:"db_id"`
	UserName string `json:"userName"`
	Password string `json:"password"`
	Network string `json:"network"`
	Host string `json:"host"`
	Port int `json:"port"`
	Database string `json:"database"`
	ConnMaxLifetime int `json:"conn_max_lifetime"`
	MaxOpenConns int `json:"max_open_conns"`
	MaxIdleConns int `json:"max_idle_conns"`
}

var defaultMysqlConfig *MysqlConfig = &MysqlConfig{
	1,
	"root",
	"123456",
	"tcp",
	"localhost",
	23306,
	"*",
	100,
	64,
	8,
}

/**
 	定义一个连接数组
 */
var defaultDbsNum = 32

var Dbs = make(map[string]*sqlx.DB, defaultDbsNum)

func GetMysqlConn(database string, dbId int) (db *sqlx.DB,err error) {
	db,err = nil, nil
	errStr := "not found conn"
	if dbId > 0 {
		key := genKey(dbId, database)
		if v, ok := Dbs[key]; ok {
			db = v
		} else {
			err = errors.New(errStr)
		}
	} else {
		for dbKey, dbs := range Dbs {
			if strings.HasPrefix(dbKey, database) {
				db = dbs
				break
			}
		}

		if db == nil {
			err = errors.New(errStr)
		}
	}
	return
}

func InitMysql(mysqlcfgs ...*MysqlConfig)  {
	if len(mysqlcfgs) == 0 {
		mysqlcfgs = []*MysqlConfig{defaultMysqlConfig}
	}
	for _, dbcfg := range mysqlcfgs {
		key := genKey(dbcfg.DbId, dbcfg.Database)
		dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s",dbcfg.UserName,dbcfg.Password,dbcfg.Network,dbcfg.Host,dbcfg.Port,dbcfg.Database)
		dsn = dsn + "?parseTime=true&loc=Asia%2FShanghai&charset=utf8"
		DB,err := sqlx.Open("mysql",dsn)
		if err != nil{
			errorStr := fmt.Sprintf("Open mysql failed,err:%v\n", err)
			Logger.Error(errorStr)
			panic(errors.New(errorStr))
		}
		DB.SetConnMaxLifetime(time.Duration(dbcfg.ConnMaxLifetime) * time.Second)  //最大连接周期，超过时间的连接就close
		DB.SetMaxOpenConns(dbcfg.MaxOpenConns) //设置最大连接数
		DB.SetMaxIdleConns(dbcfg.MaxIdleConns) //设置闲置连接数
		Dbs[key] = DB
	}
}

func genKey(id int, database string) string {
	return database + strconv.Itoa(id)
}
