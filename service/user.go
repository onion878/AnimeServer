package service

import (
	"../structs"
	"../utils"
	"github.com/gin-gonic/gin"
)

type User struct {
}

func (*User) ListByPage(page structs.Page) interface{} {
	engine := utils.GetCon()
	var user structs.User
	var users []structs.User
	total, _ := engine.Count(&user)
	engine.Limit(page.Limit, page.Start).OrderBy("`user_name` asc").Find(&users)
	return gin.H{
		"total": total,
		"rows":  users,
	}
}

func (*User) Update(user structs.User) bool {
	engine := utils.GetCon()
	_, err := engine.ID(user.UserName).Update(&user)
	if err == nil {
		return true
	} else {
		return false
	}
}

func (*User) Create(user structs.User) bool {
	engine := utils.GetCon()
	total, err := engine.Insert(&user)
	if err != nil || total == 0 {
		return false
	} else {
		return true
	}
}

func (*User) Remove(user structs.User) bool {
	engine := utils.GetCon()
	total, err := engine.Delete(&user)
	if err != nil || total == 0 {
		return false
	} else {
		return true
	}
}
