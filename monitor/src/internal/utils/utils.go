package utils

import (
	"log"

	"github.com/davecgh/go-spew/spew"
)

const (
	lError = iota
	lWarn
	lInfo
	lDebug
	lDump
)

var Level = 0

func SetLevel(l int) {
	Level = l
}

func Dump(n string, v interface{}) {
	if Level >= lDump {
		log.Printf("<<<<<<<<<<<<<<<<<< DUMP %s START >>>>>>>>>>>>>>>>>>>>>", n)
		spew.Dump(v)
		log.Printf("<<<<<<<<<<<<<<<<<< DUMP %s END >>>>>>>>>>>>>>>>>>>>>>>", n)
	}

}

func Warn(s string, a ...interface{}) {
	if Level >= lWarn {
		log.Printf(s, a...)
	}
}

func Info(s string, a ...interface{}) {
	if Level >= lInfo {
		log.Printf(s, a...)
	}
}

func Debug(s string, a ...interface{}) {
	if Level >= lDebug {
		log.Printf(s, a...)
	}
}

func Error(s string, a ...interface{}) {
	if Level >= lError {
		log.Printf(s, a...)
	}
}
