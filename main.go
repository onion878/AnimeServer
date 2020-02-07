package main

import (
	"./service"
	"./structs"
	"./utils"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"github.com/jasonlvhit/gocron"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const path = "https://anime1.me"

var runing = false
var identityKey = "id"

type User struct {
	UserName  string
	FirstName string
	LastName  string
}

func main() {
	r := gin.Default()
	utils.StartPool()
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

	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				UserName: claims["id"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals structs.User
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			if utils.Login(loginVals) {
				userID := loginVals.UserName

				// 设置token
				return &User{
					UserName: userID,
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			//用于判断是否有权限
			if v, ok := data.(*User); ok && v.UserName == "2214839296@qq.com" {
				return true
			}
			//
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"auth":    false,
				"message": "未登录!",
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	r.POST("/login", authMiddleware.LoginHandler)

	r.Use(authMiddleware.MiddlewareFunc())
	{

		r.POST("/auth", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"data": jwt.ExtractClaims(c)["id"],
			})
		})
		h := new(service.HistoryService)
		r.POST("/saveHistory", func(c *gin.Context) {
			var history structs.History
			c.ShouldBindJSON(&history)
			history.UserName = jwt.ExtractClaims(c)["id"].(string)
			c.JSON(200, gin.H{
				"data": h.SaveHistory(history),
			})
		})

		r.POST("/deleteHistory", func(c *gin.Context) {
			var history structs.History
			c.ShouldBindJSON(&history)
			history.UserName = jwt.ExtractClaims(c)["id"].(string)
			c.JSON(200, gin.H{
				"data": h.Delete(history),
			})
		})

		r.POST("/deleteAllHistory", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"data": h.DeleteAll(jwt.ExtractClaims(c)["id"].(string)),
			})
		})

		r.GET("/getHistory/:page", func(c *gin.Context) {
			page, _ := strconv.Atoi(c.Param("page"))
			c.JSON(200, h.GetHistory(page))
		})

		f := new(service.FavoriteService)
		r.POST("/saveFavorite", func(c *gin.Context) {
			var favorite structs.Favorite
			c.ShouldBindJSON(&favorite)
			favorite.UserName = jwt.ExtractClaims(c)["id"].(string)
			c.JSON(200, gin.H{
				"data": f.SaveFavorite(favorite),
			})
		})

		r.POST("/deleteFavorite", func(c *gin.Context) {
			var favorite structs.Favorite
			c.ShouldBindJSON(&favorite)
			favorite.UserName = jwt.ExtractClaims(c)["id"].(string)
			c.JSON(200, gin.H{
				"data": f.Delete(favorite),
			})
		})

		r.POST("/deleteAllFavorite", func(c *gin.Context) {
			userName := jwt.ExtractClaims(c)["id"].(string)
			c.JSON(200, gin.H{
				"data": f.DeleteAll(userName),
			})
		})

		r.GET("/getFavorite/:page", func(c *gin.Context) {
			page, _ := strconv.Atoi(c.Param("page"))
			c.JSON(200, f.GetFavorite(page))
		})
	}
	go func() {
		gocron.Every(1).Second().Do(taskWithParams)
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
	engine := utils.GetCon()
	engine.Exec("DROP TABLES IF EXISTS `index`,chapter")
	engine.CreateTables(new(structs.Cookies))
	engine.CreateTables(new(structs.Index))
	engine.CreateTables(new(structs.Chapter))
	getMenu()
	getAllIndex()
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
			utils.SaveChapter(name, pid, src, i, true)
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
					utils.SaveChapter(name, pid, arr[i].File, num, false)
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
					utils.SaveChapter(name, pid, file, num, false)
				} else {
					getMpdChapter(e.Request.URL.Query().Get("dash"), name, pid, num)
				}
			}
		} else {
			start := strings.Index(data, `,file:"`)
			end := strings.Index(data, `",controls:true`)
			if start > 0 && end > 0 {
				file := data[start+7 : end]
				utils.SaveChapter(name, pid, file, num, false)
			}
		}
	})

	c.OnHTML("body", func(e *colly.HTMLElement) {
		println(e.Text)
	})

	c.OnHTML("video source", func(e *colly.HTMLElement) {
		file := e.Attr("src")
		if len(file) > 0 {
			utils.SaveChapter(name, pid, file, num, false)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("cookie", utils.GetCookie())
		fmt.Println("Visiting", r.URL)
	})
	c.Visit(url)
}

// update time: 2018-12-06
func getMpdChapter(url string, name string, pid string, num int) {
	r, err := http.Get(url)
	if err == nil {
		doc, e := goquery.NewDocumentFromReader(io.Reader(r.Body))
		if e == nil {
			var file = ""
			flag := false
			p := doc.Find("representation")
			p.Each(func(i int, s *goquery.Selection) {
				band, ok := s.Attr("fbqualityclass")
				if ok && strings.Trim(band, "\n") == "hd" {
					flag = true
					file = s.Find("baseurl").Eq(0).Text()
				}
			})
			if !flag {
				file = p.Eq(0).Find("BaseURL").Eq(0).Text()
			}
			utils.SaveChapter(name, pid, file, num, false)
		}
	}

}
