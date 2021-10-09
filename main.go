package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"github.com/jasonlvhit/gocron"
	"strconv"
	"strings"
	"structs"
	"utils"
)

var path = ""

var runing = false
var runDetail = false
var identityKey = "id"

func main() {
	r := gin.Default()
	utils.StartPool()
	props, _ := utils.ReadPropertiesFile("properties")
	path = props["source"]
	r.Use(CORSMiddleware())
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
				engine := utils.GetCon()
				engine.Exec("delete from `chapter`")
				engine.Exec("delete from `index`")
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
		c.JSON(200, gin.H{
			"success": true,
			"msg":     "获取成功!",
		})
	})

	go func() {
		getAllSource()
	}()

	go func() {
		gocron.Every(30).Seconds().Do(taskWithParams)
		<-gocron.Start()
	}()

	go func() {
		gocron.Every(600).Seconds().Do(checkChapter)
		<-gocron.Start()
	}()
	r.Run(":8060")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func taskWithParams() {
	if !runing {
		getFirstMenu()
	}
}

func checkChapter() {
	if !runDetail && !runing {
		runDetail = true
		list := utils.GetNotChapter()
		for i := range list {
			getOneSource(list[i])
		}
		runDetail = false
	}
}

func getFirstMenu() {
	c := colly.NewCollector()
	flag := false
	c.OnHTML(".area .topli ul li", func(e *colly.HTMLElement) {
		if !flag {
			area := e.DOM.Find("span a").Text()
			if area == "日本" {
				name := e.DOM.Find("a").Eq(1).Text()
				chapter := e.DOM.Find("b a").Text()
				var newRow = new(structs.Index)
				newRow.Name = name
				newRow.Chapter = chapter
				if !runing && utils.JudgeNew(*newRow) {
					println("获取更新")
					getNew()
				}
			}
		}
		flag = true
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("cookie", utils.GetCookie())
		fmt.Println("Visiting", r.URL)
	})
	c.Visit(path + "/new")
}

