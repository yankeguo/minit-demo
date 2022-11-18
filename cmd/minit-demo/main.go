package main

import (
	"flag"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Logger struct {
	*os.File
	*log.Logger
}

func createLogger(gbk bool, name string) *Logger {
	var out io.Writer

	if gbk {
		out = os.Stdout
	} else {
		out = simplifiedchinese.GBK.NewEncoder().Writer(os.Stdout)
	}
	return &Logger{
		File:   os.Stdout,
		Logger: log.New(out, "["+name+"] ", log.LstdFlags),
	}
}

func main() {
	var (
		optName string
		optGBK  bool
		optOnce bool
	)

	flag.StringVar(&optName, "name", "noname", "set name")
	flag.BoolVar(&optGBK, "gbk", false, "set gbk")
	flag.BoolVar(&optOnce, "once", false, "set once")
	flag.Parse()

	LOG := createLogger(optGBK, optName)

	var err error
	defer func() {
		if err == nil {
			LOG.Println("exiting")
		} else {
			LOG.Println("exited with error:", err.Error())
		}
		_ = LOG.Sync()
		if err != nil {
			os.Exit(1)
		}
	}()
	defer func() {
		if e := recover(); e != nil {
			var ok bool
			if err, ok = e.(error); !ok {
				err = fmt.Errorf("%v", e)
			}
		}
	}()

	var wd string
	if wd, err = os.Getwd(); err != nil {
		return
	}

	LOG.Println("工作目录:", wd)
	LOG.Println("启动参数:", strings.Join(os.Args, ", "))
	LOG.Println("环境变量:", strings.Join(os.Environ(), ", "))

	if optOnce {
		LOG.Println("测试消息")
		time.Sleep(time.Second * 5)
	} else {
		chSig := make(chan os.Signal, 1)
		signal.Notify(chSig, syscall.SIGTERM, syscall.SIGINT)
		tick := time.NewTicker(time.Second * 3)
		for {
			select {
			case t := <-tick.C:
				LOG.Println("嘀嗒:", t.String())
			case sig := <-chSig:
				LOG.Println("信号:", sig.String())
				time.Sleep(time.Second * 3)
				return
			}
		}
	}

}
