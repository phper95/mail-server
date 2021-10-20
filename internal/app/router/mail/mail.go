package mail

import (
	"fmt"

	"github.com/midoks/imail/internal/app/context"
	"github.com/midoks/imail/internal/app/form"
	"github.com/midoks/imail/internal/db"
	tmail "github.com/midoks/imail/internal/tools/mail"
	"github.com/midoks/imail/internal/tools/paginater"
)

const (
	MAIL        = "mail/list"
	MAIL_NEW    = "mail/new"
	MAIL_CONENT = "mail/content"
)

type MailSearchOptions struct {
	page     int
	PageSize int
	OrderBy  string
	TplName  string
	Type     int
}

func RenderMailSearch(c *context.Context, opts *MailSearchOptions) {
	page := c.QueryInt("page")
	if page <= 1 {
		page = 1
	}

	var (
		mail  []*db.Mail
		count int64
		err   error
	)

	keyword := c.Query("q")
	opt := &db.MailSearchOptions{
		Keyword:  keyword,
		OrderBy:  opts.OrderBy,
		Page:     page,
		PageSize: opts.PageSize,
		Type:     opts.Type,
		Uid:      c.User.Id,
	}

	if len(keyword) == 0 {
		mail, err = db.MailList(page, opts.PageSize, opt)
		count = db.MailCountWithOpts(opt)
	} else {
		mail, count, err = db.MailSearchByName(opt)
		if err != nil {
			c.Error(err, "search user by name")
			return
		}
	}
	c.Data["Keyword"] = keyword
	c.Data["Total"] = count
	c.Data["Page"] = paginater.New(int(count), opts.PageSize, page, 5)
	c.Data["Mail"] = mail

	c.Success(opts.TplName)
}

func Flags(c *context.Context) {
	c.Data["Title"] = c.Tr("mail.flags")
	c.Data["PageIsMail"] = true

	RenderMailSearch(c, &MailSearchOptions{
		PageSize: 10,
		OrderBy:  "id ASC",
		TplName:  MAIL,
		Type:     db.MailSearchOptionsTypeFlags,
	})
}

func Sent(c *context.Context) {
	c.Data["Title"] = c.Tr("mail.sent")
	c.Data["PageIsMail"] = true

	RenderMailSearch(c, &MailSearchOptions{
		PageSize: 10,
		OrderBy:  "id ASC",
		TplName:  MAIL,
		Type:     db.MailSearchOptionsTypeSend,
	})
}

func Deleted(c *context.Context) {
	c.Data["Title"] = c.Tr("mail.deleted")
	c.Data["PageIsMail"] = true

	RenderMailSearch(c, &MailSearchOptions{
		PageSize: 10,
		OrderBy:  "id ASC",
		TplName:  MAIL,
		Type:     db.MailSearchOptionsTypeDeleted,
	})
}

func Junk(c *context.Context) {
	c.Data["Title"] = c.Tr("mail.junk")
	c.Data["PageIsMail"] = true

	RenderMailSearch(c, &MailSearchOptions{
		PageSize: 10,
		OrderBy:  "id ASC",
		TplName:  MAIL,
		Type:     db.MailSearchOptionsTypeJunk,
	})
}

func Mail(c *context.Context) {
	c.Data["Title"] = c.Tr("mail.write_letter")
	c.Data["PageIsMail"] = true

	c.Success(MAIL)
}

func New(c *context.Context) {
	c.Data["Title"] = c.Tr("mail.write_letter")
	c.Data["PageIsWriteMail"] = true

	c.Success(MAIL_NEW)
}

func NewPost(c *context.Context, f form.SendMail) {
	c.Data["Title"] = c.Tr("mail.write_letter")
	c.Data["PageIsWriteMail"] = true

	from, err := db.DomainGetMainForDomain()
	if err != nil {
		c.RenderWithErr(c.Tr("mail.new.default_not_set"), MAIL_NEW, &f)
		return
	}

	mail_from := fmt.Sprintf("%s@%s", c.User.Name, from)
	tc := tmail.GetMailRealContent(mail_from, f.ToMail, f.Subject, f.Content)

	_, err = db.MailPushSend(c.User.Id, mail_from, f.ToMail, tc)
	if err != nil {
		c.RenderWithErr(err.Error(), MAIL_NEW, &f)
		return
	}
	c.Flash.Success("OK")
	c.Success(MAIL_NEW)
}

func Content(c *context.Context) {
	c.Data["Title"] = c.Tr("mail.write_letter")
	c.Data["PageIsMailContent"] = true

	id := c.ParamsInt64(":id")

	fmt.Println(id)

	c.Success(MAIL_CONENT)
}
