package zlog

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zapcore"
	"time"
)

var (
	GinOutPut *LoggerWriter
)

func init() {
	GinOutPut = NewWriter(DefaultLogger, zapcore.DebugLevel)
}
func GetGinLogger(conf ...gin.LoggerConfig) gin.HandlerFunc {
	if len(conf) == 0 {
		return gin.LoggerWithConfig(gin.LoggerConfig{
			Formatter: LogFormatter,
			Output:    GinOutPut,
		})
	}
	return gin.LoggerWithConfig(conf[0])
}

var LogFormatter = func(param gin.LogFormatterParams) string {
	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}
	data := fmt.Sprintf("[GIN] statusCode:%v path:%s method:%s cost:%v clientIp:%s", param.StatusCode, param.Path, param.Method, param.Latency, param.ClientIP)
	return data
}
