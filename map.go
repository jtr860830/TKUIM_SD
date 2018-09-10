package main // import "github.com/jtr860830/SD-Backend"

import (
	"log"
	"net/http"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func anlMap(c *gin.Context) {
	db, err := gorm.Open("mysql", "root:password@/sd?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer db.Close()

	name := c.Query("name")

	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	data := []mapData{}

	if name == "" {
		user := User{}
		if err := db.Where(&User{Username: username}).Preload("Schedule").First(&user).Error; err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"message": "User not found"})
			return
		}

		for _, v := range user.Schedule {
			data = append(data, mapData{
				Event: v.Event,
				Type:  v.Type,
				E:     v.Location.E,
				N:     v.Location.N,
			})
		}
	} else {
		group := Group{}
		if err := db.Where(&Group{Name: name}).Preload("Schedule").First(&group).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Group not found"})
			return
		}

		for _, v := range group.Schedule {
			data = append(data, mapData{
				Event: v.Event,
				Type:  v.Type,
				E:     v.Location.E,
				N:     v.Location.N,
			})
		}
	}

	c.JSON(http.StatusOK, data)
}
