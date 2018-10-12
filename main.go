package main

import (
	"./utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

func main() {
	r := gin.Default()
	utils.StartPool()
	r.GET("/getIndex/:page", func(c *gin.Context) {
		page, _ := strconv.Atoi(c.Param("page"))
		c.JSON(200, utils.GetIndex(page))
	})
	r.GET("/getChapter/:pid", func(c *gin.Context) {
		c.JSON(200, utils.GetChapter(c.Param("pid")))
	})
	r.GET("/searchByName/:name", func(c *gin.Context) {
		c.JSON(200, utils.SearchByName(c.Param("name")))
	})
	r.GET("/getByName/:name", func(c *gin.Context) {
		c.JSON(200, utils.GetByName(c.Param("name")))
	})
	r.Run(":8060")
}
