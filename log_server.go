package zlog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
)

func InitLogServer(port int32) {
	go func() {
		http.HandleFunc("/updateLevel", func(w http.ResponseWriter, r *http.Request) {
			// 解析查询参数
			queryParams := r.URL.Query()
			// 获取单个参数
			level := queryParams.Get("level")
			l, err := zapcore.ParseLevel(level)
			if err != nil {
				w.Write([]byte(`{"code":400,"message":"parse level failed"}`))
				return
			}

			UpdateLoggerLevel(l)
		})

		if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil); err != nil {
			zap.S().Error("init log server start failed", zap.Error(err))
		}
	}()
}
