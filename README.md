# ylog
ylog是一个轻量级的go语言日志库，输出打印日志行所在的包路径，而非一般日志库的文件路径。
* 支持日志级别打印和全局日志级别设置。
* 支持按时间分隔日志文件。
* 打印日志所在行的包路径。
* 日志格式模板：时间 级别 包路径 --> 消息内容

# Install
go get github.com/liugogt/ylog

# Example

```go
package main

import "github.com/liubogt/ylog"

func main() {
    // 2022-03-11 14:25:33.277 [DEBUG] github.com/liubogt/ylog/main.go:7 --> This is a debug msg.
    ylog.Debug("This is a debug msg.")
    ylog.Debug("This is a debug msg,", "Hello.")
    ylog.DebugF("This is a debug msg: %s.", "Hi")
    
    
    // ylog.Info("This is a  info msg.")
    // ...
    
    // ylog.Warn("This is a  warn msg.")
    // ...
    
    // ylog.Error("This is a  error msg.")
    // ...    
    
}

func init() {
    // set log level
    ylog.SetLevel(ylog.LevelDebug)
}
```

