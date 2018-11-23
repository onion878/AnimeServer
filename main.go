package main

import (
	"./structs"
	"./utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"github.com/jasonlvhit/gocron"
	"strconv"
	"strings"
)

const path = "https://anime1.me"

var runing = false

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
	r.GET("/getAllSource", func(c *gin.Context) {
		if !runing {
			go func() {
				getAllSource()
			}()
			c.JSON(200, gin.H{
				"success": true,
				"msg":     "获取所有资源中...",
			})
		} else {
			c.JSON(200, gin.H{
				"success": true,
				"msg":     "正在获取所有资源中,请稍等一下!",
			})
		}
	})
	r.GET("/getOneSource/:name", func(c *gin.Context) {
		name := c.Param("name")
		utils.DeleteByName(name)
		getOneSource(name)
		println("return true")
		c.JSON(200, gin.H{
			"success": true,
			"msg":     "获取成功!",
		})
	})
	go func() {
		gocron.Every(1).Second().Do(taskWithParams, 1, "hello")
		<-gocron.Start()
	}()
	r.Run(":8060")
}

func taskWithParams(a int, b string) {
	if !runing {
		getFirstMenu()
	}
}

func getFirstMenu() {
	c := colly.NewCollector()
	flag := false
	c.OnHTML(".entry-content table tbody tr", func(e *colly.HTMLElement) {
		if !flag {
			name := e.DOM.Find(".column-1 a").Text()
			chapter := e.DOM.Find(".column-2").Text()
			var newRow = new(structs.Index)
			newRow.Name = name
			newRow.Chapter = chapter
			if utils.JudgeNew(*newRow) {
				getAllSource()
			}
		}
		flag = true
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("cookie", utils.GetCookie())
		fmt.Println("Visiting", r.URL)
	})
	c.Visit(path)
}

func getAllSource() {
	runing = true
	go func() {
		utils.SendMail("开始获取资源!")
	}()
	engine := utils.GetCon()
	engine.Exec("DROP TABLES IF EXISTS `index`,chapter")
	engine.CreateTables(new(structs.Cookies))
	engine.CreateTables(new(structs.Index))
	engine.CreateTables(new(structs.Chapter))
	getMenu()
	getAllIndex()
	go func() {
		utils.SendMail("获取资源完成!")
	}()
	runing = false
}

func getOneSource(n string) {
	c := colly.NewCollector()
	println("获取所有目录")
	c.OnHTML(".entry-content table tbody tr", func(e *colly.HTMLElement) {
		href, _ := e.DOM.Find(".column-1 a").Attr("href")
		name := e.DOM.Find(".column-1 a").Text()
		chapter := e.DOM.Find(".column-2").Text()
		if strings.EqualFold(n, name) {
			index := utils.SaveOrUpdateIndex(name, chapter, href, e.Index)
			getChapter(path+href, index.Id)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("cookie", utils.GetCookie())
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(path)
}

func getMenu() {
	c := colly.NewCollector()
	println("获取所有目录")
	c.OnHTML(".entry-content table tbody tr", func(e *colly.HTMLElement) {
		href, _ := e.DOM.Find(".column-1 a").Attr("href")
		name := e.DOM.Find(".column-1 a").Text()
		chapter := e.DOM.Find(".column-2").Text()
		utils.SaveIndex(name, chapter, href, e.Index)
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("cookie", utils.GetCookie())
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(path)
}

func getAllIndex() {
	index := utils.GetAllIndex()
	println("获取详情...")
	for i := 0; i < len(index); i++ {
		data := index[i]
		getChapter(path+data.Url, data.Id)
	}
}

func getChapter(url string, pid string) {
	c := colly.NewCollector()
	// Find and visit all links
	c.OnHTML("main", func(e *colly.HTMLElement) {
		s := e.DOM.Find("iframe[src]")
		d := e.DOM.Find(".entry-title a[href]")
		for i := 0; i < s.Length(); i++ {
			src, _ := s.Eq(i).Attr("src")
			name := d.Eq(i).Text()
			getChapterUrl(src, name, pid, i)
		}
	})

	c.OnHTML(".nav-previous a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("cookie", utils.GetCookie())
		fmt.Println("Visiting", r.URL)
	})
	c.Visit(url)
}

func getChapterUrl(url string, name string, pid string, num int) {
	c := colly.NewCollector(colly.Async(true))
	// Find and visit all links
	c.OnHTML("body script", func(e *colly.HTMLElement) {
		data := e.Text
		start := strings.Index(data, "sources:")
		end := strings.Index(data, ",controls:true")
		if start > 0 && end > 0 {
			s := data[start+8 : end]
			s = strings.Replace(s, ",label:", `,"label":`, -1)
			var arr []structs.UrlData
			_ = json.Unmarshal([]byte(s), &arr)
			var flag = false
			for i := 0; i < len(arr); i++ {
				if arr[i].Default == "true" {
					utils.SaveChapter(name, pid, arr[i].File, num)
					flag = true
				}
			}
			if !flag {
				some := 0
				file := ""
				for i := 0; i < len(arr); i++ {
					s := arr[i].Label
					if len(s) > 0 {
						hd, _ := strconv.Atoi(s[0 : len(s)-1])
						if hd > some {
							file = arr[i].File
						}
						some = hd
					}
				}
				if len(file) > 0 {
					utils.SaveChapter(name, pid, file, num)
				}
			}
		} else {
			start := strings.Index(data, `,file:"`)
			end := strings.Index(data, `",controls:true`)
			if start > 0 && end > 0 {
				file := data[start+7 : end]
				utils.SaveChapter(name, pid, file, num)
			}
		}
	})

	c.OnHTML("video source", func(e *colly.HTMLElement) {
		file := e.Attr("src")
		if len(file) > 0 {
			utils.SaveChapter(name, pid, file, num)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("cookie", utils.GetCookie())
		fmt.Println("Visiting", r.URL)
	})
	c.Visit(url)
}
