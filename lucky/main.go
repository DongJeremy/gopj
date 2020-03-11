package main

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

const (
	// MyUser database user name
	MyUser = "cloud"
	// Password database user password
	Password = "passwd"
	// Host database host ip
	Host = "192.168.1.61"
	// MPort database host port
	MPort = 3306
	// DbName database db name
	DbName = "mydb"
	// SUCCESS return value
	SUCCESS = "success"
	// Port listen port
	Port = "8080"
)

// Person model
type Person struct {
	ID         int64  `json:"id" gorm:"auto-increment"`
	Name       string `json:"name"`
	Image      string `json:"image"`
	ThumbImage string `json:"thumb_image" gorm:"column:thumb_image"`
	Lucky      bool
	Level      int64
}

// PersonGroup all persons
type PersonGroup struct {
	PersonList *[]Person
}

var personGroup PersonGroup

var (
	errInitPerson = errors.New("初始化错误")
)

func getRandom(num int) int {
	if num < 0 {
		num = 1
	}
	return rand.Intn(num)
}

func (p *PersonGroup) setLuckyNumber(luckySize, luckylevel int64) ([]Person, error) {
	rand.Seed(time.Now().UnixNano())
	var ps []Person
	personArr := p.PersonList
	if personArr == nil {
		return nil, errInitPerson
	}
	size := len(*personArr)
	var i int64
	for i < luckySize {
		r1 := getRandom(size)
		for (*personArr)[r1].Lucky {
			r1 = getRandom(size)
		}
		(*personArr)[r1].Lucky = true
		(*personArr)[r1].Level = luckylevel
		i++
		ps = append(ps, (*personArr)[r1])
	}
	return ps, nil
}

func (p *PersonGroup) getAvailableNumber() int {
	var sum int
	pl := p.PersonList
	for _, p := range *pl {
		if !p.Lucky {
			sum++
		}
	}
	return sum
}

func mysqlConnectString(driver, host string) (db *gorm.DB, err error) {
	connArgs := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", MyUser, Password, host, MPort, DbName)
	db, err = gorm.Open(driver, connArgs)
	return
}

func sqliteConnectString(driver, host string) (db *gorm.DB, err error) {
	db, err = gorm.Open(driver, host)
	return
}

func dbConn(host string) (db *gorm.DB) {
	db, err := mysqlConnectString("mysql", host)
	//db, err := sqliteConnectString("sqlite3", "gorm.db")

	if err != nil {
		panic(err.Error())
	}
	db.Table("t_person").AutoMigrate(&Person{})
	return db
}

// GetTotalNum get lucky number by random
func GetTotalNum(c *gin.Context) {
	count := personGroup.getAvailableNumber()
	if count < 0 {
		c.AbortWithStatus(http.StatusInternalServerError)
	} else {
		c.JSON(http.StatusOK, count)
	}
}

// GetAll get all persons
func GetAll(c *gin.Context) {
	if personGroup.PersonList != nil {
		c.JSON(http.StatusOK, personGroup.PersonList)
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

// GetAllPersonFromDB 获取所有人员信息
func GetAllPersonFromDB(db *gorm.DB) ([]Person, error) {
	var personArr []Person
	if err := db.Table("t_person").Find(&personArr).Error; err != nil {
		return nil, err
	}
	return personArr, nil
}

// GetLuckyNum 获取remain人数信息
func GetLuckyNum(c *gin.Context) {
	luckySize := c.DefaultQuery("luckyNum", "0")
	luckyPrize := c.DefaultQuery("luckyPrize", "0")
	size, err := strconv.ParseInt(luckySize, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	level, err := strconv.ParseInt(luckyPrize, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t, err := personGroup.setLuckyNumber(size, level)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)
}

func main() {
	var host string
	if len(os.Args) > 1 {
		host = os.Args[1]
	} else {
		host = Host
	}

	db := dbConn(host)
	defer db.Close()

	gin.SetMode(gin.ReleaseMode)
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	router := gin.Default()

	//setResource(router)
	router.LoadHTMLGlob("view/*")
	router.Static("/vender/static", "./static")

	router.GET("/", func(c *gin.Context) {
		personList, _ := GetAllPersonFromDB(db)
		personGroup = PersonGroup{&personList}
		c.HTML(http.StatusOK,
			"Index", gin.H{"count": len(personList)})
	})
	v1 := router.Group("/lucky")
	{
		v1.GET("/all", GetAll)
		v1.GET("/totalNum", GetTotalNum)
		v1.GET("/luckyNum", GetLuckyNum)

	}
	port := Port
	if len(os.Getenv("PORT")) > 0 {
		port = os.Getenv("PORT")
	}
	fmt.Println("Listening and serving HTTP on" + ":" + port)
	router.Run(":" + port)
}
