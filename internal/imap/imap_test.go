package imap

import (
	// "crypto/tls"
	// "errors"
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/midoks/imail/internal/conf"
	"github.com/midoks/imail/internal/db"
	"github.com/midoks/imail/internal/log"
)

// go test -v ./internal/imap
func init() {
	cDir, err := os.Getwd()
	appDir := filepath.Dir(filepath.Dir(cDir))

	os.Setenv("IMAIL_WORK_DIR", appDir)
	os.Chdir(appDir)
	err = conf.Init(appDir + "/custom/conf/app.conf")
	if err != nil {
		fmt.Println("TestReceivedMail config fail:", err.Error())
	}

	conf.Web.Domain = "cachecha.com"

	logger := log.Init()

	format := conf.Log.Format
	if strings.EqualFold(format, "json") {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else if strings.EqualFold(format, "text") {
		logger.SetFormatter(&logrus.TextFormatter{})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{})
	}

	if strings.EqualFold(conf.App.RunMode, "dev") {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	go Start(10143)

	time.Sleep(1 * time.Second)
}

func imapCmd(domain string, port string, name string, password string) (bool, error) {
	addr := fmt.Sprintf("%s:%s", domain, port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return false, err
	}
	defer conn.Close()
	var content string

	data, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return false, err
	}
	fmt.Println("S:", data)

	cmd := fmt.Sprintf("B CAPABILITY\r\n")
	fmt.Println("C:", cmd)
	_, err = conn.Write([]byte(cmd))
	if err != nil {
		return false, err
	}

	for {

		b := make([]byte, 4096)
		n, err := conn.Read(b[0:])
		if err != nil {
			break
		}

		v := strings.TrimSpace(string(b[:n]))
		content += fmt.Sprintf("%s\r\n", v)
		fmt.Println("S-v:", v)
		if strings.Contains(strings.ToLower(v), "completed") {
			break
		}
	}

	cmd = fmt.Sprintf("a1 login %s %s\r\n", name, password)
	fmt.Println("C:", cmd)
	_, err = conn.Write([]byte(cmd))
	if err != nil {
		return false, err
	}

	data, err = bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return false, err
	}
	fmt.Println("S:", data)

	cmd = fmt.Sprintf("a1 LIST \"\" * \r\n")
	fmt.Println("C:", cmd)
	_, err = conn.Write([]byte(cmd))
	if err != nil {
		return false, err
	}

	for {

		b := make([]byte, 4096)
		n, err := conn.Read(b[0:])
		if err != nil {
			break
		}

		v := strings.TrimSpace(string(b[:n]))
		content += fmt.Sprintf("%s\r\n", v)
		fmt.Println("S-v:", v)

		if strings.Contains(strings.ToLower(v), "completed") {
			break
		}
	}

	cmd = fmt.Sprintf("a1 logout\r\n")
	fmt.Println("C:", cmd)
	_, err = conn.Write([]byte(cmd))
	if err != nil {
		return false, err
	}

	for {

		b := make([]byte, 4096)
		n, err := conn.Read(b[0:])
		fmt.Println(n, err)
		if err != nil {
			break
		}

		v := strings.TrimSpace(string(b[:n]))
		content += fmt.Sprintf("%s\r\n", v)
		fmt.Println("S-v:", v)
	}
	return false, err
}

// func TestRunImap163(t *testing.T) {
// 	imapCmd("imap.163.com", "143", "midoks@163.com", "mm123123")
// }

// go test -run TestRunImap
func TestRunImap(t *testing.T) {
	host := "127.0.0.1"
	port := "10143"
	name := "admin"
	password := "admin"

	addr := fmt.Sprintf("%s:%s", host, port)
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		t.Errorf("link err!")
	}

	defer conn.Close()

	cmd := fmt.Sprintf("a1 login %s %s\r\n", name, password)
	_, err = conn.Write([]byte(cmd))

	if err != nil {
		t.Errorf("user or password err!")
	}

	cmd = fmt.Sprintf("a1 select \"%s\"\r\n", "INBOX")
	_, err = conn.Write([]byte(cmd))

	if err != nil {
		t.Errorf("select err!")
	}

	cmd = "D UID FETCH 1:* (UID FLAGS)"
	_, err = conn.Write([]byte(cmd))

	if err != nil {
		t.Errorf("D UID FETCH 1:* (UID FLAGS) err!")
	}

	_, err = bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		t.Errorf("select err!")
	}
}
