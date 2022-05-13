package log

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

const (
	InfoLevel = iota
	DebugLevel
	ErrorLevel
	Disabled
)

var (
	infoLog  = log.New(os.Stdout, "\033[34m[info]\033[0m ", log.LstdFlags|log.Lshortfile)  // blue "\033[34m text \033[0m"
	debugLog = log.New(os.Stdout, "\033[33m[debug]\033[0m ", log.LstdFlags|log.Lshortfile) // yellow
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile) // red
	loggers  = []*log.Logger{infoLog, debugLog, errorLog}                                  // ordered by level ascending
	mu       sync.Mutex
)

// log methods
var (
	Info   = infoLog.Println
	Infof  = infoLog.Printf
	Debug  = debugLog.Println
	Debugf = debugLog.Printf
	Error  = errorLog.Println
	Errorf = errorLog.Printf
)

func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()

	for i := InfoLevel; i < Disabled; i++ {
		if i >= level {
			loggers[i].SetOutput(os.Stdout)
		} else {
			loggers[i].SetOutput(ioutil.Discard)
		}
	}

}
