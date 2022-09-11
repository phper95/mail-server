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

	"github.com/phper95/mail-server/internal/conf"
	"github.com/phper95/mail-server/internal/db"
	"github.com/phper95/mail-server/internal/log"
	"github.com/phper95/mail-server/internal/tools"
)

// go test -v ./internal/go test -v ./internal/imap
func init() {

	cDir, err := os.Getwd()
	appDir := filepath.Dir(filepath.Dir(cDir))

	os.Setenv("IMAIL_WORK_DIR", appDir)
	os.Chdir(appDir)

	if tools.IsExist(appDir + "/custom/conf/app.conf") {
		err = conf.Init(appDir + "/custom/conf/app.conf")
		if err != nil {
			fmt.Println("test init config fail:", err.Error())
			return
		}
	} else {
		err = conf.Init("")
		if err != nil {
			fmt.Println("test init config fail:", err.Error())
			return
		}
	}

}

func initDb() {
	conf.Log.RootPath = conf.WorkDir() + "/logs"
	os.MkdirAll(conf.Log.RootPath, os.ModePerm)
	//conf.Database.Type = "sqlite3"
	//conf.Database.Path = "data/imail.db3"
	conf.Database.Type = "mysql"
	conf.Database.Host = "127.0.0.1:3306"
	conf.Database.User = "root"
	conf.Database.Password = "admin123"
	conf.Database.Name = "mail"
	conf.Database.Charset = "utf8mb4"

	conf.Smtp.Debug = false

	log.Init()
	db.Init()

	// create default user
	db.CreateUser(&db.User{
		Name:     "admin",
		Password: "admin123",
		Salt:     "123123",
		Code:     "admin",
	})

	d := &db.Domain{
		Domain:    "tt.com",
		Mx:        true,
		A:         true,
		Spf:       true,
		Dkim:      true,
		Dmarc:     true,
		IsDefault: true,
	}

	err := db.DomainCreate(d)
	if err != nil {
		return
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
// 	imapCmd("imap.163.com", "143", "phper95@163.com", "mm123123")
// }

// go test -run TestRunImap
// go test -v ./internal/imap -bench TestRunImap
func TestRunImap(t *testing.T) {
	initDb()

	host := "127.0.0.1"
	port := "10143"
	name := "admin"
	password := "admin123"

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
