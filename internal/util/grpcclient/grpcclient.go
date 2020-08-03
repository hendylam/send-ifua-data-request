package grpcclient

import (
	"time"

	pool "github.com/processout/grpc-go-pool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// ConnPoolConfig struct
type ConnPoolConfig struct {
	Addr         string
	InitConn     int
	MaxConn      int
	IddleTimeOut uint
	Cert         string
}

// CreateConnPool func
func (cfg *ConnPoolConfig) CreateConnPool() *pool.Pool {
	return createConnPool(cfg.Addr, cfg.InitConn, cfg.MaxConn, cfg.IddleTimeOut, cfg.Cert)
}

// CreateConnPool func
func CreateConnPool(addr string) *pool.Pool {
	return createConnPool(addr, 0, 50, 5, "")
}

func createConnPool(addr string, intCn int, maxCn int, iddleTimeOut uint, cert string) *pool.Pool {
	if len(cert) != 0 {
		creds, err := credentials.NewClientTLSFromFile(cert, "")
		if err != nil {
			panic(err)
		}
		gpool, err := pool.New(func() (*grpc.ClientConn, error) {
			return grpc.Dial(addr, grpc.WithTransportCredentials(creds))
		}, intCn, maxCn, time.Duration(iddleTimeOut)*time.Second)

		if err != nil {
			panic(err)
		}
		return gpool
	}

	gpool, err := pool.New(func() (*grpc.ClientConn, error) {
		return grpc.Dial(addr, grpc.WithInsecure())
	}, intCn, maxCn, time.Duration(iddleTimeOut)*time.Second)

	if err != nil {
		panic(err)
	}
	return gpool
}
