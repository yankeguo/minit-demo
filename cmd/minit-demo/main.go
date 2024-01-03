package main

import (
	"flag"
	"fmt"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Logger struct {
	pfx string
	enc *encoding.Encoder
}

func (l *Logger) Println(items ...any) {
	buf := []byte(l.pfx + fmt.Sprintln(items...))
	if l.enc != nil {
		buf, _ = l.enc.Bytes(buf)
	}
	os.Stdout.Write(buf)
	os.Stdout.Sync()
}

func newLogger(gbk bool, name string) *Logger {
	l := &Logger{}
	if name != "" {
		l.pfx = "[" + name + "] "
	}
	if gbk {
		l.enc = simplifiedchinese.GBK.NewEncoder()
	}
	return l
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

	LOG := newLogger(optGBK, optName)

	var err error
	defer func() {
		if err == nil {
			LOG.Println("退出")
		} else {
			LOG.Println("错误退出:", err.Error())
		}
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
