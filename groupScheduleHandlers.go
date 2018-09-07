package main // import "github.com/jtr860830/SD-Backend"

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func getGroupScheduleHdlr(c *gin.Context) {
	db, err := gorm.Open("mysql", "root:password@/sd?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer db.Close()

	name := c.Query("name")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Querystring name can't be empty"})
		return
	}

	group := Group{}

	if err := db.Where(&Group{Name: name}).Preload("Schedule").First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Group not found"})
		return
	}

	c.JSON(http.StatusOK, group.Schedule)
}

func addGroupScheduleHdlr(c *gin.Context) {
	db, err := gorm.Open("mysql", "root:password@/sd?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer db.Close()

	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	name := c.PostForm("name")
	event := c.PostForm("event")
	startTime, _ := time.Parse(time.RFC3339, c.PostForm("start"))
	endTime, _ := time.Parse(time.RFC3339, c.PostForm("end"))
	location := c.PostForm("location")
	color := c.PostForm("color")
	note := c.PostForm("note")

	if strings.Trim(event, " ") == "" || strings.Trim(startTime.String(), " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Event and time can't be empty"})
		return
	}

	user := User{}
	group := Group{}

	if err := db.Where(&User{Username: username}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	if err := db.Where(&Group{Name: name}).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Group not found"})
		return
	}

	if err := db.Model(&group).Association("Schedule").Append(groupSchedule{
		SponsorID: user.ID,
		Event:     event,
		StartTime: startTime,
		EndTime:   endTime,
		Location:  location,
		Color:     color,
		Note:      note,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func rmGroupScheduleHdlr(c *gin.Context) {
	db, err := gorm.Open("mysql", "root:password@/sd?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer db.Close()

	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	name := c.PostForm("name")
	event := c.PostForm("event")
	startTime, _ := time.Parse(time.RFC3339, c.PostForm("start"))
	endTime, _ := time.Parse(time.RFC3339, c.PostForm("end"))
	location := c.PostForm("location")
	color := c.PostForm("color")
	note := c.PostForm("note")

	user := User{}
	group := Group{}
	schedule := groupSchedule{}

	if err := db.Where(&User{Username: username}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	if err := db.Where(&Group{Name: name}).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Group not found"})
		return
	}

	if err := db.Where(&groupSchedule{
		SponsorID: user.ID,
		Event:     event,
		StartTime: startTime,
		EndTime:   endTime,
		Location:  location,
		Color:     color,
		Note:      note,
	}).First(&schedule).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "You don't have this schedule"})
		return
	}

	if err := db.Delete(&schedule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}
