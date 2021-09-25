package utils

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"github.com/satori/go.uuid"
	"gopkg.in/gomail.v2"
	"log"
	"math/rand"
	"os"
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
	return uuid.NewV4().String()
}

func SendMail(msg string) {
	d := gomail.NewDialer("smtp.qq.com", 587, "genaretor@qq.com", "nbvlluxakyzgebji")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	m := gomail.NewMessage()
	m.SetHeader("From", "genaretor@qq.com")
	m.SetHeader("To", "2419186601@qq.com")
	m.SetHeader("Subject", "资源获取通知!")
	m.SetBody("text/html", "<b>获取通知</b><br><i>"+msg+"</i>!")

	// Send emails using d.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

type AppConfigProperties map[string]string

func ReadPropertiesFile(filename string) (AppConfigProperties, error) {
	config := AppConfigProperties{}

	if len(filename) == 0 {
		return config, nil
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				config[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return config, nil
}
