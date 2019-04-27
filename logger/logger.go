package logger

import (
	"log"
	"os"
)

//全局log
var DebugLog *log.Logger

func init() {
	//日志文件
	logFile, err := os.Create("./debug.log")
	//defer logFile.Close()
	if err != nil {
		log.Fatalln("open file error !")
	}
	DebugLog = log.New(logFile, "[Debug]", log.LstdFlags|log.Lshortfile)
}
