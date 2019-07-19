package log

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
)

/**
根据传入的参数设定日志级别
*/
func convertLogLevel(level string) int {
	switch (level) {
	case "debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "trace":
		return logs.LevelTrace
	}
	return logs.LevelDebug
}

/**
日志初始化
*/
func InitLogger(LogPath string, LogLevel string) (err error) {
	config := make(map[string]interface{})
	config["fileName"] = LogPath
	config["level"] = convertLogLevel(LogLevel)
	configStr, err := json.Marshal(config)
	if err != nil {
		fmt.Println("initLogger failed , err", err)
		return
	}
	_ = logs.SetLogger(logs.AdapterConsole)
	_ = logs.SetLogger(logs.AdapterFile, string(configStr))
	return
}
