package utils

import (
	"../structs"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

var connect *xorm.Engine

func StartPool() {
	engine, err := xorm.NewEngine("mysql", "onion:1234@tcp(localhost:3306)/anime?charset=utf8")
	if err != nil {
		fmt.Println(err)
		return
	}
	//连接测试
	if err := engine.Ping(); err != nil {
		fmt.Println(err)
		return
	}
	//日志打印SQL
	engine.ShowSQL(true)
	//设置连接池的空闲数大小
	engine.SetMaxIdleConns(5)
	//设置最大打开连接数
	engine.SetMaxOpenConns(5)
	//名称映射规则主要负责结构体名称到表名和结构体field到表字段的名称映射
	engine.SetTableMapper(core.SnakeMapper{})
	engine.CreateTables(new(structs.Index))
	engine.CreateTables(new(structs.Chapter))
	connect = engine
}

func GetCon() *xorm.Engine {
	return connect
}
