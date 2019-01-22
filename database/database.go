package database

import (
	"log"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// DBLC Database location
var DBLC = os.Getenv("DBLC")

// DBMS type of database management system
var DBMS = os.Getenv("DBMS")
var DB *gorm.DB

func init() {
	if DBLC == "" {
		DBLC = "root:password@/sd?charset=utf8&parseTime=True&loc=Local"
	}
	if DBMS == "" {
		DBMS = "mysql"
	}

	var err error
	DB, err = gorm.Open(DBMS, DBLC)
	for err != nil {
		log.Println(err)
		time.Sleep(time.Duration(5) * time.Second)
		DB, err = gorm.Open(DBMS, DBLC)
	}

	if !DB.HasTable(&User{}) {
		DB.Set("gorm:table_options", "charset=utf8").AutoMigrate(&User{}, &Group{}, &UserSchedule{}, &GroupSchedule{}, &Backup{}, &Location{})
		DB.Model(&UserSchedule{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
		DB.Model(&Backup{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
		DB.Model(&GroupSchedule{}).AddForeignKey("group_id", "groups(id)", "RESTRICT", "RESTRICT")
	}
}

func GetDB() *gorm.DB {
	return DB
}
