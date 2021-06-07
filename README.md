#### 介绍
>这是一个用Golang实现的Linux监控系统后台，适用于x86_64，目前功能还较为简单，至于前端，本人抽不出来时间写，苦涩。
>
>**常规功能**：该系统实现了进程资源，系统资源实时监控的接口，类似与Linux的top命令，通过websocket推送到前端。集成了mysql与mongodb数据库，用于存储进程资源等信息，便于后期对某个进程做性能分析。
>也可以通过接口进行注册进程使用资源阈值以及处理方式，超过阈值进行处理。
>
>**特色功能**：可以查看某个进程的系统调用。通过Go的syscall包，使用ptrace系统调用追踪某个进程的系统调用情况，但是该功能并不完善，比如当被跟踪进程进入系统调用不出来，这时候跟踪进程就没办法退出跟踪了（除非把自己杀死），具体解决方法思路：启动一个子进程来跟踪别的进程，
>通过通道或者RPC通讯，目前在想解决步骤。
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

### 主要接口

####  进程资源推送

+ URL: IP:Port/monitor/process/stream

+ Method: GET

+ Response
```json
{
  "code":200,
  "message":"success",
  "data":
    [
      {"pid":26356,"pPid":26354,"name":"c:/Program Files/Tencent/WeChat/WeChat.exe\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\0000\u0000\u0000\u0000","tgid":26356,"state":"S","uid":1000,"gid":1000,"threads":42,"vmPeak":3848840,"vmSize":2882272,"vmLck":64,"vmPin":0,"vmHWM":512524,"vmRss":388792,"vmData":326764,"vmStk":132,"vmExe":1904360,"vmLib":0,"vmPTE":1260,"vmSwap":0,"voluntaryCS":16609838,"noVoluntaryCS":6623,"cpuUsage":0.18170656,"cpuUsed":34011,"time":0,"ifAlarm":false,"alarmMessage":{"cpuMsg":"","vmMsg":""}},
      {"pid":26361,"pPid":1,"name":"/usr/lib/i386-linux-gnu/deepin-wine/./wineserver.real\u0000-p0\u0000","tgid":26361,"state":"S","uid":1000,"gid":1000,"threads":1,"vmPeak":17860,"vmSize":8820,"vmLck":0,"vmPin":0,"vmHWM":15108,"vmRss":5564,"vmData":4360,"vmStk":132,"vmExe":484,"vmLib":3772,"vmPTE":48,"vmSwap":0,"voluntaryCS":17262641,"noVoluntaryCS":4111,"cpuUsage":0.07366482,"cpuUsed":11203,"time":0,"ifAlarm":false,"alarmMessage":{"cpuMsg":"","vmMsg":""}},
      {"pid":4473,"pPid":3605,"name":"/opt/google/chrome/chrome --type=renderer --field-trial-handle=914862911565256638,18363096547406141597,131072 --lang=zh-CN --enable-crash-reporter=00eb75a2-7563-43c0-a281-cd132289c04e, --origin-trial-disabled-features=SecurePaymentConfirmation --num-raster-threads=4 --enable-main-frame-before-activation --renderer-client-id=190 --no-v8-untrusted-code-mitigations --shared-files=v8_context_snapshot_data:100","tgid":4473,"state":"S","uid":1000,"gid":1000,"threads":19,"vmPeak":42430596,"vmSize":38423504,"vmLck":0,"vmPin":0,"vmHWM":175040,"vmRss":157116,"vmData":244088,"vmStk":132,"vmExe":161068,"vmLib":9404,"vmPTE":1292,"vmSwap":5764,"voluntaryCS":1628,"noVoluntaryCS":75,"cpuUsage":0.06875384,"cpuUsed":173,"time":0,"ifAlarm":false,"alarmMessage":{"cpuMsg":"","vmMsg":""}},
      ......
    ]
}
```


####  进程快照获取

+ URL: IP:Port/monitor/process/snapshot/{pid}

+ Method: GET

+ Response
```json
{
  "code":200,
  "message":"success",
  "data":
    [
      { 
        "pid":26356,
        "pPid":26354,
        "name":"c:/Program Files/Tencent/WeChat/WeChat.ex0000",
        "tgid":26356,
        "state":"S",
        "uid":1000,
        "gid":1000,
        "threads":42,
        "vmPeak":3848840,
        "vmSize":2882272,
        "vmLck":64,"vmPin":0,
        "vmHWM":512524,
        "vmRss":385484,
        "vmData":326764,
        "vmStk":132,
        "vmExe":1904360,
        "vmLib":0,
        "vmPTE":1260,
        "vmSwap":0,
        "voluntaryCS":20189781,
        "noVoluntaryCS":12706,
        "cpuUsage":0.18742293,
        "cpuUsed":38731,
        "time":0,
        "ifAlarm":false
      }
    ]
}
```


####  强制停止进程

+ URL: IP:Port/monitor/process/kill/{pid}

+ Method: GET

+ Response
```json
{
  "code":200,
  "message":"success",
  "data":null
}
```


####  获取进程近况

+ URL: IP:Port/monitor/process/status/{pid}

+ Method: GET

+ Response
```json

```


####  注册进程资源阈值

+ URL: IP:Port/monitor/process/alarm/register

+ Method: POST

+ Request
```json
{ 
  "pid":4084,         #进程pid
  "vmLimit":9999999,  #进程物理内存阈值
  "cpuLimit":0.01,    #进程CPU使用阈值
  "operate":2         #告警操作
}
```
+ Response
```json

```


####  跟踪某进程

+ URL: IP:Port/monitor/process/ptrace/{pid}

+ Method: GET

+ Response
```json

```


####  查看主机资源快照

+ URL: IP:Port/monitor/host/snapshot

+ Method: GET

+ Response
```json

```



####  查看主机资源流

+ URL: IP:Port/monitor/host/stream

+ Method: GET

+ Response
```json

```