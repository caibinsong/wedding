package utils

import (
	"fmt"
	"os"
	"sync"
)

var oneLogWriter logWriter

<<<<<<< HEAD
const MAX_SIZE int64 = 1024 * 1024 * 10 //文件大小
=======
<<<<<<< HEAD
const MAX_SIZE int64 = 1024 * 1024 * 10 //文件大小
=======
const MAX_SIZE int64 = 1024 * 1024 * 10
>>>>>>> dd12374ac95f08e4145cdb3fa4b628e5d98bd4d3
>>>>>>> a64d7c5df01427534bebc1ec23b5463de6ce4777

//获得logwriter
func GetLogWriter() *logWriter {
	return &oneLogWriter
}

//日志记录类
type logWriter struct {
	muxLog   sync.Mutex
	pLogFile *os.File
	fileName string
}

<<<<<<< HEAD
//设置日志文件保存位置
=======
<<<<<<< HEAD
//设置日志文件保存位置
=======
>>>>>>> dd12374ac95f08e4145cdb3fa4b628e5d98bd4d3
>>>>>>> a64d7c5df01427534bebc1ec23b5463de6ce4777
func (l *logWriter) SetLogFile(f string) {
	l.muxLog.Lock()
	defer l.muxLog.Unlock()
	if l.pLogFile != nil {
		l.pLogFile.Close()
		l.pLogFile = nil
	}
	l.fileName = f
	l.pLogFile, _ = os.OpenFile(f, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
}

<<<<<<< HEAD
//写日志
=======
<<<<<<< HEAD
//写日志
=======
>>>>>>> dd12374ac95f08e4145cdb3fa4b628e5d98bd4d3
>>>>>>> a64d7c5df01427534bebc1ec23b5463de6ce4777
func (l *logWriter) Write(b []byte) (int, error) {
	l.muxLog.Lock()
	defer l.muxLog.Unlock()
	fmt.Print(string(b))
	if b != nil && l.pLogFile != nil {
		size, err := l.pLogFile.Seek(0, 2)
		if err != nil {
			return 0, err
		}
		if size > MAX_SIZE {
			l.pLogFile.Close()
			os.Rename(l.fileName, l.fileName+".old")
			l.pLogFile, _ = os.OpenFile(l.fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
		}
		//这里写入到文件....
		l.pLogFile.Write(b)
		l.pLogFile.Sync()
	}
	return len(b), nil
}
