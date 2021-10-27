package main

import (
	"crypto/tls"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/url"
)


func checkErr(err error) {
	if err != nil {
		log.Fatal("err:", err)
		//panic(err)
	}
}
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.RequestURI
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

// 请求二进制
func _request(u string) (*http.Response, error) {

	tr:=&http.Transport{TLSClientConfig: &tls.Config{
		InsecureSkipVerify: true,
	}}

	if len(CONF.Proxy) > 0{
		proxyUrl, err := url.Parse(CONF.Proxy)
		checkErr(err)
		tr.Proxy = http.ProxyURL(proxyUrl)
	}

	client := &http.Client{
		Transport: tr,
	}
	//提交请求
	reqest, err := http.NewRequest("GET", u, nil)
	//增加header选项
	reqest.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36")

	if err != nil {
		return nil, err
	}
	//处理返回结果
	response, err := client.Do(reqest)

	if err != nil {
		return response, err
	}
	if response.StatusCode != 200 {
		return response, http.ErrMissingFile
	}
	return response, err
}
