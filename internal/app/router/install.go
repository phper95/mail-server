package router

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/ini.v1"
	"gopkg.in/macaron.v1"

	"github.com/phper95/mail-server/internal/app/context"
	"github.com/phper95/mail-server/internal/app/form"
	"github.com/phper95/mail-server/internal/conf"
	"github.com/phper95/mail-server/internal/db"
	"github.com/phper95/mail-server/internal/imap"
	"github.com/phper95/mail-server/internal/log"
	"github.com/phper95/mail-server/internal/pop3"
	"github.com/phper95/mail-server/internal/smtpd"
	"github.com/phper95/mail-server/internal/task"
	"github.com/phper95/mail-server/internal/tools"
	"github.com/phper95/mail-server/internal/tools/debug"
)

const (
	INSTALL = "install"
)

func startService(name string) {

	if strings.EqualFold(name, "smtpd") && conf.Smtp.Enable {
		go smtpd.Start(conf.Smtp.Port)
	} else if strings.EqualFold(name, "pop3") && conf.Pop3.Enable {
		go pop3.Start(conf.Pop3.Port)
	} else if strings.EqualFold(name, "imap") && conf.Imap.Enable {
		go imap.Start(conf.Imap.Port)
	}

	log.Infof("listen %s success!", name)

	if strings.EqualFold(name, "smtpd") && conf.Smtp.SslEnable {
		go smtpd.StartSSL(conf.Smtp.SslPort)
	} else if strings.EqualFold(name, "pop3") && conf.Pop3.SslEnable {
		go pop3.StartSSL(conf.Pop3.SslPort)
	} else if strings.EqualFold(name, "imap") && conf.Imap.SslEnable {
		go imap.StartSSL(conf.Imap.SslPort)
	}

	log.Infof("listen %s ssl success!", name)

}

func checkRunMode() {
	if conf.IsProdMode() {
		macaron.Env = macaron.PROD
		macaron.ColorLog = false
	} else {
		macaron.Env = macaron.DEV
	}
	log.Infof("Run mode: %s", strings.Title(macaron.Env))
}

func GlobalInit(customConf string) error {

	err := conf.Init(customConf)
	if err != nil {
		return errors.Wrap(err, "init configuration")
	}

	log.Init()

	// format := conf.Log.Format
	// if strings.EqualFold(format, "json") {
	// 	logger.SetFormatter(&logrus.JSONFormatter{})
	// } else if strings.EqualFold(format, "text") {
	// 	logger.SetFormatter(&logrus.TextFormatter{})
	// } else {
	// 	logger.SetFormatter(&logrus.TextFormatter{})
	// }

	if !strings.EqualFold(conf.App.RunMode, "prod") {
		go debug.Pprof()
	}

	if tools.IsFile("./install.lock") {

		db.Init()
		task.Init()

		startService("smtpd")
		startService("pop3")
		startService("imap")

	}

	checkRunMode()
	if !tools.IsFile("./install.lock") {
		return nil
	}

	return nil
}

func InstallInit(c *context.Context) {
	if tools.IsFile("./install.lock") {
		c.NotFound()
		return
	}

	c.Title("install.install")
	c.PageIs("Install")

	c.Data["DbOptions"] = []string{"MySQL", "PostgreSQL", "SQLite3"}
}

func Install(c *context.Context) {
	f := form.Install{}

	// Database settings
	f.DbHost = conf.Database.Host
	f.DbUser = conf.Database.User
	f.DbName = conf.Database.Name
	f.DbPath = conf.Database.Path

	c.Data["CurDbOption"] = "PostgreSQL"
	switch conf.Database.Type {
	case "mysql":
		c.Data["CurDbOption"] = "MySQL"
	case "sqlite3":
		c.Data["CurDbOption"] = "SQLite3"
	}

	// Application general settings
	f.AppName = conf.App.BrandName

	// Note(unknwon): it's hard for Windows users change a running user,
	// 	so just use current one if config says default.
	if conf.IsWindowsRuntime() && conf.App.RunUser == "git" {
		f.RunUser = tools.CurrentUsername()
	} else {
		f.RunUser = conf.App.RunUser
	}

	f.Domain = conf.Web.Domain
	f.HttpPort = conf.Web.HttpPort
	f.LogRootPath = fmt.Sprintf("%s/%s", conf.WorkDir(), conf.Log.RootPath)

	// Server and other services settings
	form.Assign(f, c.Data)
	c.Success(INSTALL)
}

