package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"strings"
	"time"
)

type netLog struct {
	ID         uint `gorm:"primarykey"`
	Datetime   int
	ClientUA   string
	ClientIP   string
	RequestURL string
	HttpCode   string
	DeviceId   string
}

func main() {
	r := gin.Default()

	// Get System Variable
	dbUser := os.Getenv("logger-dbUser")
	dbPassword := os.Getenv("logger-dbPassword")
	dbHost := os.Getenv("logger-dbHost")
	dbPort := os.Getenv("logger-dbPort")
	dbName := os.Getenv("logger-dbName")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&netLog{})

	r.POST("/debugger", func(c *gin.Context) {
		clientUA := strings.Join(c.Request.Header["User-Agent"], " ")
		clientIP := c.ClientIP()
		createTime := int(time.Now().Unix())
		requestURL, reqValid := c.GetPostForm("requestURL")
		if !reqValid {
			c.Status(403)
			return
		}
		httpCode, codeValid := c.GetPostForm("httpCode")
		if !codeValid {
			c.Status(403)
			return
		}
		DeviceId, msgValid := c.GetPostForm("deviceID")
		if !msgValid {
			c.Status(403)
			return
		}

		thisLog := netLog{
			ClientUA: clientUA, ClientIP: clientIP, RequestURL: requestURL,
			HttpCode: httpCode, DeviceId: DeviceId, Datetime: createTime,
		}

		db.Create(&thisLog)
		return

	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
