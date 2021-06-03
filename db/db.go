package db

import (
  "fmt"
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
  "time"
)

var db *gorm.DB
var err error

func Init() {
  dsn := "root:root@tcp(127.0.0.1:3306)/imail?charset=utf8mb4&parseTime=True"
  db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
  // defer db.Close()
  if err != nil {
    fmt.Println("init db err,link error!")
    return
  }

  fmt.Println("init db success!")

  sqlDB, sqlErr := db.DB()
  // SetMaxIdleConns 设置空闲连接池中连接的最大数量
  sqlDB.SetMaxIdleConns(10)
  // SetMaxOpenConns 设置打开数据库连接的最大数量。
  sqlDB.SetMaxOpenConns(100)
  // SetConnMaxLifetime 设置了连接可复用的最大时间。
  sqlDB.SetConnMaxLifetime(time.Hour)

  if sqlErr != nil {
    fmt.Println(sqlErr)
    return
  }

  db.AutoMigrate(&User{})
  db.AutoMigrate(&Mail{})

  //创建默认账户

  var user User
  d := db.First(&user, "name = ?", "admin")
  // fmt.Println(d.Error)
  if d.Error != nil {
    db.Create(&User{
      Name:     "admin",
      Password: "21232f297a57a5a743894a0e4a801fc3",
      Code:     "admin",
    })
  }

  LoginWithCode("admin@xxx.com", "admin")

  MailPush("admin@xxx.com", "midoks@163.com", "tedmm")
}