func InstallPost(c *context.Context, f form.Install) {
	c.Data["CurDbOption"] = f.DbType

	if c.HasError() {
		if c.HasValue("Err_AdminName") ||
			c.HasValue("Err_AdminPasswd") ||
			c.HasValue("Err_AdminEmail") {
			c.FormErr("Admin")
		}

		c.Success(INSTALL)
		return
	}

	if _, err := exec.LookPath("git"); err != nil {
		c.RenderWithErr(c.Tr("install.test_git_failed", err), INSTALL, &f)
		return
	}

	// Pass basic check, now test configuration.
	// Test database setting.
	dbTypes := map[string]string{
		"PostgreSQL": "postgres",
		"MySQL":      "mysql",
		"SQLite3":    "sqlite3",
	}
	conf.Database.Type = dbTypes[f.DbType]
	conf.Database.Host = f.DbHost
	conf.Database.User = f.DbUser
	conf.Database.Password = f.DbPasswd
	conf.Database.Name = f.DbName
	conf.Database.SslMode = f.SslMode
	conf.Database.Path = f.DbPath

	if conf.Database.Type == "sqlite3" && len(conf.Database.Path) == 0 {
		c.FormErr("DbPath")
		c.RenderWithErr(c.Tr("install.err_empty_db_path"), INSTALL, &f)
		return
	}

	// Set test engine.
	// if err := db.NewTestEngine(); err != nil {
	// 	if strings.Contains(err.Error(), `Unknown database type: sqlite3`) {
	// 		c.FormErr("DbType")
	// 		c.RenderWithErr(c.Tr("install.sqlite3_not_available", "https://github.com/phper95/mail-server/wiki"), INSTALL, &f)
	// 	} else {
	// 		c.FormErr("DbSetting")
	// 		c.RenderWithErr(c.Tr("install.invalid_db_setting", err), INSTALL, &f)
	// 	}
	// 	return
	// }

	// Test log root path.
	f.LogRootPath = strings.Replace(f.LogRootPath, "\\", "/", -1)
	if err := os.MkdirAll(f.LogRootPath, os.ModePerm); err != nil {
		c.FormErr("LogRootPath")
		c.RenderWithErr(c.Tr("install.invalid_log_root_path", err), INSTALL, &f)
		return
	}

	currentUser, match := conf.CheckRunUser(f.RunUser)
	if !match {
		c.FormErr("RunUser")
		c.RenderWithErr(c.Tr("install.run_user_not_match", f.RunUser, currentUser), INSTALL, &f)
		return
	}

	if len(f.AdminName) == 0 {
		c.FormErr("Admin", "AdminName")
		c.RenderWithErr(c.Tr("install.admin_setting_desc"), INSTALL, f)
		return
	}

	// Check admin password.
	if len(f.AdminName) > 0 && len(f.AdminPassword) == 0 {
		c.FormErr("Admin", "AdminPassword")
		c.RenderWithErr(c.Tr("install.err_empty_admin_password"), INSTALL, f)
		return
	}
	if f.AdminPassword != f.AdminConfirmPassword {
		c.FormErr("Admin", "AdminPassword")
		c.RenderWithErr(c.Tr("form.password_not_match"), INSTALL, f)
		return
	}

	// Save settings.
	cfg := ini.Empty()
	if tools.IsFile(conf.CustomConf) {
		// Keeps custom settings if there is already something.
		if err := cfg.Append(conf.CustomConf); err != nil {
			log.Errorf("Failed to load custom conf %s: %s", conf.CustomConf, err)
		}
	}

	cfg.Section("").Key("app_name").SetValue("mail-server")
	cfg.Section("").Key("brand_name").SetValue(f.AppName)
	cfg.Section("").Key("run_user").SetValue(f.RunUser)
	cfg.Section("").Key("run_mode").SetValue("prod")

	cfg.Section("database").Key("type").SetValue(conf.Database.Type)
	cfg.Section("database").Key("host").SetValue(conf.Database.Host)
	cfg.Section("database").Key("name").SetValue(conf.Database.Name)
	cfg.Section("database").Key("user").SetValue(conf.Database.User)
	cfg.Section("database").Key("password").SetValue(conf.Database.Password)
	cfg.Section("database").Key("ssl_mode").SetValue(conf.Database.SslMode)
	cfg.Section("database").Key("path").SetValue(conf.Database.Path)

	cfg.Section("web").Key("domain").SetValue(f.Domain)
	cfg.Section("web").Key("http_port").SetValue(fmt.Sprintf("%d", f.HttpPort))
	cfg.Section("web").Key("mail_save_mode").SetValue(f.MailSaveMode)

	cfg.Section("session").Key("provider").SetValue("file")
	cfg.Section("session").Key("cookie_secure").SetValue("false")

	cfg.Section("auth").Key("disable_registration").SetValue("true")

	mode := "file"
	if f.EnableConsoleMode {
		mode = "console, file"
	}
	cfg.Section("log").Key("format").SetValue("text")
	cfg.Section("log").Key("mode").SetValue(mode)
	cfg.Section("log").Key("level").SetValue("Info")
	cfg.Section("log").Key("root_path").SetValue(f.LogRootPath)

	cfg.Section("security").Key("install_lock").SetValue("true")
	secretKey := tools.RandString(15)
	cfg.Section("security").Key("secret_key").SetValue(secretKey)

	os.MkdirAll(filepath.Dir(conf.CustomConf), os.ModePerm)
	if err := cfg.SaveTo(conf.CustomConf); err != nil {
		c.RenderWithErr(c.Tr("install.save_config_failed", err), INSTALL, &f)
		return
	}

	// NOTE: We reuse the current value because this handler does not have access to CLI flags.
	if err := GlobalInit(conf.CustomConf); err != nil {
		c.RenderWithErr(c.Tr("install.init_failed", err), INSTALL, &f)
		return
	}

	// Create admin account
	if len(f.AdminName) < 1 {
		c.FormErr("AdminName", "AdminEmail")
		return
	}

	u := &db.User{
		Name:     f.AdminName,
		Password: f.AdminPassword,
		IsAdmin:  true,
		IsActive: true,
	}

	if err := db.CreateUser(u); err != nil {
		c.FormErr("AdminName", "AdminEmail")
		c.RenderWithErr(c.Tr("install.invalid_admin_setting", err), INSTALL, &f)
		return
	}

	//install lock
	err := tools.WriteFile("./install.lock", "")
	if err != nil {
		log.Errorf("create install lock file error %v", err)
		c.RenderWithErr(c.Tr("install.init_failed", err), INSTALL, &f)
		return
	}

	log.Info("first-time run install finished!")
	c.Flash.Success(c.Tr("install.install_success"))
	// Auto-login for admin
	_ = c.Session.Set("uid", u.Id)
	_ = c.Session.Set("uname", u.Name)

	c.Redirect("/user/login")
}
