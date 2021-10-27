package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/http"
	"os"
	"path"
)

type configuration struct {
	Port   string
	Path   string
	DbName string
	Proxy  string
	BusUrl string
}

var CONF = configuration{}

// 定义 图片目录
var imgPath string

func init() {
	// 读取配置文件
	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&CONF)
	checkErr(err)
	fmt.Println(CONF)

	appPath, _ := os.Getwd()
	imgPath = path.Join(appPath, CONF.Path)

	// 判断目录是否存在
	_, err = os.Stat(imgPath)
	if err != nil {
		err := os.MkdirAll(imgPath, os.ModePerm)
		checkErr(err)
	}

	// connection
	db, err := gorm.Open(sqlite.Open(CONF.DbName), &gorm.Config{})
	checkErr(err)
	db.AutoMigrate(&fanhao{})
	sqlDB, err := db.DB()
	if err == nil {
		defer sqlDB.Close()
	}

	fmt.Println("photoPath:", imgPath)
	fmt.Println("dataBases:",CONF.DbName)
}
func main() {

	// 设置静态资源
	static := &StaticResource{staticFS: WebUI, path: "static"}
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(Cors())

	r.StaticFS("/static/", http.FS(static))
	r.Static("/photos/", CONF.Path)
	r.GET("/", handleIndex)

	r.GET("/api/getList", handleApiGetList)
	r.GET("/api/set", handleApiSet)
	r.GET("/api/search", handleApiSearch)
	r.GET("/api/del", handleApiDel)
	r.GET("/api/test", handleApiTest)

	r.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })
	err := r.Run(":" + CONF.Port)
	checkErr(err)
}

func main1() {
	//tmi := time.Unix(1492442108,0).Format("2006-01-02 15:04:05")
	//fmt.Println(tmi)
	//tmi1,_ := time.Parse("2006-01-02 15:04:05",tmi)
	//fmt.Println(tmi1)
	//fmt.Println(tmi1.Unix())
	//fmt.Println(time.Now().Unix())

	//db,err := gorm.Open(sqlite.Open(CONF.DbName),&gorm.Config{
	//	Logger: logger.Default.LogMode(logger.Info),
	//})
	//checkErr(err)

	//err = db.AutoMigrate(&fanhao{})
	//checkErr(err)

	//var fan fanhao
	//db.First(&fan)
	//fmt.Println(fan.UpdatedAt)
	//fan.Starnum=2
	//db.Save(fan)

	// all code
	//type onlyCode struct { Code string }
	//var codes []onlyCode
	//db.Model(&fanhao{}).Find(&codes)
	//fmt.Println(codes)

	//var fans []fanhao
	//db.Find(&fans)
	//fmt.Println(fans)

	// connection
	db, err := gorm.Open(sqlite.Open(CONF.DbName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	checkErr(err)

	var fans fanhao
	db.First(&fans)
	fmt.Print(fans)
}

func main2() {

	// 关闭文件
	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&CONF)
	checkErr(err)
	fmt.Println(CONF)

}

func main3() {
	fmt.Println("123")
}