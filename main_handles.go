package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func handleIndex(c *gin.Context) {
	indexHTML, _ := WebUI.ReadFile("web/dist/index.html")
	_, err := c.Writer.Write(indexHTML)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
	} else {
		c.Writer.WriteHeader(http.StatusOK)
	}
	c.Writer.Header().Add("Accept", "text/html")
	c.Writer.Flush()
}
func handleApiGetList(c *gin.Context) {
	// connection
	db, err := gorm.Open(sqlite.Open(CONF.DbName), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		c.String(http.StatusBadRequest, "failed open db")
		return
	}
	sqlDB, err := db.DB()
	if err == nil {
		defer sqlDB.Close()
	}

	var querys = make(map[string]interface{})

	if star, ex := c.GetQuery("star"); ex {
		if len(star) > 0 {
			querys["star"] = star
		}
	}
	if code, ex := c.GetQuery("code"); ex {
		if len(code) > 0 {
			querys["code"] = code
		}
	}
	if ma, ex := c.GetQuery("ma"); ex {
		if len(ma) > 0 {
			querys["ma"] = ma
		}
	}
	if dn, ex := c.GetQuery("downed"); ex {
		if len(dn) > 0 {
			querys["downed"] = dn
		}
	}
	//fmt.Println(querys)

	size := 48
	if re, ex := c.GetQuery("size"); ex {
		if ai, err := strconv.Atoi(re); err == nil {
			size = ai
		}
	}
	index := 1
	if re, ex := c.GetQuery("index"); ex {
		if ai, err := strconv.Atoi(re); err == nil {
			index = ai
		}
	}

	var fans []fanhao
	db.Limit(size).Offset((index - 1) * size).Where(querys).Order("id desc").Find(&fans)
	//fmt.Print(fans)

	var count int64
	db.Model(&fanhao{}).Where(querys).Count(&count)

	//// all code
	type onlyCode struct{ Code string }
	var codes []onlyCode
	db.Model(&fanhao{}).Order("Code").Find(&codes)
	//fmt.Println(codes)
	//// all star
	type onlyStar struct{ Star string }
	var stars []onlyStar
	db.Model(&fanhao{}).Group("Star").Order("Star").Find(&stars)
	//fmt.Println(stars)

	c.JSON(200, gin.H{
		"codes":  codes,
		"stars":  stars,
		"fans":   fans,
		"size":   size,
		"index":  index,
		"count":  count,
		"busurl": CONF.BusUrl,
	})
}
func handleApiSet(c *gin.Context) {
	tp, has := c.GetQuery("t")
	if has != true {
		c.String(http.StatusBadRequest, "null: t")
		return
	}

	id, has := c.GetQuery("id")
	if has != true {
		c.String(http.StatusBadRequest, "null: id")
		return
	}
	iId, err := strconv.Atoi(id)
	if err != nil {
		c.String(http.StatusBadRequest, "must be number: id")
		return
	}

	flag, has := c.GetQuery("flag")
	if has != true {
		c.String(http.StatusBadRequest, "null: flag")
		return
	}
	iFlag, err := strconv.Atoi(flag)
	if err != nil {
		c.String(http.StatusBadRequest, "must be number: flag")
		return
	}

	// connection
	db, err := gorm.Open(sqlite.Open(CONF.DbName), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		c.String(http.StatusBadRequest, "failed open db")
		return
	}
	sqlDB, err := db.DB()
	if err == nil {
		defer sqlDB.Close()
	}

	var fan fanhao
	res := db.Where("id", iId).First(&fan)
	if res.RowsAffected != 1 {
		c.String(http.StatusBadRequest, "arg err: id")
		return
	}
	switch tp {
	case "st":
		fan.Starnum = iFlag
		break
	case "ma":
		fan.Ima = iFlag
	case "fa":
		fan.Iface = iFlag
	case "dn":
		fan.Downed = iFlag
	}
	res = db.Save(&fan)

	//fmt.Println(res)
	if res.RowsAffected != 1 {
		c.String(http.StatusBadRequest, "update err")
		return
	} else {
		c.String(http.StatusOK, "ok")
	}

}
func handleApiSearch(c *gin.Context) {
	//-- 校验参数
	code, has := c.GetQuery("c")
	if has != true {
		c.String(http.StatusBadRequest, "null: c")
		return
	}
	code = strings.ToUpper(strings.Trim(code, " "))
	url := CONF.BusUrl + code

	_, justUp := c.GetQuery("u")

	//-- 请求页面数据
	res, err := _request(url)
	if err != nil {
		c.String(http.StatusBadRequest, "request err")
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "request read body err")
		return
	}

	html := string(body)
	
	//-- 解析页面数据
	regCode, err := regexp.Compile("<span class=\"header\">識別碼:</span>.*?<span style=\"color:#CC0000;\">(.*?)</span>")
	if err != nil {
		c.String(http.StatusBadRequest, "failed match code")
		return
	}
	matchCode := regCode.FindStringSubmatch(html)

	regTitle, err := regexp.Compile("<h3>(.*?)</h3>")
	if err != nil {
		c.String(http.StatusBadRequest, "failed match title")
		return
	}
	matchTitle := regTitle.FindStringSubmatch(html)

	regStar, err := regexp.Compile("<div class=\"star-name\"><a href=\"https://.*?/star/.*?\" title=\".*?\">(.*?)</a></div>")
	if err != nil {
		c.String(http.StatusBadRequest, "failed match star")
		return
	}
	matchStar := regStar.FindStringSubmatch(html)

	regStarCode, err := regexp.Compile("<div class=\"star-name\"><a href=\"https://.*?/star/(.*?)\" title=\".*?\">.*?</a></div>")
	if err != nil {
		c.String(http.StatusBadRequest, "failed match starcode")
		return
	}
	matchStarCode := regStarCode.FindStringSubmatch(html)

	regImg, err := regexp.Compile("<a class=\"bigImage\" href=\"(.*?)\">")
	if err != nil {
		c.String(http.StatusBadRequest, "failed match img")
		return
	}
	matchImg := regImg.FindStringSubmatch(html)

	// connection
	db, err := gorm.Open(sqlite.Open(CONF.DbName), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		c.String(http.StatusBadRequest, "failed open db")
		return
	}
	sqlDB, err := db.DB()
	if err == nil {
		defer sqlDB.Close()
	}

	//-- 查询数据
	var fan fanhao
	dbGetRes := db.Where("Code", code).First(&fan)

	if !justUp && dbGetRes.RowsAffected == 1 {
		c.String(http.StatusOK, "has")
		return
	}

	if len(matchCode) > 1 {
		fan.Code = matchCode[1]
	}
	if len(matchTitle) > 1 {
		fan.Title = matchTitle[1]
	}
	if len(matchStar) > 1 {
		fan.Star = matchStar[1]
	} else {
		fan.Star = "-暂无-"
	}
	if len(matchStarCode) > 1 {
		fan.StarCode = matchStarCode[1]
	}
	if len(matchImg) > 1 {
		fan.Img = matchImg[1]
	}
	fan.Fname = fan.Code + ".jpg"

	dbSaveRes := db.Save(&fan)

	if dbSaveRes.RowsAffected != 1 {
		c.String(http.StatusBadRequest, "update err")
		return
	}

	// 开始保存图片
	imgUrl := CONF.BusUrl + fan.Img
	resImg, err := _request(imgUrl)
	if err != nil {
		c.String(http.StatusBadRequest, "request img err")
		return
	}
	defer resImg.Body.Close()
	imgBody, err := ioutil.ReadAll(resImg.Body)
	checkErr(err)
	if err != nil {
		c.String(http.StatusBadRequest, "request img read err")
		return
	}

	fPath := path.Join(CONF.Path, fan.Fname)
	file, err := os.Create(fPath)
	defer file.Close()
	if err != nil {
		c.String(http.StatusBadRequest, "request img create err")
		return
	}
	_, err = io.Copy(file, bytes.NewReader(imgBody))
	if err != nil {
		c.String(http.StatusBadRequest, "request img save err")
		return
	}

	c.String(http.StatusOK, "ok")

}
func handleApiDel(c *gin.Context) {

	id, has := c.GetQuery("qazxsw")
	if has != true {
		c.String(http.StatusBadRequest, "null: id")
		return
	}
	iId, err := strconv.Atoi(id)
	if err != nil {
		c.String(http.StatusBadRequest, "must be number: id")
		return
	}

	// connection
	db, err := gorm.Open(sqlite.Open(CONF.DbName), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		c.String(http.StatusBadRequest, "failed open db")
		return
	}
	sqlDB, err := db.DB()
	if err == nil {
		defer sqlDB.Close()
	}

	dbDelRes := db.Delete(&fanhao{}, iId)

	c.String(http.StatusOK, string(dbDelRes.RowsAffected))

}
func handleApiTest(c *gin.Context) {
	// connection
	db, err := gorm.Open(sqlite.Open(CONF.DbName), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		c.String(http.StatusBadRequest, "failed open db")
		return
	}
	sqlDB, err := db.DB()
	if err == nil {
		defer sqlDB.Close()
	}

	var fans fanhao
	db.First(&fans)
	//fmt.Print(fans)

	c.JSON(200, gin.H{
		"fans": fans,
	})
}
