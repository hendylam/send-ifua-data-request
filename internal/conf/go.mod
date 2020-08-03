module tapera.integrasi/conf

go 1.14

replace tapera/util => ../util

require (
	tapera/util v1.0.0
	github.com/jinzhu/gorm v1.9.12
	github.com/lib/pq v1.1.1
	github.com/rs/zerolog v1.19.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)
