package cmd

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/midoks/imail/internal/app"
	"github.com/midoks/imail/internal/config"
	"github.com/midoks/imail/internal/db"
	"github.com/midoks/imail/internal/debug"
	"github.com/midoks/imail/internal/imap"
	"github.com/midoks/imail/internal/log"
	"github.com/midoks/imail/internal/pop3"
	"github.com/midoks/imail/internal/smtpd"
	"github.com/midoks/imail/internal/task"
	"github.com/urfave/cli"
	"strings"
)

var Service = cli.Command{
	Name:        "service",
	Usage:       "This command starts all services",
	Description: `Start POP3, IMAP, SMTP, web and other services`,
	Action:      runAllService,
	Flags: []cli.Flag{
		stringFlag("config, c", "", "Custom configuration file path"),
	},
}

func runAllService(c *cli.Context) error {

	confFile, err := initConfig(c, "")
	if err != nil {
		panic("imail config file load error")
		return err
	}

	log.Init()
	go ConfigFileStartMonitor(confFile)

	err = db.Init()
	if err != nil {
		return err
	}

	task.Init()

	runmode := config.GetString("runmode", "dev")
	if strings.EqualFold(runmode, "dev") {
		go debug.Pprof()
	}

	startService("smtpd")
	startService("pop3")
	startService("imap")

	http_enable, err := config.GetBool("http.enable", false)
	if http_enable {
		http_port, err := config.GetInt("http.port", 80)
		if err == nil {
			log.Info("listen http success!")
			app.Start(http_port)
		} else {
			log.Errorf("listen http erorr:%s", err)
		}
	}
	return nil
}

func ServiceDebug() {

	confFile, err := initConfig(nil, "conf/app.conf")
	if err != nil {
		panic("imail config file load error")
	}

	log.Init()

	go ConfigFileStartMonitor(confFile)

	err = db.Init()
	if err != nil {
		return
	}

	task.Init()

	runmode := config.GetString("runmode", "dev")

	if strings.EqualFold(runmode, "dev") {
		go debug.Pprof()
	}

	startService("smtpd")
	startService("pop3")
	startService("imap")

	http_enable, err := config.GetBool("http.enable", false)
	if http_enable {
		http_port, err := config.GetInt("http.port", 80)
		if err == nil {
			log.Infof("listen http[%d] success!", http_port)
			app.Start(http_port)
		} else {
			log.Errorf("listen http[%d] erorr:%s", http_port, err)
		}
	}
}

func startService(name string) {
	config_enable := fmt.Sprintf("%s.enable", name)
	enable, err := config.GetBool(config_enable, false)
	if err == nil && enable {

		config_port := fmt.Sprintf("%s.port", name)
		port, err := config.GetInt(config_port, 25)
		if err == nil {
			log.Infof("listen %s port:%d success!", name, port)

			if strings.EqualFold(name, "smtpd") {
				go smtpd.Start(port)
			} else if strings.EqualFold(name, "pop3") {
				go pop3.Start(port)
			} else if strings.EqualFold(name, "imap") {
				go imap.Start(port)
			}
		} else {
			log.Errorf("listen %s erorr:%s", name, err)
		}
	}

	config_ssl_enable := fmt.Sprintf("%s.ssl_enable", name)
	ssl_enable, err := config.GetBool(config_ssl_enable, false)
	if err == nil && ssl_enable {

		config_ssl_port := fmt.Sprintf("%s.ssl_port", name)
		ssl_port, err := config.GetInt(config_ssl_port, 25)
		if err == nil {
			log.Infof("listen %s ssl port:%d success!", name, ssl_port)

			if strings.EqualFold(name, "smtpd") {
				go smtpd.StartSSL(ssl_port)
			} else if strings.EqualFold(name, "pop3") {
				go pop3.StartSSL(ssl_port)
			} else if strings.EqualFold(name, "imap") {
				go imap.StartSSL(ssl_port)
			}
		} else {
			log.Errorf("listen %s ssl erorr:%s", name, err)
		}
	}
}

func reloadService(path string) {
	log.Infof("reloadService:%s", path)

	err := config.Load(path)
	if err != nil {
		log.Errorf("imail config file reload error:%s", err)
		return
	}

	// fmt.Println("imap reload start")
	// err = imap.Close()
	// fmt.Println("[reloadService]err", err)
	// fmt.Println("startService")
	// startService("imap")
	// fmt.Println("imap reload end")

}

func ConfigFileStartMonitor(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Errorf("fsnotify.NewWatcher errors:%s", err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case e := <-watcher.Events:

				if e.Op&fsnotify.Chmod == fsnotify.Chmod {
					reloadService(path)
				}
			case err = <-watcher.Errors:
				if err != nil {
					log.Errorf("watcher errors:%s", err)
				}
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		log.Errorf("failed to watch directory error:%s", err)
	}
	<-done
}
