package utils

import (
	"fmt"
	"strings"
	"structs"
	"time"
)

func Login(user structs.User) bool {
	engine := GetCon()
	var users structs.User
	engine.Id(user.UserName).Get(&users)
	if strings.EqualFold(user.Password, users.Password) && len(user.Password) > 0 {
		return true
	}
	return false
}

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

func FindByName(name string) []structs.Index {
	var list []structs.Index
	engine := GetCon()
	engine.Where("name = ?", name).Find(&list)
	return list
}

func GetNotChapter() []string {
	var list []string
	engine := GetCon()
	res, err := engine.QueryString("select a.name from `index` a left join chapter b on a.id=b.pid where b.id is null")
	if err == nil {
		for i := range res {
			list = append(list, res[i]["name"])
		}
	}
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

func GetCookie() string {
	var cookie []structs.Cookies
	engine := GetCon()
	engine.Where("id = 1").Find(&cookie)
	return cookie[0].Value
}

func SaveIndex(name string, chapter string, url string, order int, image string) structs.Index {
	var index structs.Index
	engine := GetCon()
	index.Url = url
	index.Id = NewKeyId()
	index.Update = time.Now()
	index.Total = 0
	index.Name = name
	index.Chapter = chapter
	index.Index = order
	index.Image = image
	index.UpdateFlag = true
	engine.Insert(&index)
	fmt.Printf("%+v\n", index)
	return index
}

func GetAllIndex() []structs.Index {
	var indexs []structs.Index
	engine := GetCon()
	engine.Where("flag = 0").OrderBy("`index` asc").Find(&indexs)
	return indexs
}

func FindIndexByUpadate(flag bool) []structs.Index {
	var indexs []structs.Index
	engine := GetCon()
	engine.Where("flag = 0").Where("`update_flag`=?", flag).OrderBy("`index` asc").Find(&indexs)
	return indexs
}

func SaveOrUpdateIndex(name string, chapter string, url string, order int, image string, flag bool) structs.Index {
	var index structs.Index
	var indexs []structs.Index
	engine := GetCon()
	engine.Where("name = ?", name).Find(&indexs)
	if len(indexs) == 0 {
		index.Id = NewKeyId()
		index.Update = time.Now()
		index.Total = 0
		index.Url = url
		index.Name = name
		index.Chapter = chapter
		index.Image = image
		index.Index = order
		index.UpdateFlag = flag
		engine.Insert(&index)
		index.Flag = true
	} else {
		index = indexs[0]
		index.Chapter = chapter
		index.Url = url
		index.Update = time.Now()
		index.Index = order
		index.Image = image
		index.UpdateFlag = flag
		engine.ID(index.Id).Cols("`update_flag`", "`chapter`", "`url`", "`update`", "`index`").Update(&index)
		index.Flag = false
	}
	fmt.Printf("%+v\n", index)
	return index
}

func UpdateIndexInfo(id string, image string, label string, info string, date string) {
	var index structs.Index
	engine := GetCon()
	index.Id = id
	index.Image = image
	index.Label = label
	index.Info = info
	index.Date = date
	engine.ID(index.Id).Cols("`info`", "`date`", "`label`").Update(&index)
}

func UpdateIndexOrder(id string, order int) {
	var index structs.Index
	engine := GetCon()
	index.Update = time.Now()
	index.Index = order
	engine.ID(id).Update(&index)
	fmt.Printf("%+v\n", index)
}

func UpdateIndexFlag(id string, flag bool) {
	var index structs.Index
	engine := GetCon()
	index.Update = time.Now()
	index.UpdateFlag = flag
	engine.ID(id).Cols("`update_flag`").Update(&index)
	fmt.Printf("%+v\n", index)
}

func SaveChapter(name string, pid string, url string, num int, webFlag bool) structs.Chapter {
	var chapter structs.Chapter
	var chapters []structs.Chapter
	engine := GetCon()
	engine.Where("name = ?", name).Where("pid = ?", pid).Find(&chapters)
	if len(chapters) == 0 {
		chapter.Id = NewKeyId()
		chapter.Pid = pid
		chapter.Name = name
		chapter.Path = url
		chapter.Num = num
		chapter.Flag = webFlag
		engine.Insert(&chapter)
	} else {
		chapter = chapters[0]
		chapter.Path = url
		chapter.Flag = webFlag
		engine.Id(chapters[0].Id).Cols("flag", "path").Update(&chapter)
	}
	return chapter
}

func DeleteByName(name string) {
	var list []structs.Index
	engine := GetCon()
	engine.Where("name = ?", name).Find(&list)
	for i := 0; i < len(list); i++ {
		engine.Exec("DELETE FROM `chapter` WHERE `pid`=?", list[i].Id)
	}
}

func JudgeNew(newRow structs.Index) bool {
	var indexs []structs.Index
	var index structs.Index
	engine := GetCon()
	engine.Where("`index` = 0").Find(&indexs)
	if len(indexs) == 1 {
		index = indexs[0]
	}
	if index.Name == newRow.Name && index.Chapter == newRow.Chapter {
		return false
	} else {
		return true
	}
}
