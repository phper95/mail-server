package db

import (
	"fmt"
	"strings"
	"time"

	"github.com/midoks/imail/internal/tools"
)

type User struct {
	Id       int64  `gorm:"primaryKey"`
	Name     string `gorm:"unique;size:50;comment:用户名"`
	Password string `gorm:"size:32;comment:用户密码"`
	Code     string `gorm:"size:50;comment:编码"`
	Token    string `gorm:"unique;size:50;comment:Token"`
	Status   int    `gorm:"comment:状态"`
	Salt     string `gorm:"type:varchar(10)"`

	IsActive bool
	IsAdmin  bool

	Created     time.Time `xorm:"-" gorm:"-" json:"-"`
	CreatedUnix int64     `gorm:"autoCreateTime;comment:创建时间"`
	Updated     time.Time `xorm:"-" gorm:"-" json:"-"`
	UpdatedUnix int64     `gorm:"autoCreateTime;comment:更新时间"`
}

func (User) TableName() string {
	return "im_users"
}

// CreateUser creates record of a new user.
// Deprecated: Use Users.Create instead.
func CreateUser(u *User) (err error) {
	data := db.First(u, "name = ?", u.Name)

	u.Salt = tools.RandString(10)
	u.Password = tools.Md5(tools.Md5(u.Password) + u.Salt)
	if data.Error != nil {
		db.Create(u)
	}
	return nil
}

func LoginWithCode(name string, code string) (bool, int64) {
	list := strings.SplitN(name, "@", 2)

	var user User
	err := db.First(&user, "name = ?", list[0]).Error

	if err != nil {
		return false, 0
	}

	if user.Code == code {
		return true, user.Id
	}

	return false, 0
}

func LoginByUserPassword(name string, password string) (bool, int64) {

	var u User
	err := db.First(&u, "name = ?", name).Error

	if err != nil {
		return false, 0
	}

	inputPwd := tools.Md5(tools.Md5(password) + u.Salt)

	fmt.Println(password, inputPwd, u)
	if inputPwd == u.Password {
		return true, u.Id
	}
	return false, 0
}

func UserCheckIsExist(name string) bool {
	var user User
	err := db.First(&user, "name = ?", name).Error
	if err == nil {
		return true
	}
	return false
}

func UserUpdateTokenGetByName(name string, token string) bool {
	db.Model(&User{}).Where("name = ?", name).Update("token", token)
	return true
}

func UserUpdateCodeGetByName(name string, code string) bool {
	db.Model(&User{}).Where("name = ?", name).Update("code", code)
	return true
}

func UserGetByName(name string) (User, error) {
	list := strings.SplitN(name, "@", 2)
	var user User
	err := db.First(&user, "name = ?", list[0]).Error
	if err != nil {
		return user, err
	}
	return user, nil
}
