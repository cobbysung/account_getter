
#### 介绍
- golang练手项目，抓取网络上的百度网盘、迅雷....等账号
- 代码仅供学习参考使用，请勿用于非法用途
- 代码仅供学习参考使用，请勿用于非法用途
- 代码仅供学习参考使用，请勿用于非法用途

#### 软件架构

<p align="center">
	<img src="https://raw.githubusercontent.com/cobbysung/account_getter/master/structure.png"  width="80%" height="80%">
	<p align="center">
		<em>架构图</em>
	</p>
</p>

* engine
* scheduler
* parser
* fetcher
* models
* proxy-pool
* http-server


#### 安装教程

1. 相关依赖
```
gopm list
Dependency list (3):
-> github.com/PuerkitoBio/goquery
-> golang.org/x/net
-> golang.org/x/text
```

2. 安装依赖
```go
go get github.com/PuerkitoBio/goquery
go get golang.org/x/net
go get golang.org/x/text
```
3. 编译
```golang
go build entry.go
```

#### 使用说明



```
Usage of ./entry:
  -http-port int
        HTTP服务端口 (default 8000)
  -update
        是否需要重新更新数据 (default true)
  -worker-num int
        普通worker数 (default 5)
 ```
 
- 抓取信息
 `./entry`
- 不重新抓取，只展现已有账户
 `./entry -update=0`
- 指定web-server端口
`./enty -http-port=8111`

- 查看获取账号，打开地址：[http://localhost:8000/home](http://localhost:8000/home)
 
- 演示
 ![img](https://raw.githubusercontent.com/cobbysung/account_getter/master/demo.gif)

