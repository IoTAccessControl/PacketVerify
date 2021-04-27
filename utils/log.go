package utils

import (
	"fmt"
	"log"
	"os"
)

const LOGDIR = "log"

// https://www.honeybadger.io/blog/golang-logging/
func InitLogger(tag string)  {
	name := LOGDIR + "/" + tag + ".log"
	file, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	log.SetOutput(file)
	log.Println("======================================")
}

// 用来评估计时
func Statistic(f string, args ... interface{}) {
	f = "[STATISTIC] " + f
	fmt.Printf(f, args...)
	log.Printf(f, args...)
}

func LogInfo(f string, args ... interface{})  {
	fmt.Printf(f, args...)
	log.Printf(f, args...)
}

func Fatal(s string, a ... interface{}) {
	fmt.Fprintf(os.Stderr, "netfwd: %s\n", fmt.Sprintf(s, a))
	os.Exit(2)
}

