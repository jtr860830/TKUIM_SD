package main // import "github.com/jtr860830/SD-Backend"

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

func getScheduleHdlr(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	user := User{}

	if err := db.Set("gorm:auto_preload", true).Where(&User{Username: username}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user.Schedule)
}

func addScheduleHdlr(c *gin.Context) {
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

	user := User{}

	if err := db.Where(&User{Username: username}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	if err := db.Model(&user).Association("Schedule").Append(userSchedule{
		Event:     event,
		StartTime: startTime,
		EndTime:   endTime,
		Location: location{
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

func udScheduleHdlr(c *gin.Context) {

}

func rmScheduleHdlr(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	event := c.PostForm("event")
	startTime, _ := time.Parse(time.RFC3339, c.PostForm("start"))
	endTime, _ := time.Parse(time.RFC3339, c.PostForm("end"))
	//location := c.PostForm("location")
	color := c.PostForm("color")
	t := c.PostForm("type")

	user := User{}
	schedule := userSchedule{}

	if err := db.Where(&User{Username: username}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	if err := db.Where(&userSchedule{
		UserID:    user.ID,
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
