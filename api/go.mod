module tapera.integrasi/api

go 1.14

replace tapera/util => ../internal/util

replace tapera/conf => ../internal/conf

replace tapera.integrasi/service => ../internal/service

replace tapera.integrasi/grpc/client => ../internal/grpc/client

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-playground/validator/v10 v10.3.0
	github.com/golang/protobuf v1.4.2
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.4
	github.com/gorilla/schema v1.1.0
	github.com/joho/godotenv v1.3.0
	github.com/mwitkow/go-proto-validators v0.3.0 // indirect
	github.com/processout/grpc-go-pool v1.2.1
	github.com/rs/zerolog v1.19.0
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.6.0
	github.com/swaggo/http-swagger v0.0.0-20200308142732-58ac5e232fba
	github.com/swaggo/swag v1.6.3
	google.golang.org/grpc v1.30.0
	gopkg.in/DATA-DOG/go-sqlmock.v1 v1.3.0
	tapera.integrasi/grpc/client v1.0.0
	tapera.integrasi/service v1.0.0
	tapera/conf v1.0.0
	tapera/util v1.0.0
)
