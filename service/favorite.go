package service

import (
	"structs"
	"utils"
)

type FavoriteService struct {
}

func (*FavoriteService) SaveFavorite(f structs.Favorite) bool {
	engine := utils.GetCon()
	total, _ := engine.Count(&f)
	if total == 0 {
		f.Id = utils.NewKeyId()
		t, err := engine.Insert(&f)
		if err != nil || t == 0 {
			return false
		} else {
			return true
		}
	} else {
		engine.Exec("DELETE FROM favorite WHERE user_name=? and `index`=?", f.UserName, f.Index)
		return true
	}
}

func (*FavoriteService) GetFavorite(page int) []structs.Favorite {
	var list []structs.Favorite
	engine := utils.GetCon()
	engine.OrderBy("`created` desc").Limit(20, page*20).Find(&list)
	return list
}

func (*FavoriteService) DeleteAll(userName string) bool {
	engine := utils.GetCon()
	engine.Exec("DELETE FROM favorite WHERE user_name=?", userName)
	return true
}

func (*FavoriteService) Delete(history structs.Favorite) bool {
	engine := utils.GetCon()
	engine.Exec("DELETE FROM favorite WHERE user_name=? and `index`=?", history.UserName, history.Index)
	return true
}
