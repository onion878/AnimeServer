package utils

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"structs"
	"time"
	"xorm.io/core"
)

var connect *xorm.Engine

func StartPool() {
	props, err := ReadPropertiesFile("properties")
	if err != nil {
		fmt.Println("Error while reading properties file")
	}
	engine, err := xorm.NewEngine("mysql", props["username"]+":"+props["password"]+"@tcp("+props["url"]+":"+props["port"]+")/"+props["database"])
	if err != nil {
		fmt.Println(err)
		return
	}
	timeZone, _ := time.LoadLocation("Asia/Shanghai")
	engine.TZLocation = timeZone
	engine.Charset("UTF-8")
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
	engine.CreateTables(new(structs.Cookies))
	engine.CreateTables(new(structs.Index))
	engine.CreateTables(new(structs.Chapter))
	engine.CreateTables(new(structs.User))
	engine.CreateTables(new(structs.History))
	engine.CreateTables(new(structs.Favorite))
	connect = engine
}

func GetCon() *xorm.Engine {
	return connect
}
