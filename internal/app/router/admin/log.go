package admin

import (
	// "strings"

	"github.com/phper95/mail-server/internal/app/context"
	// "github.com/phper95/mail-server/internal/app/form"
	"github.com/phper95/mail-server/internal/db"
	"github.com/phper95/mail-server/internal/tools/paginater"
)

const (
	MAIL_LOG_LIST = "admin/log/list"
)

func RenderLogSearch(c *context.Context, opts *db.LogSearchOptions) {
	page := c.QueryInt("page")
	if page <= 1 {
		page = 1
	}

	var (
		log   []db.MailLog
		count int64
		err   error
	)

	keyword := c.Query("q")
	if len(keyword) == 0 {
		log, err = db.LogList(page, opts.PageSize)
		count = db.LogCount()
	} else {
		log, count, err = db.LogSearchByName(&db.LogSearchOptions{
			Keyword:  keyword,
			OrderBy:  opts.OrderBy,
			Page:     page,
			PageSize: opts.PageSize,
		})
		if err != nil {
			c.Error(err, "search user by name")
			return
		}
	}
	c.Data["Keyword"] = keyword
	c.Data["Total"] = count
	c.Data["Page"] = paginater.New(int(count), opts.PageSize, page, 5)
	c.Data["Logs"] = log

	c.Success(opts.TplName)
}

func Log(c *context.Context) {
	c.Data["Title"] = c.Tr("admin.log.manage_panel")
	c.Data["PageIsAdmin"] = true
	c.Data["PageIsAdminUsers"] = true

	// c.Success(MAIL_LOG_LIST)
	RenderLogSearch(c, &db.LogSearchOptions{
		PageSize: 20,
		OrderBy:  "id ASC",
		TplName:  MAIL_LOG_LIST,
	})
}

func LogDelete(c *context.Context) {
	id := c.ParamsInt64(":id")

	if err := db.LogDeleteById(id); err != nil {
		c.Fail(-1, c.Tr("common.fail"))
	} else {
		c.Flash.Success(c.Tr("common.success"))
	}

	c.Redirect("/admin/log")
}

func LogClear(c *context.Context) {
	db.LogClear()
	c.Redirect("/admin/log")
}
