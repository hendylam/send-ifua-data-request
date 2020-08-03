package conf

import (
	"io/ioutil"
	"net"
	"os"

	"tapera/util/constant"
	"tapera/util/env"
)

// AppConf var
var AppConf AppConfig

// EnvFactory struct
type EnvFactory struct {
	factory *Factory
}

// Factory struct
type Factory struct {
	cfg *Config
}

// AppConfig interface
type AppConfig interface {
	AppName() string
	AppVer() string
	Mode() string
	Port() int
	Secret() string
	IP() string
}

// Config struct
type Config struct {
	appName string
	appVer  string
	mode    string
	port    int
	secret  string
	ip      string
}

// EnvConfig func
func EnvConfig() *Config {
	return &Config{
		appName: env.Str(constant.EnvAppName, true, nil),
		appVer:  env.Str(constant.EnvAppVer, true, nil),
		mode:    env.Str(constant.EnvAppMode, false, "debug"),
		port:    env.Int(constant.EnvAppPort, true, nil),
		secret:  env.Str(constant.EnvAppSecret, false, nil),
		ip:      findIPAddr(),
	}
}

// NewConfig function
func NewConfig(appName string, appVer string, mode string, port int, secret string) *Config {
	ip := findIPAddr()
	return &Config{
		appName: appName,
		appVer:  appVer,
		mode:    mode,
		port:    port,
		secret:  secret,
		ip:      ip,
	}
}

// NewEnvFactory func
func (c *Config) NewEnvFactory() *EnvFactory {
	return &EnvFactory{&Factory{c}}
}

// NewFactory func
func (c *Config) NewFactory() *Factory {
	return &Factory{c}
}

// AppName func
func (c *Config) AppName() string {
	return c.appName
}

// AppVer func
func (c *Config) AppVer() string {
	return c.appVer
}

// Mode func
func (c *Config) Mode() string {
	return c.mode
}

// Port func
func (c *Config) Port() int {
	return c.port
}

// Secret func
func (c *Config) Secret() string {
	return c.secret
}

// IP func
func (c *Config) IP() string {
	return c.ip
}

func findIPAddr() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func readTxtFromFile(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	return string(b), nil
}
