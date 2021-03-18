package service

import (
	"../structs"
	"../utils"
)

type HistoryService struct {
}

func (*HistoryService) SaveHistory(history structs.History) bool {
	engine := utils.GetCon()
	history.Id = utils.NewKeyId()
	engine.Exec("DELETE FROM history WHERE user_name=? and `index`=?", history.UserName, history.Index)
	total, err := engine.Insert(&history)
	if err != nil || total == 0 {
		return false
	} else {
		return true
	}
}

func (*HistoryService) GetHistory(page int) []structs.History {
	var list []structs.History
	engine := utils.GetCon()
	engine.OrderBy("`created` desc").Limit(20, page*20).Find(&list)
	return list
}

func (*HistoryService) Delete(history structs.History) bool {
	engine := utils.GetCon()
	engine.Exec("DELETE FROM history WHERE user_name=? and `index`=?", history.UserName, history.Index)
	return true
}

func (*HistoryService) DeleteAll(userName string) bool {
	engine := utils.GetCon()
	engine.Exec("DELETE FROM history WHERE user_name=?", userName)
	return true
}
