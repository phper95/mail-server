package task

import (
	"fmt"
	"github.com/Shopify/go-rspamd"
	"github.com/midoks/imail/internal/config"
	"github.com/midoks/imail/internal/db"
	"github.com/midoks/imail/internal/libs"
	"github.com/midoks/imail/internal/smtpd"
	"github.com/robfig/cron"
	// "sync"
	"context"
	"os"
	// "time"
)

func TaskQueueeSendMail() {
	domain := config.GetString("mail.domain", "xxx.com")
	postmaster := fmt.Sprintf("postmaster@%s", domain)

	result := db.MailSendListForStatus(2, 1)
	if len(result) == 0 {

		result = db.MailSendListForStatus(0, 1)
		for _, val := range result {
			db.MailSetStatusById(val.Id, 2)
			err := smtpd.Delivery("", val.MailFrom, val.MailTo, []byte(val.Content))
			if err != nil {

				content, _ := libs.GetMailReturnToSender(val.MailFrom, val.MailTo, val.Content, err.Error())
				db.MailPush(val.Uid, 1, postmaster, val.MailFrom, content, 1)
			}
			db.MailSetStatusById(val.Id, 1)
			// fmt.Println("send mail:", err)
		}
	}
}

func TaskRspamdCheck() {
	fmt.Println("TaskRspamdCheck")
	client := rspamd.New("http://127.0.0.1:11334")
	pong, err := client.Ping(context.Background())
	fmt.Println(pong, err)

	f, err := os.Open("/Users/midoks/Desktop/t.eml")
	fmt.Println(f, err)

	email := rspamd.NewEmailFromReader(f)
	checkRes, _ := client.Check(context.Background(), email)
	fmt.Println(checkRes)
}

func Init() {
	c := cron.New()

	c.AddFunc("*/5 * * * * * ", func() {
		// fmt.Println(fmt.Sprintf("TaskQueueeSendMail! time:%d", time.Now().Unix()))
		TaskQueueeSendMail()
	})

	c.AddFunc("*/1 * * * * * ", func() {
		TaskRspamdCheck()
	})

	c.Start()
}
