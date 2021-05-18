module ligomonitor

go 1.14

require (
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575
	github.com/gin-gonic/gin v1.6.3
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/stretchr/testify v1.5.1 // indirect
	golang.org/x/tools v0.1.0 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
)

replace (
	github.com/containerd/containerd => /home/guozhaocoder/go/src/github.com/containerd
	github.com/docker/docker v1.13.1 => github.com/docker/engine v17.12.0-ce-rc1.0.20200916142827-bd33bbf0497b+incompatible
	github.com/gogo/protobuf => /home/guozhaocoder/go/src/github.com/gogo/protobuf
	github.com/google/go-cmp => /home/guozhaocoder/go/src/github.com/go-cmp
	github.com/gorilla/mux => /home/guozhaocoder/go/src/github.com/mux
	github.com/sirupsen/logrus => /home/guozhaocoder/go/src/github.com/logrus
	golang.org/x/time => github.com/golang/time v0.0.0-20201208040808-7e3f01d25324
)
