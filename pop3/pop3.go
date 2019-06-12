package pop3

import (
	"bufio"
	"fmt"
	"github.com/midoks/imail/app/models"
	"github.com/midoks/imail/libs"
	"log"
	"net"
	"runtime"
	"strings"
	"time"
)

const (
	CMD_READY = iota
	CMD_USER  = iota
	CMD_PASS  = iota
	CMD_STAT  = iota
	CMD_LIST  = iota
	CMD_RETR  = iota
	CMD_DELE  = iota
	CMD_NOOP  = iota
	CMD_RSET  = iota
	CMD_TOP   = iota
	CMD_UIDL  = iota
	CMD_APOP  = iota
	CMD_QUIT  = iota
	CMD_CAPA  = iota
)

var stateList = map[int]string{
	CMD_READY: "READY",
	CMD_USER:  "USER",
	CMD_PASS:  "PASS",
	CMD_STAT:  "STAT",
	CMD_LIST:  "LIST",
	CMD_RETR:  "RETR",
	CMD_DELE:  "DELE",
	CMD_NOOP:  "NOOP",
	CMD_RSET:  "RSET",
	CMD_TOP:   "TOP",
	CMD_UIDL:  "UIDL",
	CMD_APOP:  "APOP",
	CMD_QUIT:  "QUIT",
	CMD_CAPA:  "CAPA",
}

const (
	MSG_INIT          = "Welcome to coremail Mail Pop3 Server (imail)"
	MSG_OK            = "core mail"
	MSG_BAD_SYNTAX    = "500"
	MSG_LOGIN_OK      = "%d message(s) [%d byte(s)]"
	MSG_LOGIN_DISABLE = "Unable to log on"
	MSG_CMD_NOT_VALID = "Command not valid in this state"
	MSG_RETR_DATA     = "%s\r\n."
	MSG_CAPA          = "Capability list follows"
)

var GO_EOL = GetGoEol()

func GetGoEol() string {
	if "windows" == runtime.GOOS {
		return "\r\n"
	}
	return "\n"
}

type Pop3Server struct {
	debug         bool
	conn          net.Conn
	state         int
	startTime     time.Time
	errCount      int
	recordCmdUser string
	recordCmdPass string

	// user id
	userID int64
}

func (this *Pop3Server) setState(state int) {
	this.state = state
}

func (this *Pop3Server) getState() int {
	return this.state
}

func (this *Pop3Server) D(a ...interface{}) (n int, err error) {
	return fmt.Println(a...)
}

func (this *Pop3Server) Debug(d bool) {
	this.debug = d
}

func (this *Pop3Server) w(msg string) {
	_, err := this.conn.Write([]byte(msg))

	if err != nil {
		log.Fatal(err)
	}
}

func (this *Pop3Server) writeArgs(code string, args ...interface{}) {
	info := fmt.Sprintf("+OK "+code+"\r\n", args...)
	this.w(info)
}

func (this *Pop3Server) ok(code string) {
	info := fmt.Sprintf("+OK %s\r\n", code)
	this.w(info)
}

func (this *Pop3Server) error(code string) {
	info := fmt.Sprintf("-ERR %s\r\n", code)
	this.w(info)
}

func (this *Pop3Server) getString() (string, error) {
	input, err := bufio.NewReader(this.conn).ReadString('\n')
	if err != nil {
		return "", err
	}
	inputTrim := strings.TrimSpace(input)
	return inputTrim, err
}

func (this *Pop3Server) getString0() (string, error) {
	buffer := make([]byte, 2048)

	n, err := this.conn.Read(buffer)
	if err != nil {
		log.Fatal(this.conn.RemoteAddr().String(), " connection error: ", err)
		return "", err
	}

	input := string(buffer[:n])
	inputTrim := strings.TrimSpace(input)
	return inputTrim, err
}

func (this *Pop3Server) close() {
	this.conn.Close()
}

func (this *Pop3Server) cmdCompare(input string, cmd int) bool {
	if strings.EqualFold(input, stateList[cmd]) {
		return true
	}
	return false
}

func (this *Pop3Server) stateCompare(input int, cmd int) bool {
	if input == cmd {
		return true
	}
	return false
}

func (this *Pop3Server) checkUserLogin() bool {
	name := this.recordCmdUser
	pwd := this.recordCmdPass

	fmt.Println(name, pwd)
	info, err := models.UserGetByName(name)
	if err != nil {
		return false
	}

	pwdMd5 := libs.Md5str(pwd)
	if pwdMd5 != info.Password {
		return false
	}

	this.userID = info.Id
	return true
}

