package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"LifePrint-Server/database"
)

func GetGroupScheduleHdlr(c *gin.Context) {
	db := database.GetDB()
	name := c.Query("name")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Querystring name can't be empty"})
		return
	}

	group := database.Group{}

	if err := db.Set("gorm:auto_preload", true).Where(&database.Group{Name: name}).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Group not found"})
		return
	}

	c.JSON(http.StatusOK, group.Schedule)
}

func AddGroupScheduleHdlr(c *gin.Context) {
	db := database.GetDB()
	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	name := c.PostForm("name")
	event := c.PostForm("event")
	startTime, _ := time.Parse(time.RFC3339, c.PostForm("start"))
	endTime, _ := time.Parse(time.RFC3339, c.PostForm("end"))
	lc := c.PostForm("location")
	n, _ := strconv.ParseFloat(c.PostForm("n"), 64)
	e, _ := strconv.ParseFloat(c.PostForm("e"), 64)
	color := c.PostForm("color")
	t := c.PostForm("type")

	if strings.Trim(event, " ") == "" || strings.Trim(startTime.String(), " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Event and time can't be empty"})
		return
	}

	user := database.User{}
	group := database.Group{}

	if err := db.Where(&database.User{Username: username}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	if err := db.Where(&database.Group{Name: name}).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Group not found"})
		return
	}

	if err := db.Model(&group).Association("Schedule").Append(database.GroupSchedule{
		SponsorID: user.ID,
		Event:     event,
		StartTime: startTime,
		EndTime:   endTime,
		Location: database.Location{
			Name: lc,
			E:    e,
			N:    n,
		},
		Color: color,
		Type:  t,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func AddAllGroupScheduleHdlr(c *gin.Context) {
	db := database.GetDB()
	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	event := c.PostForm("event")
	startTime, _ := time.Parse(time.RFC3339, c.PostForm("start"))
	endTime, _ := time.Parse(time.RFC3339, c.PostForm("end"))
	lc := c.PostForm("location")
	n, _ := strconv.ParseFloat(c.PostForm("n"), 64)
	e, _ := strconv.ParseFloat(c.PostForm("e"), 64)
	color := c.PostForm("color")
	t := c.PostForm("type")

	if strings.Trim(event, " ") == "" || strings.Trim(startTime.String(), " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Event and time can't be empty"})
		return
	}

	user := database.User{}

	if err := db.Where(&database.User{Username: username}).Preload("Groups").First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	for _, v := range user.Groups {
		if err := db.Model(&v).Association("Schedule").Append(database.GroupSchedule{
			SponsorID: user.ID,
			Event:     event,
			StartTime: startTime,
			EndTime:   endTime,
			Location: database.Location{
				Name: lc,
				E:    e,
				N:    n,
			},
			Color: color,
			Type:  t,
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func RmGroupScheduleHdlr(c *gin.Context) {
	db := database.GetDB()
	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	name := c.PostForm("name")
	event := c.PostForm("event")
	startTime, _ := time.Parse(time.RFC3339, c.PostForm("start"))
	endTime, _ := time.Parse(time.RFC3339, c.PostForm("end"))
	//location := c.PostForm("location")
	color := c.PostForm("color")
	t := c.PostForm("type")

	user := database.User{}
	group := database.Group{}
	schedule := database.GroupSchedule{}

	if err := db.Where(&database.User{Username: username}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	if err := db.Where(&database.Group{Name: name}).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Group not found"})
		return
	}

	if err := db.Where(&database.GroupSchedule{
		SponsorID: user.ID,
		Event:     event,
		StartTime: startTime,
		EndTime:   endTime,
		Color:     color,
		Type:      t,
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
