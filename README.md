<p align="center">
	<h1 align="center">番械库</h1>
	<h3 align="center">私人番号收藏</h3>
</p>

## 简介
每个人都有自己喜欢的番号，每个人都想有自己记录番号的小本本。。

当小本本记录了太多之后我们想找到今日想要的哪一部时，究竟哪行神秘代码对应哪部电影，想选择却无从下手

现在有个番号收藏，管理番号，so easy !

添加番号收藏后会自动添加番号对应的封面图片，演员，还可手动标记骑兵步兵

## 注意事项
根目录下"fanhao.db"文件为数据库文件,请注意存档备份;如使用docker部署可存放在其他目录更安全,操作方式下面有

## 部署方式1:docker(推荐!推荐!推荐!)
### 1. 编译前端 必须先编译前端
```
docker-compose -f docker-compose-build-1web.yml up

直到看到"build done"容器自动退出则编译成功
```

### 2. 编译后端
```
docker-compose -f docker-compose-build-2go.yml up
直到看到"build done"容器自动退出则编译成功
```

### 3. 运行 使用nginx代理静态资源
```
docker-compose -f docker-compose-run-1nginx.yml up

启动成功后访问 http://ip:8888

端口在 docker-compose-run-1nginx.yml 中修改

```

### 4. 运行 不使用nginx代理静态资源
```
docker-compose -f docker-compose-run-2single.yml up

启动成功后访问 http://ip:8888

端口在 docker-compose-run-2single.yml 中修改

```

此部署模式下 数据库文件"fanhao.db"和图片目录"photos" 都可以在yml文件中指定和修改

### 5. 修改数据存储位置

```
在"docker-compose-run-1nginx.yml"和"docker-compose-run-2single.yml"里均有被注释的
"
#      - /data/dcdb/fanhao/photos:/app/photos
#      - /data/dcdb/fanhao/fanhao.db:/app/fanhao.db
"
以上两行,注释掉就是使用项目目录存放数据,解开注释就是使用自定义目录
在使用自定义目录前,需要先启动一次生成数据库文件,然后自己建好目录,将fanhao.db和photos移动到新目录,再修改路径,后面的"/app/photos"和"/app/fanhao.db"与"conf.json"配置文件挂钩不建议修改
```


## 部署方法2:手动
### 1. 编译前端 必须先编译前端
```
cd web # 进入前端目录
npm i # 安装相关组件,速度和网速,代理等有关
npm run build # 开始编译
```

### 2. 编译后端
```
go build # 开始编译,速度和网速,代理等有关
```
### 3. 运行

```
运行当前目录下生成的"fanhao"或"番fahao.exe"文件即可
```

### 4. 配置nginx反向代理

```
既然选择手动部署,我相信你动手能力一定很强,加油鸭
```

## 问题处理
### 获取番号信息时报错,
```
1.请检查 busUrl 地址是否被墙,busUrl地址可能会更换,使用浏览器打开busUrl网址看是否能访问,如果还不能访问请发issue,题主会定期更换

2.程序所在环境是否能正常访问到小飞机的代理,

3.如果代理没有问题,则可能是代理环境无法访问配置文件中配置的"busUrl"网址,则把小飞机更换为全局代理模式实施
```

## 技术参考
> golang v1.17+ , node v16.10+
> docker 需要 配套 docker-compose
> 数据库使用sqlite,默认保存在根目录 "fanhao.db" 中,注意备份!!!注意备份!!!注意备份!!!


## 关于代理
在 "conf.json" 中配置 proxy,默认为"",配置如"socks5://127.0.0.1:1080" 或 "http://127.0.0.1:1081"

确保本地小飞机软件已打开,如代理不在本机则确保设置"允许其他设备接入"已勾选,

socks代理端口默认为1080 , http代理端口默认为1081,请根据实际情况填写 (具体方法请绅士自查)

修改配置需要重启程序

## license
> 本作品仅供学习交流使用，对使用后产生的任何后果不承担任何责任; 前方净空,允许进入,祝君武运昌隆

## todo:
> ~~配置项集中到配置文件~~
>
> ~~换用sqlite存储方式~~
>
> ~~自动化初始脚本~~