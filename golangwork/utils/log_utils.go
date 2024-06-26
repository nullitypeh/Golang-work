package utils

import (
	"log"
	"os"
)

var Logger *log.Logger

func InitLogger(logFilePath string) {
	// 创建日志文件
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	// 创建 Logger
	Logger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Logger.Println("Logger initialized")
}
