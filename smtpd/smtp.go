package smtpd

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

const (
	deliveryDebug = true
)

func DeliveryDebug(args ...interface{}) {
	if deliveryDebug {
		fmt.Println("deliveryDebug:")
		fmt.Println(args...)
	}
}

func DnsQuery(domain string) (string, error) {

	mx, _ := net.LookupMX(domain)
	mxLen := len(mx)

	if 0 == mxLen {
		return "", errors.New("not find mx record")
	}

	rand.Seed(time.Now().UnixNano())
	rand := rand.Intn(mxLen)
	mxSelect := mx[rand]

	mxHost := fmt.Sprintf("%s", mxSelect.Host)
	mxHost = strings.Trim(mxHost, ".")

	return mxHost, nil
}

// Delivery of mail to external mail
func Delivery(domain string, port string, from string, to string, content string) (bool, error) {

	addr := fmt.Sprintf("%s:%s", domain, port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	data, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return false, err
	}

	_, err = conn.Write([]byte("HELO IMAIL\r\n"))
	if err != nil {
		return false, err
	}

	// data, err = bufio.NewReader(conn).ReadString('\n')
	// if err != nil {
	// 	return false, err
	// }

	data = ""
	for {

		b := make([]byte, 4096)
		n, err := conn.Read(b[0:])
		fmt.Println(n, err)
		if err != nil {
			break
		}

		v := strings.TrimSpace(string(b[:n]))
		data += fmt.Sprintf("%s\r\n", v)
		inputN := strings.Split(v, "\r\n")
		last := inputN[len(inputN)-1:][0]
		fmt.Println(last)
		if strings.EqualFold(last, "250 8BITMIME") {
			break
		}

		if strings.EqualFold(last, "250 OK") {
			break
		}
	}

	fmt.Println(data)

	// if !strings.HasPrefix(data, "250") {
	// 	return false, errors.New(data)
	// }

	mailfrom := fmt.Sprintf("MAIL FROM: <%s>\r\n", from)
	DeliveryDebug(mailfrom)

	_, err = conn.Write([]byte(mailfrom))
	if err != nil {
		return false, err
	}

	data, err = bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return false, err
	}

	if !strings.HasPrefix(data, "250") {
		return false, errors.New(data)
	}

	rcpt_to := fmt.Sprintf("RCPT TO: <%s>\r\n", to)
	DeliveryDebug(rcpt_to)
	_, err = conn.Write([]byte(rcpt_to)) //向服务端发送数据。用n接受返回的数据大小，用err接受错误信息。
	if err != nil {
		return false, err
	}

	data, err = bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return false, err
	}

	if !strings.HasPrefix(data, "250") {
		return false, errors.New(data)
	}

	_, err = conn.Write([]byte("DATA\r\n")) //向服务端发送数据。用n接受返回的数据大小，用err接受错误信息。
	if err != nil {
		return false, err
	}

	data, err = bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return false, err
	}
	if !strings.HasPrefix(data, "354") {
		return false, errors.New(data)
	}

	content = fmt.Sprintf("%s\r\n\r\n", content)
	DeliveryDebug(content)
	_, err = conn.Write([]byte(content))
	if err != nil {
		return false, err
	}

	_, err = conn.Write([]byte(".\r\n"))
	if err != nil {
		return false, err
	}

	data, err = bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return false, err
	}

	fmt.Println(".", data)
	if !strings.HasPrefix(data, "354") {
		return false, errors.New(data)
	}

	_, err = conn.Write([]byte("QUIT\r\n"))
	if err != nil {
		return false, err
	}

	data, err = bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return false, err
	}

	if !strings.HasPrefix(data, "221") {
		return false, errors.New(data)
	}

	return true, nil
}
