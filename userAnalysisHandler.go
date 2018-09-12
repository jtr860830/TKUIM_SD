package main // import "github.com/jtr860830/SD-Backend"

import (
	"log"
	"net/http"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func getUserAnalysisHdlr(c *gin.Context) {
	db, err := gorm.Open("mysql", "root:password@/sd?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer db.Close()

	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	user := User{}

	if err := db.Where(&User{Username: username}).Preload("Groups").First(&user).Error; err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "User not found"})
		return
	}

	data := []usrAnlData{}

	for _, v := range user.Groups {
		data = append(data, usrAnlData{Groupname: (*v).Name, Cnt: db.Model(v).Association("Schedule").Count()})
	}

	c.JSON(http.StatusOK, data)
}

func getUserAnalysis2Hdlr(c *gin.Context) {
	db, err := gorm.Open("mysql", "root:password@/sd?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer db.Close()

	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	user := User{}

	if err := db.Where(&User{Username: username}).Preload("Groups").First(&user).Error; err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "User not found"})
		return
	}

	data := []usrAnl2Data{}

	for k, g := range user.Groups {
		db.Model(g).Related(&(*g).Schedule)
		data = append(data, usrAnl2Data{Groupname: (*g).Name, Cnt: 0})
		for _, v := range (*g).Schedule {
			data[k].Cnt += v.EndTime.Sub(v.StartTime).Seconds()
		}
	}

	c.JSON(http.StatusOK, data)
}
