package logger

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

const VER = 10
const DATEFORMAT = "2006-01-02"

const TIMESTAMP = "2006-01-02 15:04:05.999"

const (
	DEBUG uint8 = iota
	INFO
	WARN
	ERROR
	FATAL
	OFF
)

const (
	K int64 = 1 << 10
	M int64 = K << 10
	G int64 = M << 10
)

var logExts = [5]string{"debug", "info", "warn", "error", "fatal"}

var logFlag = [5]string{"D", "I", "W", "E", "F"}

/*
字体颜色 30:黑 31:红 32:绿 33:黄 34:蓝色 35:紫色 36:深绿
背景颜色 40:黑 41:深红 42:绿 43:黄色 44:蓝色 45:紫色 46:深绿 47:白色
显示方式 0：关闭所有属性 1：设置高亮 4：下划线 5：闪烁 7：反显 8：消隐
0x1B[背景颜色;字体颜色;显示方式m
返回正常:0x1B[0m
*/

var Color = [5]string{
	"\x1B[40;32m",
	"\x1B[40;37m",
	"\x1B[43;30;5m",
	"\x1B[41;37;1m",
	"\x1B[45;37;1m",
}
var ColorNormal = "\x1B[0m"

var logLevel uint8
var logPath, logPrefix, logFile string
var enableConsole bool = true
var logObj = make(map[uint8]*logger)
var maxFileSize int64

type logger struct {
	level     uint8
	mu        sync.Mutex
	fileIndex int
	length    int
	fp        *os.File
	fileName  string
	buf       bytes.Buffer
	fileSize  int64
	log       *log.Logger
}

func SetLevel(level uint8) {
	logLevel = level
}

func DisableConsole() {
	enableConsole = false
}

func EnableConsole() {
	enableConsole = true
}

func SetLogFile(path, prefix string, maxSize int64, unit int64) {
	logPrefix = prefix
	maxFileSize = maxSize * unit
	if !strings.HasSuffix(path, "/") {
		logPath = path + "/"
	} else {
		logPath = path
	}
	for i := logLevel; i < OFF; i++ {
		initObj(i)
	}
}

func initObj(level uint8) {
	err := mkdir(logPath)
	if err != nil {
		panic(err)
	}
	var obj *logger
	if logObj[level] == nil {
		logObj[level] = &logger{
			level: level,
		}
	}
	obj = logObj[level]
	obj.mu.Lock()
	defer obj.mu.Unlock()
	var buf bytes.Buffer
	/*
		文件名
	*/
	buf.WriteString(logPath)
	buf.WriteString(logPrefix)
	buf.WriteByte('.')
	buf.WriteString(logExts[level])
	obj.fileName = buf.String()
	fSt, err := os.Stat(obj.fileName)
	if err == nil {
		obj.fileSize = fSt.Size()
	}
	obj.fp, _ = os.OpenFile(obj.fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	obj.log = log.New(obj.fp, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func newFiles(level uint8) {
	i := 1
	obj := logObj[level]
	obj.mu.Lock()
	defer obj.mu.Unlock()
	filePrefix := obj.fileName
	for {
		if !isExist(filePrefix + "." + strconv.Itoa(i)) {
			obj.fp.Close()
			os.Rename(filePrefix, filePrefix+"."+strconv.Itoa(i))
			obj.fp, _ = os.OpenFile(filePrefix, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
			obj.log = log.New(obj.fp, "", log.Ldate|log.Ltime|log.Lshortfile)
			obj.fileSize = 0
			break
		}
		i++
	}
}

func mkdir(dir string) (e error) {
	_, er := os.Stat(dir)
	b := er == nil || os.IsExist(er)
	if !b {
		if err := os.MkdirAll(dir, 0777); err != nil {
			if os.IsPermission(err) {
				fmt.Println("permission denied:", err.Error())
				e = err
			}
		}
	}
	return
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func isFull(file string) bool {
	f, e := os.Stat(file)
	if e == nil || os.IsExist(e) {
		return false
	}
	if e != nil {
		fmt.Println(e.Error())
		return true
	}
	if f.Size() >= int64(maxFileSize) {
		return true
	}
	return false
}

func catchError() {
	if err := recover(); err != nil {
		log.Println("err", err)
	}
}

func Log(level uint8, value ...interface{}) {
	str := fmt.Sprintln(value)
	if logObj[level] != nil {
		defer catchError()
		fileCheck(level)
		logObj[level].log.Output(2, str)
		logObj[level].fileSize += int64(len(str))
	}
	/*
		show console
	*/
	if enableConsole {
		_, file, line, _ := runtime.Caller(2)
		index := strings.LastIndex(file, "/")
		short := file[index+1:]
		log.Printf("%s [%s]%s(%d)  %s %s", Color[level],
			logFlag[level], short, line, str,
			ColorNormal)
	}
}

func isNew(level uint8) bool {
	if logObj[level].fileSize >= maxFileSize {
		return true
	}
	return false
}

func fileCheck(level uint8) {
	if logObj[level] == nil {
		return
	}
	if isNew(level) {
		newFiles(level)
	}
}

func Debug(value ...interface{}) {
	Log(DEBUG, value)
}

func Info(value ...interface{}) {
	Log(INFO, value)
}

func Warn(value ...interface{}) {
	Log(WARN, value)
}

func Error(value ...interface{}) {
	Log(ERROR, value)
}

func Fatal(value ...interface{}) {
	Log(FATAL, value)
}