func getNew() {
	runing = true
	c := colly.NewCollector()
	var data []structs.Index
	i := 0
	c.OnHTML(".area .topli ul li", func(e *colly.HTMLElement) {
		area := e.DOM.Find("span a").Text()
		if area == "日本" {
			name := e.DOM.Find("a").Eq(1).Text()
			href, _ := e.DOM.Find("a").Eq(1).Attr("href")
			chapter := e.DOM.Find("b a").Text()
			list := utils.GetByName(name)
			println("更新资源:" + name)
			println(e.Index)
			if list != nil && len(list) > 0 {
				println("更新资源:" + name)
				data = append(data, structs.Index{Id: list[0].Id, Name: name, Chapter: chapter, Url: href, Index: i, UpdateFlag: true, Image: list[0].Image})
				i++
			} else {
				data = append(data, structs.Index{Id: "", Name: name, Chapter: chapter, Url: href, Index: i, UpdateFlag: true})
				i++
			}
		}
		if e.Index == 99 {
			if data != nil && len(data) > 0 {
				changeIndex(data)
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("cookie", utils.GetCookie())
		fmt.Println("Visiting", r.URL)
	})
	c.Visit(path + "/new")
}

func changeIndex(data []structs.Index) {
	for i := range data {
		name := data[i].Name
		chapter := data[i].Chapter
		href := data[i].Url
		image := data[i].Image
		println(data[i].Id)
		if data[i].Id != "" {
			utils.DeleteByName(name)
			utils.SaveOrUpdateIndex(name, chapter, href, i, image, true)
			getOneSource(name)
		} else {
			utils.SaveIndex(name, chapter, href, i, image)
			getOneSource(name)
		}
	}
	newList := utils.FindIndexByUpadate(true)
	oldList := utils.FindIndexByUpadate(false)
	newList = append(newList, oldList...)
	for i := range newList {
		utils.UpdateIndexOrder(newList[i].Id, i)
		utils.UpdateIndexFlag(newList[i].Id, false)
	}
	runing = false
}

func getAllSource() {
	runing = true
	getMenu()
	getAllIndex()
	runing = false
}

func getOneSource(n string) {
	println("重新获取:" + n)
	list := utils.FindByName(n)
	if len(list) > 0 {
		getChapter(path+list[0].Url, list[0].Id)
	}
}

func getMenu() {
	c := colly.NewCollector()
	println("获取所有目录")
	c.OnHTML(".area .fire .lpic ul li", func(e *colly.HTMLElement) {
		href, _ := e.DOM.Find("h2 a").Attr("href")
		p := e.Request.URL.Path
		name := e.DOM.Find("h2 a").Text()
		chapter := e.DOM.Find("span font").Text()
		image, _ := e.DOM.Find("a img").Attr("src")
		i := utils.GetIntFromString(p) - 1
		list := utils.FindByName(name)
		if len(list) > 0 {
			if list[0].Chapter != chapter {
				utils.SaveOrUpdateIndex(name, chapter, href, i*15+e.Index, image, true)
			}
		} else {
			utils.SaveIndex(name, chapter, href, i*15+e.Index, image)
		}
	})

	c.OnHTML(".area .fire .pages a", func(e *colly.HTMLElement) {
		href, _ := e.DOM.Attr("href")
		name := e.DOM.Text()
		if name == "下一页" {
			c.Visit(path + href)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("cookie", utils.GetCookie())
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(path + "/japan")
}

func getAllIndex() {
	index := utils.GetAllIndex()
	println("获取详情...")
	for i := 0; i < len(index); i++ {
		data := index[i]
		list := utils.FindByName(data.Name)
		if len(list) > 0 {
			if list[0].UpdateFlag {
				utils.UpdateIndexFlag(data.Id, false)
				getChapter(path+data.Url, data.Id)
			}
		}
	}
}

func getChapter(url string, pid string) {
	c := colly.NewCollector()
	// Find and visit all links
	c.OnHTML(".area .fire .tabs .main0 .movurl ul", func(e *colly.HTMLElement) {
		s := e.DOM.Find("li")
		for i := 0; i < s.Length(); i++ {
			src, _ := s.Eq(i).Find("a").Attr("href")
			name := s.Eq(i).Find("a").Text()
			utils.SaveChapter(name, pid, src, i, true)
			getChapterUrl(src, name, pid, i)
		}
	})

	c.OnHTML(".area .fire", func(e *colly.HTMLElement) {
		labels := e.DOM.Find(".sinfo span").Eq(2).Find("a")
		var rows []string
		for i := 0; i < labels.Length(); i++ {
			labels.Eq(i).Text()
			rows = append(rows, labels.Eq(i).Text())
		}
		label := strings.Join(rows, ",")
		image, _ := e.DOM.Find(".thumb img").Attr("src")
		info := e.DOM.Find(".info").Text()
		date := e.DOM.Find(".sinfo span").Eq(0).Text()
		utils.UpdateIndexInfo(pid, image, label, strings.Trim(strings.Replace(info, "\n", "", -1), " "), strings.Replace(date, "上映:", "", -1))
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
	c.OnHTML(".play .area .bofang", func(e *colly.HTMLElement) {
		d, _ := e.DOM.Find("div").Attr("data-vid")
		p := strings.Replace(d, "$mp4", "", -1)
		webFlag := strings.HasSuffix(p, ".m3u8") || strings.HasSuffix(p, ".mp4")
		utils.SaveChapter(name, pid, p, num, !webFlag)
		c.Visit(p)
	})

	c.OnHTML("body script", func(e *colly.HTMLElement) {
		data := e.Text
		start := strings.Index(data, "url:")
		end := strings.Index(data, ",\n                pic")
		if start > 0 && end > 0 {
			file := data[start+6 : end-1]
			utils.SaveChapter(name, pid, file, num, false)
		} else {
			// 未发现不处理
		}
	})

	c.OnHTML("body", func(e *colly.HTMLElement) {
		//println(e.Text)
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("cookie", utils.GetCookie())
		fmt.Println("Visiting", r.URL)
	})
	c.Visit(path + url)
}
