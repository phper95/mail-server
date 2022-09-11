package user

import (
	"github.com/phper95/mail-server/internal/app/context"
	"github.com/phper95/mail-server/internal/app/form"
	"github.com/phper95/mail-server/internal/db"
	"github.com/phper95/mail-server/internal/tools"
)

const (
	SETTINGS_PROFILE      = "user/settings/profile"
	SETTINGS_AVATAR       = "user/settings/avatar"
	SETTINGS_AUTHPASSWORD = "user/settings/auth_password"
	SETTINGS_PASSWORD     = "user/settings/password"
	SETTINGS_EMAILS       = "user/settings/email"
	SETTINGS_DELETE       = "user/settings/delete"
	SETTINGS_CLIENT       = "user/settings/client"
)

func Settings(c *context.Context) {
	c.Title("settings.profile")
	c.PageIs("SettingsProfile")
	c.Success(SETTINGS_PROFILE)

}

func SettingsPost(c *context.Context, f form.UpdateProfile) {
	c.Title("settings.profile")
	c.PageIs("SettingsProfile")

	if c.HasError() {
		c.Success(SETTINGS_PROFILE)
		return
	}

	if err := db.UserUpdateNickGetByName(c.User.Name, f.Nick); err != nil {
		msg := err.Error()
		c.RenderWithErr(msg, SETTINGS_PROFILE, &f)
		c.Errorf(err, "update user nick")
		return
	}

	c.Flash.Success(c.Tr("settings.update_profile_success"))
	c.RedirectSubpath("/user/settings")
	c.Success(SETTINGS_PROFILE)
}

func SettingsAuthPassword(c *context.Context) {
	c.Title("settings.auth_password")
	c.PageIs("SettingsAuthPassword")

	token := c.Query("token")
	c.Data["AuthPassword"] = token

	c.Success(SETTINGS_AUTHPASSWORD)
}

func SettingsAuthPasswordPost(c *context.Context, f form.Empty) {
	c.Title("settings.auth_password")
	c.PageIs("SettingsAuthPassword")

	if c.HasError() {
		c.Success(SETTINGS_PASSWORD)
		return
	}

	token := tools.RandString(8)

	err := db.UserUpdateCodeGetById(c.User.Id, token)
	if err != nil {
		c.Flash.Error(err.Error())
		c.Success(SETTINGS_AUTHPASSWORD)
		return
	}

	c.Data["AuthPassword"] = token
	c.Flash.Success(c.Tr("settings.auth_password_success"))
	c.RedirectSubpath("/user/settings/authpassword?token=" + token)
	c.Success(SETTINGS_AUTHPASSWORD)
}

func SettingsPassword(c *context.Context) {
	c.Title("settings.password")
	c.PageIs("SettingsPassword")

	c.Success(SETTINGS_PASSWORD)
}

func SettingsPasswordPost(c *context.Context, f form.ChangePassword) {
	c.Title("settings.password")
	c.PageIs("SettingsPassword")

	if c.HasError() {
		c.Success(SETTINGS_PASSWORD)
		return
	}

	if !c.User.ValidPassword(f.OldPassword) {
		c.Flash.Error(c.Tr("settings.password_incorrect"))
	} else if f.Password != f.Retype {
		c.Flash.Error(c.Tr("form.password_not_match"))
	} else {

		salt := tools.RandString(10)
		c.User.Password = tools.Md5(tools.Md5(f.Password) + salt)
		if err := db.UserUpdater(c.User); err != nil {
			c.Errorf(err, "update user")
			return
		}
		c.Flash.Success(c.Tr("settings.change_password_success"))
	}

	c.RedirectSubpath("/user/settings/password")
}

func SettingsClient(c *context.Context) {
	c.Title("settings.client")
	c.PageIs("SettingsClient")
	c.Success(SETTINGS_CLIENT)
}
