package database

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
	Schedule []UserSchedule
	Backup   []Backup
}

// Group data model
type Group struct {
	gorm.Model
	Name     string `gorm:"unique;not null"`
	Color    string
	Sticker  string
	Users    []*User `gorm:"many2many:user_group;"`
	Schedule []GroupSchedule
}

type UserSchedule struct {
	gorm.Model
	UserID     uint
	Event      string
	StartTime  time.Time
	EndTime    time.Time
	Location   Location
	LocationID uint
	Color      string
	Type       string
}

type GroupSchedule struct {
	gorm.Model
	GroupID    uint
	SponsorID  uint
	Event      string
	StartTime  time.Time
	EndTime    time.Time
	Location   Location
	LocationID uint
	Color      string
	Type       string
}

type Backup struct {
	gorm.Model
	UserID     uint
	Title      string
	Info       string
	Importance int
}

type Location struct {
	gorm.Model
	Name string
	E    float64
	N    float64
}

// Data structs
type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type Payload struct {
	UserID   uint
	Username string
}

type CdItem struct {
	BelongsTo string
	Event     string
	StartTime time.Time
	CD        int
}

type GpAnlData struct {
	Username string
	Cnt      int
}

type GpAnl2Data struct {
	Cnt int
}

type UsrAnlData struct {
	Groupname string
	Cnt       int
}

type UsrAnl2Data struct {
	Groupname string
	Cnt       float64
}

type MapData struct {
	Event string
	Type  string
	E     float64
	N     float64
}
