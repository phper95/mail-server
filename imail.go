package main

import (
	"fmt"
	"github.com/midoks/imail/internal/app"
	"github.com/midoks/imail/internal/config"
	"github.com/midoks/imail/internal/db"
	// "github.com/midoks/imail/internal/dkim"
	"github.com/fsnotify/fsnotify"
	"github.com/midoks/imail/internal/imap"
	"github.com/midoks/imail/internal/pop3"
	"github.com/midoks/imail/internal/smtpd"
	// "log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/trace"
	"strconv"
	"strings"
)

func startService(name string) {
	config_enable := fmt.Sprintf("%s.enable", name)
	enable, err := config.GetBool(config_enable, false)
	if err == nil && enable {

		config_port := fmt.Sprintf("%s.port", name)
		port, err := config.GetInt(config_port, 25)
		if err == nil {
			fmt.Printf("listen %s port:%d success!\n", name, port)

			if strings.EqualFold(name, "smtpd") {
				go smtpd.Start(port)
			} else if strings.EqualFold(name, "pop3") {
				go pop3.Start(port)
			} else if strings.EqualFold(name, "imap") {
				go imap.Start(port)
			}
		} else {
			fmt.Printf("listen %s erorr:%s\n", name, err)
		}
	}

	config_ssl_enable := fmt.Sprintf("%s.ssl_enable", name)
	ssl_enable, err := config.GetBool(config_ssl_enable, false)
	if err == nil && ssl_enable {

		config_ssl_port := fmt.Sprintf("%s.ssl_port", name)
		ssl_port, err := config.GetInt(config_ssl_port, 25)
		if err == nil {
			fmt.Printf("listen %s ssl port:%d success!\n", name, ssl_port)

			if strings.EqualFold(name, "smtpd") {
				go smtpd.StartSSL(ssl_port)
			} else if strings.EqualFold(name, "pop3") {
				go pop3.StartSSL(ssl_port)
			} else if strings.EqualFold(name, "imap") {
				go imap.StartSSL(ssl_port)
			}
		} else {
			fmt.Printf("listen %s ssl erorr:%s\n", name, err)
		}
	}
}

func StartMonitor(path string) {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("StartMonitor:err", err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case e := <-watcher.Events:

				if e.Op&fsnotify.Chmod == fsnotify.Chmod {
					fmt.Printf("%s had change content!", path)
				}

			case err = <-watcher.Errors:
				if err != nil {
					fmt.Println("错误:", err)
				}

			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		fmt.Println("Failed to watch directory: ", err)
	}
	<-done
}

func main() {
	// go mod init
	// go mod tidy
	// go mod vendor

	err := config.Load("conf/app.conf")
	if err != nil {
		return
	}
	// err := dkim.MakeDkimConfFile("biqu.xyz")
	// fmt.Println(err)

	// tomail := "627293072@qq.com"

	// msg := []byte("from:admin@cachecha.com\r\n" +
	// 	"to: " + tomail + "\r\n" +
	// 	"Subject: hello,imail!\r\n" +
	// 	"Content-Type:multipart/mixed;boundary=a\r\n" +
	// 	"Mime-Version:1.0\r\n" +
	// 	"\r\n" +
	// 	"--a\r\n" +
	// 	"Content-type:text/plain;charset=utf-8\r\n" +
	// 	"Content-Transfer-Encoding:quoted-printable\r\n" +
	// 	"\r\n" +
	// 	"此处为正文内容D!\r\n")

	// err := smtpd.Delivery("admin@cachecha.com", tomail, msg)
	// fmt.Println("err:", err)

	// auth := smtpd.PlainAuth("", "yuludejia@gmail.com", "pmroenyllybhlwub", "smtp.gmail.com")

	// msg := []byte("from:yuludejia@gmail.com\r\n" +
	// 	"to: midoks@163.com\r\n" +
	// 	"Subject: hello,subject!\r\n" +
	// 	"Content-Type:multipart/mixed;boundary=a\r\n" +
	// 	"Mime-Version:1.0\r\n" +
	// 	"\r\n" +
	// 	"--a\r\n" +
	// 	"Content-type:text/plain;charset=utf-8\r\n" +
	// 	"Content-Transfer-Encoding:quoted-printable\r\n" +
	// 	"\r\n" +
	// 	"此处为正文内容D!\r\n")

	// err := smtpd.SendMail("smtp.gmail.com", "587", auth, "yuludejia@gmail.com", []string{"midoks@163.com"}, msg)
	// fmt.Println("err:", err)

	go pprof()

	db.Init()

	startService("smtpd")
	startService("pop3")
	startService("imap")

	http_enable, err := config.GetBool("http.enable", false)
	if http_enable {
		http_port, err := config.GetInt("http.port", 80)
		if err == nil {
			app.Start(http_port)
			fmt.Println("listen http success!")
		} else {
			fmt.Println("listen http erorr:", err)
		}

	}
	fmt.Println("end", err)
}

//手动GC
func gc(w http.ResponseWriter, r *http.Request) {
	runtime.GC()
	w.Write([]byte("StartGC"))
}

//运行trace
func traceStart(w http.ResponseWriter, r *http.Request) {
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}

	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	w.Write([]byte("TrancStart"))
	fmt.Println("StartTrancs")
}

//停止trace
func traceStop(w http.ResponseWriter, r *http.Request) {
	trace.Stop()
	w.Write([]byte("TrancStop"))
	fmt.Println("StopTrancs")
}

// go tool trace trace.out

//运行pprof分析器
func pprof() {
	go func() {
		//关闭GC
		debug.SetGCPercent(-1)
		http.HandleFunc("/go_nums", func(w http.ResponseWriter, r *http.Request) {
			num := strconv.FormatInt(int64(runtime.NumGoroutine()), 10)
			w.Write([]byte(num))
		})
		//运行trace
		http.HandleFunc("/start", traceStart)
		//停止trace
		http.HandleFunc("/stop", traceStop)
		//手动GC
		http.HandleFunc("/gc", gc)
		//网站开始监听
		http.ListenAndServe(":6060", nil)
	}()
}
