# ylog
go 日志库，支持日志级别和日志文件分割。


# Getting started
go get github.com/liugogt/ylog

# Example

```go
package main

import "github.com/liubogt/ylog"

func main() {
    // 2022/03/11 14:25:33.277195 [DEBUG] github.com/liubogt/ylog/main.go:7 --> This is a debug msg.
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

