package gini

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gokit/zapi"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"
)

func Logger(logger *zapi.Logger) HandlerFunc {
	return func(c *Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		logger.WithContext(c).Info(c.Request.URL.Path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("ip", c.ClientIP()),
			zap.String("ua", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("duration", duration),
		)
	}
}

func Recovery(logger *zapi.Logger) HandlerFunc {
	return func(c *Context) {
		defer func() {
			logger = logger.WithContext(c)
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				headers := strings.Split(string(httpRequest), "\r\n")
				for idx, header := range headers {
					current := strings.Split(header, ":")
					if current[0] == "Authorization" {
						headers[idx] = current[0] + ": *"
					}
				}
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				} else {
					logger.Error("recover panic",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}

				if brokenPipe {
					c.Error(err.(error))
					c.Abort()
				} else {
					c.AbortWithStatus(http.StatusInternalServerError)
				}
			}
		}()
		c.Next()
	}
}
