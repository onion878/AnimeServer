package structs

import "time"

type Index struct {
	Id         string    `xorm:"not null pk VARCHAR(40)"`
	Name       string    `xorm:"not null VARCHAR(100)"`
	Image      string    `xorm:"not null VARCHAR(100)"`
	Chapter    string    `xorm:"text"`
	Total      int       `xorm:"not null int"`
	Update     time.Time `xorm:"TIMESTAMP"`
	Index      int       `xorm:"null int"`
	Url        string    `xorm:"null VARCHAR(100)"`
	Info       string    `xorm:"text"`
	Date       string    `xorm:"null VARCHAR(20)"`
	Label      string    `xorm:"null VARCHAR(200)"`
	Created    time.Time `xorm:"TIMESTAMP created"`
	UpdateFlag bool      `xorm:"bit"`
	Flag       bool
}

type Chapter struct {
	Id      string    `xorm:"not null pk VARCHAR(40)"`
	Pid     string    `xorm:"not null VARCHAR(40)"`
	Name    string    `xorm:"not null VARCHAR(150)"`
	Path    string    `xorm:"text"`
	Flag    bool      `xorm:"bit"`
	Num     int       `xorm:"not null int"`
	Created time.Time `xorm:"TIMESTAMP created"`
}

type UrlData struct {
	Type    string `json:"type"`
	File    string `json:"file"`
	Label   string `json:"label"`
	Default string `json:"default"`
}

type Page struct {
	Page  int    `json:"page"`
	Start int    `json:"start"`
	Limit int    `json:"limit"`
	Name  string `json:"name"`
}

type Cookies struct {
	Id    string `xorm:"not null pk int"`
	Value string `xorm:"text"`
}

type User struct {
	UserName string `xorm:"not null pk VARCHAR(40)"`
	Password string `not null VARCHAR(40)`
}

type History struct {
	Id       string    `xorm:"not null pk VARCHAR(40)"`
	Index    string    `xorm:"not null VARCHAR(100)"`
	Chapter  int64     `xorm:"not null int"`
	Duration int64     `xorm:"not null int"`
	UserName string    `xorm:"not null VARCHAR(40)"`
	Created  time.Time `xorm:"TIMESTAMP created"`
}

type Favorite struct {
	Id       string    `xorm:"not null pk VARCHAR(40)"`
	Index    string    `xorm:"not null VARCHAR(100)"`
	UserName string    `xorm:"not null VARCHAR(40)"`
	Created  time.Time `xorm:"TIMESTAMP created"`
}
