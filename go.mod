module github.com/siangyeh8818/gitlab.api

go 1.14

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

// google.golang.org/grpc 1.26后的版本是不支持clientv3的

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5

//解決 go mod 的bug
//        module declares its path as: go.etcd.io/bbolt
//                but was required as: github.com/coreos/bbolt

require (
	github.com/ThreeDotsLabs/watermill v1.1.1
	github.com/ThreeDotsLabs/watermill-nats v1.0.5
	github.com/coreos/bbolt v0.0.0-00010101000000-000000000000 // indirect
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/gin-gonic/gin v1.7.2
	github.com/go-ini/ini v1.62.0
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/btree v1.0.1 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/nats-io/nats.go v1.11.0
	github.com/nats-io/stan.go v0.9.0
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20201229170055-e5319fda7802 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.uber.org/zap v1.18.1 // indirect
	golang.org/x/time v0.0.0-20210611083556-38a9dc6acbc6 // indirect
	google.golang.org/genproto v0.0.0-20210719143636-1d5a45f8e492 // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)
