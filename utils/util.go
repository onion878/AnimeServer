package utils

import (
	"bytes"
	"github.com/satori/go.uuid"
	"math/rand"
	"strings"
	"time"
)

//RandomString(8, "A") 大写
//RandomString(8, "a0") 小写
//RandomString(20, "Aa0") 混合
func RandomString(randLength int, randType string) (result string) {
	var num string = "0123456789"
	var lower string = "abcdefghijklmnopqrstuvwxyz"
	var upper string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := bytes.Buffer{}
	if strings.Contains(randType, "0") {
		b.WriteString(num)
	}
	if strings.Contains(randType, "a") {
		b.WriteString(lower)
	}
	if strings.Contains(randType, "A") {
		b.WriteString(upper)
	}
	var str = b.String()
	var strLen = len(str)
	if strLen == 0 {
		result = ""
		return
	}

	rand.Seed(time.Now().UnixNano())
	b = bytes.Buffer{}
	for i := 0; i < randLength; i++ {
		b.WriteByte(str[rand.Intn(strLen)])
	}
	result = b.String()
	return
}

//获取uuid
func NewKeyId() string {
	if id, err := uuid.NewV4(); err == nil {
		return id.String()
	}
	return "创建失败!"
}