func (this *Pop3Server) cmdUser(input string) bool {
	inputN := strings.SplitN(input, " ", 2)

	if this.cmdCompare(inputN[0], CMD_USER) {
		if len(inputN) < 2 {
			this.ok(MSG_BAD_SYNTAX)
			return false
		}

		this.recordCmdUser = strings.TrimSpace(inputN[1])
		this.ok(MSG_OK)
		return true
	}
	return false
}

func (this *Pop3Server) cmdPass(input string) bool {
	inputN := strings.SplitN(input, " ", 2)

	if this.cmdCompare(inputN[0], CMD_PASS) {
		if len(inputN) < 2 {
			this.ok(MSG_BAD_SYNTAX)
			return false
		}
		this.recordCmdPass = strings.TrimSpace(inputN[1])

		if this.checkUserLogin() {

			count, size := models.BoxUserTotal(this.userID)
			fmt.Println("stat", count, size, this.userID)
			this.writeArgs(MSG_LOGIN_OK, count, size)
			return true
		}
		this.error(MSG_LOGIN_DISABLE)
		return false
	}
	return false
}

func (this *Pop3Server) cmdStat(input string) bool {
	if this.cmdCompare(input, CMD_STAT) {

		count, size := models.BoxUserTotal(this.userID)
		fmt.Println("stat", count, size, this.userID)
		this.writeArgs(MSG_LOGIN_OK, count, size)
		return true
	}
	return false
}

func (this *Pop3Server) cmdList(input string) bool {
	inputN := strings.SplitN(input, " ", 2)

	inputLen := len(inputN)
	if inputLen == 1 {
		if this.cmdCompare(inputN[0], CMD_LIST) {

			t1, t2 := models.BoxUserTotal(this.userID)
			fmt.Println(t1, t2)

			list, count := models.BoxList(this.userID, 1, 100)
			fmt.Println(list, count)
			return true
		}
	} else if inputLen == 2 {
		if this.cmdCompare(inputN[0], CMD_LIST) {
			return true
		}
	}
	this.error(MSG_BAD_SYNTAX)
	return false
}

func (this *Pop3Server) cmdRetr(input string) bool {
	inputN := strings.SplitN(input, " ", 2)

	if this.cmdCompare(inputN[0], CMD_RETR) {
		if len(inputN) < 2 {
			this.ok(MSG_BAD_SYNTAX)
			return false
		}
		this.ok(MSG_RETR_DATA)
		return true
	}
	return false
}

func (this *Pop3Server) cmdUidl(input string) bool {
	inputN := strings.SplitN(input, " ", 2)

	if this.cmdCompare(inputN[0], CMD_UIDL) {
		if len(inputN) < 2 {
			this.ok(MSG_BAD_SYNTAX)
			return false
		}
		this.ok(MSG_OK)
		return true
	}
	return false
}

func (this *Pop3Server) cmdQuit(input string) bool {
	if this.cmdCompare(input, CMD_QUIT) {
		this.ok(MSG_OK)
		this.close()
		return true
	}
	return false
}

func (this *Pop3Server) cmdCapa(input string) bool {
	if this.cmdCompare(input, CMD_CAPA) {
		this.ok(MSG_CAPA)
		this.w("TOP\r\n")
		this.w("USER\r\n")
		this.w("PIPELINING\r\n")
		this.w("UIDL\r\n")
		this.w("LANG\r\n")
		this.w("UTF8\r\n")
		this.w("SASL PLAIN\r\n")
		this.w("STLS\r\n")
		this.w(".\r\n")
		return true
	}
	return false
}

func (this *Pop3Server) handle() {
	for {
		state := this.getState()
		input, err := this.getString()

		if err != nil {
			break
		}

		fmt.Println("pop3:", state, input)
		if this.cmdQuit(input) {
			break
		}

		if this.cmdCapa(input) {
		}

		if this.stateCompare(state, CMD_READY) {
			if this.cmdUser(input) {
				this.setState(CMD_USER)
			}
		}

		if this.stateCompare(state, CMD_USER) {
			if this.cmdPass(input) {
				this.setState(CMD_PASS)
			}
		}

		if this.stateCompare(state, CMD_PASS) {

			if this.cmdStat(input) {
			}

			if this.cmdList(input) {
			}

			if this.cmdUidl(input) {
			}

			if this.cmdRetr(input) {
			}
		}
	}
}

func (this *Pop3Server) start(conn net.Conn) {
	conn.SetReadDeadline(time.Now().Add(time.Minute * 10))
	defer conn.Close()
	this.conn = conn

	this.startTime = time.Now()

	this.ok(MSG_INIT)
	this.setState(CMD_READY)

	this.handle()
}

func Start(port int) {
	addr := fmt.Sprintf(":%d", port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		srv := Pop3Server{}
		go srv.start(conn)
	}
}
