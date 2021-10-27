package main

import (
	"embed"
	"io/fs"
	"path"
)

// StaticResource 嵌入普通静态资源
type StaticResource struct {
	// 静态资源
	staticFS embed.FS
	// 设置embed文件到静态资源的相对路径，也就是embed注释里的路径
	path string
}

// Open 静态资源被访问逻辑
func (_this_ *StaticResource) Open(name string) (fs.File, error) {
	return _this_.staticFS.Open(path.Join("web/dist", name))
}

//go:embed web/dist
var WebUI embed.FS