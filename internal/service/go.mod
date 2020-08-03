module tapera.integrasi/service

go 1.14

replace tapera/util => ../util

replace tapera.integrasi/grpc/client => ../grpc/client

require (
	github.com/golang/protobuf v1.4.2
	github.com/processout/grpc-go-pool v1.2.1 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.6.0
	google.golang.org/genproto v0.0.0-20200608115520-7c474a2e3482 // indirect
	google.golang.org/grpc v1.30.0
	gopkg.in/DATA-DOG/go-sqlmock.v1 v1.3.0
	tapera.integrasi/grpc/client v1.0.0
	tapera/util v1.0.0
)
