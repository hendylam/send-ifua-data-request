package middleware

import (
	"net/http"
	"regexp"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"tapera/util/appcontext"
)

// ReqLoggerConfig struct
type ReqLoggerConfig struct {
	Logger *zerolog.Logger
	// UTC a boolean stating whether to use UTC time zone or local.
	UTC            bool
	SkipPath       []string
	SkipPathRegexp *regexp.Regexp
}

// ReqLogger func
func ReqLogger(config ...ReqLoggerConfig) func(http.Handler) http.Handler {
	var newConfig ReqLoggerConfig
	if len(config) > 0 {
		newConfig = config[0]
	}

	var skip map[string]struct{}
	if length := len(newConfig.SkipPath); length > 0 {
		skip = make(map[string]struct{}, length)
		for _, path := range newConfig.SkipPath {
			skip[path] = struct{}{}
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sublog := newConfig.Logger

			if sublog == nil {
				if ctxLog, ok := appcontext.LoggerWithInf(r.Context()); ok {
					sublog = ctxLog
				} else {
					sublog = &log.Logger
				}
			}

			start := time.Now()
			path := r.URL.Path
			raw := r.URL.RawQuery
			userAgnt := r.UserAgent()
			msg := "Request"

			var scheme string
			if r.TLS != nil {
				scheme = "https"
			} else {
				scheme = "http"
			}

			if raw != "" {
				path = path + "?" + raw
			}

			dumplogger := sublog.With().
				Str("scheme", scheme).
				Str("proto", r.Proto).
				Str("method", r.Method).
				Str("path", path).
				Str("ip", r.RemoteAddr).
				Str("user-agent", userAgnt).
				Logger()
			dumplogger.Info().Msg(msg)

			// continue equest
			next.ServeHTTP(w, r)

			track := true

			if _, ok := skip[path]; ok {
				track = false
			}

			if track &&
				newConfig.SkipPathRegexp != nil &&
				newConfig.SkipPathRegexp.MatchString(path) {
				track = false
			}

			if track {
				end := time.Now()
				latency := end.Sub(start)

				if newConfig.UTC {
					end = end.UTC()
				}

				msg = "Response"
				ctype := w.Header().Get("Content-Type")
				dumplogger = sublog.With().
					Str("scheme", scheme).
					Str("proto", r.Proto).
					Str("content-type", ctype).
					Str("method", r.Method).
					Str("path", path).
					Str("ip", r.RemoteAddr).
					Dur("latency", latency).
					Str("user-agent", userAgnt).
					Logger()
				dumplogger.Info().
					Msg(msg)
			}

		})
	}
}
