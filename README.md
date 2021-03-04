#### 介绍
>这是一个用Golang实现的Linux监控系统后台，适用于x86_64，目前功能还较为简单，至于前端，本人抽不出来时间写，苦涩。
>
>**常规功能**：该系统实现了进程资源，系统资源实时监控的接口，类似与Linux的top命令，通过websocket推送到前端。集成了mysql与mongodb数据库，用于存储进程资源等信息，便于后期对某个进程做性能分析。
>
>**特色功能**：可以查看某个进程的系统调用。通过Go的syscall包，使用ptrace系统调用追踪某个进程的系统调用情况。
##### 运行项目

```sh
$ cd <项目根目录>
$ cd script
$ sh build.sh
```

##### 项目目录介绍

```shell script
├── cmd                             #入口包             
│   ├── app                         #初始化包
│   │   └── ligoconf.go                #初始化配置信息
│   └── main.go                        #启动入口文件
├── configs                         #配置信息包
│   ├── conf.json                      #全局配置文件
│   ├── logcfg.xml                     #日志配置文件
│   ├── logs.log                       #日志
│   └── 配置文件说明.md
├── front                           #前端静态文件包
├── pkg                             #业务处理包
│   ├── api                         #路由包
│   │   └── router.go                  #路由接入文件    
│   ├── conn                        #连接提供包
│   │   ├── dbclient.go                #数据库客户端
│   │   └── dockerclient.go            #docker客户端
│   ├── cons                        #常量包
│   │   ├── exitcode.go                #退出码
│   │   └── httpcode.go                #http请求标识码
│   ├── model                       #实体类层
│   │   ├── conf.go                    #配置信息实体类
│   │   ├── process.go                 #进程信息实体类
│   │   ├── ptrace.go                  #trace信息实体类
│   │   ├── reqresp.go                 #请求返回信息实体类
│   │   └── resource.go                #主机信息实体类
│   └── service                     #总业务处理层
│       ├── ctl                     #控制器层，与router对接
│       │   ├── hostctl.go             #主机信息控制器
│       │   └── processctl.go          #进程控制器
│       ├── db                      #数据库访问层
│       │   ├── dbinterface.go         #数据库访问接口
│       │   ├── mongodb.go             #mongo实现
│       │   └── mysql.go               #mysql实现
│       ├── dbsync                  #数据库同步业务层
│       │   ├── dbsvc.go               #数据库同步业务文件
│       │   └── dbtimer.go             #数据库同步定时器
│       ├── host                    #主机资源访问层
│       │   ├── cpu.go                 #cpu资源提供
│       │   ├── memory.go              #内存资源提供
│       │   ├── process.go             #进程资源提供
│       │   └── ptrace.go              #系统调用资源提供
│       └── svc                     #具体业务处理层
│           ├── hostsvc.go             #主机资源业务处理
│           └── processsvc.go          #进程资源业务处理
├── script                          #脚本部署层
├── test                            #测试层
└── utils                           #工具类层
    ├── errsplitutil.go                #异常处理
    ├── fileutil.go                    #文件处理
    └── logutil.go                     #日志处理
```