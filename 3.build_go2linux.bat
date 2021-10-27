
SET GOOS=linux
SET GOARCH=amd64
go env

go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,https://goproxy.io,direct
go build
echo build done