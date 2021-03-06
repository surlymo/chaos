#Chaos

注意：microservice相关的部分即将从主项目chaos中剥离。

##QuickStart
### 私有云平台
+ 如需要可执行程序，请执行go install opensource/chaos/server (相对于$GOPATH/src下的项目放置位置。目前开发版本放置在示例的目录下)，之后到bin目录下找到server的可执行程序双击之后。打开浏览器输入`localhost:8080`即可登录私有云平台界面。同时rest api也全部生效。
+ 如希望直接执行，则请执行`go run $YOUR_PATH_PREFIX/server/*.go`之后即全部生效。

**目前界面提供如下功能**
+ 私有云首页入口
+ 罗列服务集群现状概况，并提供搜索排序功能

**目前提供如下rest api**
+ /deploy/apps/rollback 快速回滚到某一个版本或者回滚到最近的版本
+ /deploy/apps 对单个服务进行一键部署
+ /deploy/apps/updater 对单个或者多个服务进行新增或者更新
+ /deploy/groups 对单个或者多个组进行一键部署
+ /info 获取所有服务信息

### 微服务容器
每个容器都会带有如下几个基础进程：
+ sshd 便于进行ssh登录
+ supervisor 守护进程
+ consul-template 服务发现
+ registrator 服务注册
+ haproxy 路由、负载均衡
+ dnsmasq 本地dns解析，将*.service.consul域名全部解析到127.0.0.1

**步骤**
+ 修改对应Dockerfile里面的IP地址为你所用的地址。
+ 确保要部署的服务被打包成output.tar.gz,并且在解压后的根目录下带有`run.sh`和`dependency`两个文件。`dependency`代表依赖的服务，一行一个。如果是http服务如tomcat，则域名为tomcat.service.consul，`dependency`文件中写入**tomcat**；如果是tcp服务如zookeeper，则访问zk的域名为zk.service.consul:81，`dependency`文件中则写入**zk:81**(端口可自定义，不冲突即可)。若无依赖，填入"**empty**"。`run.sh`为启动脚本。
+ 每个宿主机都启动Consul Agent，其中选择3-5个作为Server节点。组成集群。
+ 启动时候容器需要携带环境变量`$SERVICE_NAME`和`$SERVICE_PORT`或在容器Dockerfile中声明。结束。至此微服务所需的负载均衡、服务发现、服务注册、路由已经都起来了。

**go get**
go get github.com/ant0ine/go-json-rest/rest
go get github.com/garyburd/redigo/redis
go get github.com/samalba/dockerclient
go get github.com/tjz101/goprop
go get gopkg.in/mgo.v2
go get github.com/fsouza/go-dockerclient(golang.org/x/net/context被墙，自行处理~)
go get github.com/hashicorp/consul/api  

##Introdution
丁丁私有云平台server。提供微服务基础组件、监控、历史记录以及上线部署等的支持。具体后续更新

## Todo Recently
+ 一键部署相关的rest接口完善，并提供异步化
+ 结合marathon和mesos以及docker，提取服务的真实ip等关键信息，并进行视图呈现
+ 搭建事件动态刷新前后端server架构。对于服务变更进行第一时间监控和实时刷新。
+ 对centos6.5镜像进行完善。对网络问题进行优化。

## Contact Us
inf@zufangit.cn

## Changelog

**v0.5** —— **2016-03-17**
+ 服务注册、服务发现、负载均衡、路由改为haproxy+consul+template+registrator+dnsmasq来进行。
+ 修复tcp长连接等其他bug。

**v0.4** —— **2016-02-24**
+ 添加registrator。用以和consul结合，实现服务发现和服务注册

**v0.3** —— **2016-01-25**
+ 提供一键部署的四个基本rest api
+ 模块抽象和拆解。代码优化

**v0.2** —— **2016-01-25**
+ 前后端架构重构。改为rest + angularjs-ajax的方式
+ 前端代码抽象和拆解。
+ 提供获取所有服务信息的接口

**v0.1** —— **2016-01-13**
+ chaos初始化