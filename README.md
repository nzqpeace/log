# log
A log library written by golang, based on official log pkg, support log level and partition by file size automatically

## Usage
```go
logFileName := "./log/example.log"
log.SetLogFile(logFileName)
log.SetLevel(log.LEVEL_INFO)
log.SetMaxFileSize(10 * 1024 * 1024)
log.SetFlags(log.LstdFlags | log.Lmicroseconds)

log.Debug("this is a debug message")
log.Info("test pkg log")
log.Error("error occured when create file[%s]", filename)
log.Panic("Panic! exit now")

i := 1
f := 1.345
s := "this is a test string"

log.Debugln("this is a debug message, ", s)
log.Infoln("test pkg log, ", f)
log.Errorln("error occured when create file ", filename)
log.Panicln("Panic! exit now")
```

Or create individual logger instance
------
###1. init logger instance by io.Writer
```go
f, err := os.OpenFile("t.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
if err != nil {
   fmt.Println("open file failed")
}
logger := log.NewLogger(f)
logger.SetFlags(log.LstdFlags | log.Lmicroseconds)
```
###2. init logger instance by filename
```go
logger, err := log.NewFileLog(logfilename)
if err != nil{
    fmt.Println(err)
}
logger.Info("create logger instance success!")
```
