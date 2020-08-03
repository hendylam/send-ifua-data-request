package main

import (
	"fmt"
	"time"

	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	pool "github.com/processout/grpc-go-pool"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"tapera.integrasi/api/constant"
	bri "tapera.integrasi/api/controller/ifua"
	_ "tapera.integrasi/api/docs"
	"tapera.integrasi/api/middleware"

	"tapera/conf"
	"tapera/util"
	"tapera/util/env"
)

// @title Tapera API
// @version v1.0.0
// @description This is Tapera API listing descriptions.
// @termsOfService http://Tapera.org/terms/

// @contact.name Tapera API Support
// @contact.url http://www.Tapera.org/support
// @contact.email support@Tapera.org

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	// load env variable if the is .env file exists
	if util.FileExists(".env") {
		err := godotenv.Load()
		if err != nil {
			panic(err)
		}
	}

	// create appconfig
	appConf := conf.EnvConfig()
	conf.AppConf = appConf

	// create factory based on env
	envf := appConf.NewEnvFactory()

	// create logger
	logger := envf.Logger()

	// create db connection
	// db := envf.CockcroachDb()
	// defer db.Close()

	// create mux router
	r := mux.NewRouter()

	// add middlewares
	r.Use(handlers.CompressHandler,
		middleware.AddRequestID,
		middleware.Logger(logger),
		middleware.ReqLogger(),
	)

	// register swagger route
	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	// create grpc connection pool
	// miGrpcAddr := env.Str(constant.EnvGrpcServerMitraIntegrasi, true, nil)
	// grpcClientCnPool := createGrpcClientCnPool(miGrpcAddr)
	// defer grpcClientCnPool.Close()

	// // create grpc client manager
	// miClientMgr := mic.NewGrpcClientManager(grpcClientCnPool)

	//create bri controller
	bri.NewController().Route(r)

	// run http server
	logger.Info().Msgf("server is listening to port %d", appConf.Port())
	if err := http.ListenAndServe(fmt.Sprintf(":%d", appConf.Port()), handlers.RecoveryHandler()(r)); err != nil {
		panic(err)
	}
	logger.Info().Msg("exit")
}

func createGrpcClientCnPool(addr string) *pool.Pool {
	// create grpc connection pool see https://github.com/processout/grpc-go-pool
	maxCn := env.Int(constant.EnvGrpcClientMaxOpenConn, false, 10)
	maxIddleCn := env.Int(constant.EnvGrpcClientMaxIdleConn, false, 0)
	idleTimeOut := env.Int(constant.EnvGrpcClientIdleTimeoutSec, false, 10)
	cert := env.Str(constant.EnvGrpcClientCert, false, nil)

	if len(cert) != 0 {
		creds, err := credentials.NewClientTLSFromFile(cert, "")
		if err != nil {
			panic(err)
		}

		gpool, err := pool.New(func() (*grpc.ClientConn, error) {
			return grpc.Dial(addr, grpc.WithTransportCredentials(creds))
		}, maxIddleCn, maxCn, time.Duration(idleTimeOut)*time.Second)

		if err != nil {
			panic(err)
		}
		return gpool
	}

	gpool, err := pool.New(func() (*grpc.ClientConn, error) {
		return grpc.Dial(addr, grpc.WithInsecure())
	}, maxIddleCn, maxCn, time.Duration(idleTimeOut)*time.Second)

	if err != nil {
		panic(err)
	}
	return gpool
}
