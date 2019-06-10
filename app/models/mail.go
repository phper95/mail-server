package models

import (
	_ "fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

type UserMail struct {
	Id         int64
	MailFrom   string `orm:"size(50);comment(邮件来源)"`
	MailTo     string `orm:"size(50);comment(接收邮件)"`
	Content    string `orm:"comment(邮件内容)"`
	Size       int    `orm:"size(50);comment(邮件内容大小)"`
	Status     int
	UpdateTime int64
	CreateTime int64
}

func (u *UserMail) TableName() string {
	return "im_mail"
}

func (u *UserMail) Update(fields ...string) error {
	u.UpdateTime = time.Now().Unix()
	if _, err := orm.NewOrm().Update(u, fields...); err != nil {
		return err
	}
	return nil
}

func MailPush(mail_from string, mail_to string, content string) (int64, error) {
	data := new(UserMail)
	data.MailFrom = mail_from
	data.MailTo = mail_to
	data.Content = content
	data.Size = len(content)

	data.UpdateTime = time.Now().Unix()
	data.CreateTime = time.Now().Unix()
	i, err := orm.NewOrm().Insert(data)

	if err != nil {
		return 0, err
	}
	return i, err
}
