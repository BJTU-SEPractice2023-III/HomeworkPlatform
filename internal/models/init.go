package models

import (
	"homework_platform/internal/bootstrap"
	"homework_platform/internal/utils"
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm/logger"

	// "gocloud.dev/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// db, err := gorm.Open(sqlite.Open("ach.db"), &gorm.Config{})
	// db, err := gorm.Open(mysql.Open(bootstrap.Config.SQLDSN), &gorm.Config{})
	var db *gorm.DB
	var err error
	if bootstrap.Sqlite {
		db, err = gorm.Open(sqlite.Open("homework-platform.db"), &gorm.Config{})
		sqlDB, err := db.DB()
		sqlDB.SetMaxOpenConns(1)
		if err != nil {
			log.Panicln("无法设置连接池")
		}
	} else if bootstrap.Mysql {
		db, err = gorm.Open(mysql.Open(bootstrap.Config.SQLDSN), &gorm.Config{})
	} else if bootstrap.Test {
		db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		log.Println("正在使用memory sqlite数据库")
	} else {
		db, err = gorm.Open(postgres.Open(bootstrap.Config.SQLDSN), &gorm.Config{})
	}

	if err != nil {
		log.Panicf("无法连接数据库，%s", err)
	}

	DB = db
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&Course{})
	DB.AutoMigrate(&Homework{})
	DB.AutoMigrate(&HomeworkSubmission{})
	DB.AutoMigrate(&Comment{})
	DB.AutoMigrate(&Complaint{})
	// 创建初始管理员账户
	addDefaultUser()
}

func addDefaultUser() {
	_, err := GetUserByID(1)
	password := utils.RandStringRunes(8)

	if err == gorm.ErrRecordNotFound {
		defaultUser := &User{}

		defaultUser.Username = "Admin"
		defaultUser.Password = utils.EncodePassword(password, utils.RandStringRunes(16))
		defaultUser.IsAdmin = true

		if err := DB.Create(&defaultUser).Error; err != nil {
			log.Panicf("创建初始管理员账户失败: %s\n", err)
		}

		log.Println("初始管理员账户创建完成")
		log.Printf("用户名: %s\n", "Admin")
		log.Printf("密码: %s\n", password)
	}
}
