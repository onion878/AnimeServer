package utils

import (
	"../structs"
)

func GetIndex(page int) []structs.Index {
	var list []structs.Index
	engine := GetCon()
	engine.OrderBy("`index` asc").Limit(20, page*20).Find(&list)
	return list
}

func SearchByName(name string) []structs.Index {
	var list []structs.Index
	engine := GetCon()
	engine.Where("name like concat('%',?,'%')", name).OrderBy("`index` asc").Find(&list)
	return list
}

func GetByName(name string) []structs.Index {
	var list []structs.Index
	engine := GetCon()
	engine.Where("name=?", name).OrderBy("`index` asc").Find(&list)
	return list
}

func GetChapter(pid string) []structs.Chapter {
	var chapters []structs.Chapter
	engine := GetCon()
	engine.Where("pid = ?", pid).OrderBy("name desc").Find(&chapters)
	return chapters
}
