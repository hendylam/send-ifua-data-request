package constant

const (
	// EnvAppName env const
	EnvAppName = "APP_NAME"
	// EnvAppVer env const
	EnvAppVer = "APP_VER"
	// EnvAppPort env const
	EnvAppPort = "APP_PORT"
	// EnvAppHost env const
	EnvAppHost = "APP_HOST"
	// EnvAppMode env const
	EnvAppMode = "APP_MODE"
	// EnvAppSecret env const
	EnvAppSecret = "APP_SECRET"
)

type correlationIDType int

const (
	// RequestIDKey is the key that holds the unique request ID in a request context.
	RequestIDKey correlationIDType = iota
	// ContextLoggerKey is the key that holds the logger with reqId info in a request context.
	ContextLoggerKey
	// UserKey const
	UserKey
)
