package cmdservice

import (
	"os"
	"log"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func debug(str ...interface{}) {
	logger.Println(str...)
}
func info(str ...interface{}) {
	logger.Println(str...)
}

func logError(str ...interface{}) {
	logger.Println(str...)
}