package main // import "github.com/jtr860830/SD-Backend"

import (
	"time"

	"github.com/jinzhu/gorm"
)

// User data model
type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null" json:"-"`
	Email    string
	Birthday time.Time `gorm:"not null"`
	Sticker  string
	Friend   []*User  `gorm:"many2many:friendships;association_jointable_foreignkey:friend_id"`
	Groups   []*Group `gorm:"many2many:user_group;"`
	Schedule []userSchedule
	Backup   []backup
}

// Group data model
type Group struct {
	gorm.Model
	Name     string `gorm:"unique;not null"`
	Color    string
	Sticker  string
	Users    []*User `gorm:"many2many:user_group;"`
	Schedule []groupSchedule
}

type userSchedule struct {
	gorm.Model
	UserID    uint
	Event     string
	StartTime time.Time
	EndTime   time.Time
	Location  string
	Color     string
	Note      string
}

type groupSchedule struct {
	gorm.Model
	GroupID   uint
	SponsorID uint
	Event     string
	StartTime time.Time
	EndTime   time.Time
	Location  string
	Color     string
	Note      string
}

type backup struct {
	gorm.Model
	UserID     uint
	Title      string
	Info       string
	Importance int
}

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type payload struct {
	UserID   uint
	Username string
}
