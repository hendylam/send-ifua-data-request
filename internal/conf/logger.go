package conf

import (
	"io"
	"os"
	"path"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"

	"tapera/util/env"
)

const (
	envLogConsoleEnbl  = "LOG_CONSOLE_ENABLE"
	envLogEncodedAsJSN = "LOG_ENCODED_AS_JSON"
	envLogFileEnbl     = "LOG_FILE_ENABLE"
	envLogFileDir      = "LOG_FILE_DIR"
	envLogFileNm       = "LOG_FILE_NAME"
	envLogFileMaxSz    = "LOG_FILE_MAX_SIZE"
	envLogFileMaxBc    = "LOG_FILE_MAX_BACKUP"
	envLogFileMaxAg    = "LOG_FILE_MAX_AGE"
	envLogCompressed   = "LOG_COMPRESSED"
)

// LoggerConfig struct
type LoggerConfig struct {
	// Enable console logging
	ConsoleLoggingEnabled bool `json:"consoleloggingenabled" yaml:"consoleloggingenabled"`

	// EncodeLogsAsJson makes the log framework log JSON
	EncodeLogsAsJSON bool `json:"encodelogsasjson" yaml:"encodelogsasjson"`

	// FileLoggingEnabled makes the framework log to a file
	// the fields below can be skipped if this value is false!
	FileLoggingEnabled bool `json:"fileloggingenabled" yaml:"fileloggingenabled"`

	// Directory to log to to when filelogging is enabled
	Directory string `json:"directory" yaml:"directory"`
	// Filename is the name of the logfile which will be placed inside the directory

	Filename string `json:"filename" yaml:"filename"`
	// MaxSize the max size in MB of the logfile before it's rolled

	MaxSize int `json:"maxsizes" yaml:"maxsizes"`
	// MaxBackups the max number of rolled files to keep

	MaxBackups int `json:"maxbackups" yaml:"maxbackups"`

	// MaxAge the max age in days to keep a logfile
	MaxAge int

	Compress bool `json:"compress" yaml:"compress"`
}

// Logger func
func (f *Factory) Logger(cfg *LoggerConfig) *zerolog.Logger {
	//time.RFC3339
	//zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if f.cfg.mode == "debug" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	var writers []io.Writer

	if cfg.ConsoleLoggingEnabled {
		if cfg.EncodeLogsAsJSON {
			writers = append(writers, log.Output(zerolog.ConsoleWriter{Out: os.Stdout}))
		} else {
			writers = append(writers, zerolog.ConsoleWriter{
				Out:     os.Stderr,
				NoColor: false,
			})
		}
	}

	if cfg.FileLoggingEnabled {
		writers = append(writers, newRollingFile(cfg))
	}

	mw := io.MultiWriter(writers...)
	logger := zerolog.New(mw).With().Timestamp().Str("app", f.cfg.appName+f.cfg.appVer).Str("host", f.cfg.ip+":"+strconv.Itoa(f.cfg.port)).Logger()

	logger.Info().
		Bool("fileLogging", cfg.FileLoggingEnabled).
		Bool("jsonLogOutput", cfg.EncodeLogsAsJSON).
		Str("logDirectory", cfg.Directory).
		Str("fileName", cfg.Filename).
		Int("maxSizeMB", cfg.MaxSize).
		Int("maxBackups", cfg.MaxBackups).
		Int("maxAgeInDays", cfg.MaxAge).
		Bool("compress", cfg.Compress).
		Msg("logging configured")

	return &logger
}

// Logger func
func (f *EnvFactory) Logger() *zerolog.Logger {
	cfg := &LoggerConfig{
		ConsoleLoggingEnabled: env.Bool(envLogConsoleEnbl, false, true),
		EncodeLogsAsJSON:      env.Bool(envLogEncodedAsJSN, false, true),
		FileLoggingEnabled:    env.Bool(envLogFileEnbl, false, nil),
	}

	if cfg.FileLoggingEnabled {
		cfg.Directory = env.Str(envLogFileDir, true, nil)
		cfg.Filename = env.Str(envLogFileNm, true, nil)
		cfg.MaxSize = env.Int(envLogFileMaxSz, true, nil)
		cfg.MaxBackups = env.Int(envLogFileMaxBc, true, nil)
		cfg.MaxAge = env.Int(envLogFileMaxAg, true, nil)
		cfg.Compress = env.Bool(envLogCompressed, false, nil)
	}
	return f.factory.Logger(cfg)
}

func newRollingFile(cfg *LoggerConfig) io.Writer {
	if err := os.MkdirAll(cfg.Directory, 0744); err != nil {
		log.Error().Err(err).Str("path", cfg.Directory).Msg("can't create log directory")
		return nil
	}

	return &lumberjack.Logger{
		Filename:   path.Join(cfg.Directory, cfg.Filename),
		MaxBackups: cfg.MaxBackups, // files
		MaxSize:    cfg.MaxSize,    // megabytes
		MaxAge:     cfg.MaxAge,     // days
		Compress:   cfg.Compress,   // compress
	}
}
