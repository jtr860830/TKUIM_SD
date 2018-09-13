package main // import "github.com/jtr860830/SD-Backend"

import (
	"log"
	"net/http"
	"strconv"
	"time"

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
		if err := db.Set("gorm:auto_preload", true).Where(&User{Username: username}).First(&user).Error; err != nil {
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
		if err := db.Set("gorm:auto_preload", true).Where(&Group{Name: name}).First(&group).Error; err != nil {
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

func anlMapTimeWeek(c *gin.Context) {
	db, err := gorm.Open("mysql", "root:password@/sd?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer db.Close()

	name := c.Query("name")
	size, _ := strconv.Atoi(c.Query("size"))

	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	upper := time.Now()
	lower := upper.AddDate(0, 0, -7*size)

	data := []mapData{}

	if name == "" {
		user := User{}
		if err := db.Set("gorm:auto_preload", true).Where(&User{Username: username}).First(&user).Error; err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"message": "User not found"})
			return
		}

		for _, v := range user.Schedule {
			if lower.Before(v.StartTime) && upper.After(v.StartTime) {
				data = append(data, mapData{
					Event: v.Event,
					Type:  v.Type,
					E:     v.Location.E,
					N:     v.Location.N,
				})
			}
		}
	} else {
		group := Group{}
		if err := db.Set("gorm:auto_preload", true).Where(&Group{Name: name}).First(&group).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Group not found"})
			return
		}

		for _, v := range group.Schedule {
			if lower.Before(v.StartTime) && upper.After(v.StartTime) {
				data = append(data, mapData{
					Event: v.Event,
					Type:  v.Type,
					E:     v.Location.E,
					N:     v.Location.N,
				})
			}
		}
	}

	c.JSON(http.StatusOK, data)
}

func anlMapTimeMonth(c *gin.Context) {
	db, err := gorm.Open("mysql", "root:password@/sd?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer db.Close()

	name := c.Query("name")
	size, _ := strconv.Atoi(c.Query("size"))

	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	upper := time.Now()
	lower := upper.AddDate(0, size, 0)

	data := []mapData{}

	if name == "" {
		user := User{}
		if err := db.Set("gorm:auto_preload", true).Where(&User{Username: username}).First(&user).Error; err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"message": "User not found"})
			return
		}

		for _, v := range user.Schedule {
			if lower.Before(v.StartTime) && upper.After(v.StartTime) {
				data = append(data, mapData{
					Event: v.Event,
					Type:  v.Type,
					E:     v.Location.E,
					N:     v.Location.N,
				})
			}
		}
	} else {
		group := Group{}
		if err := db.Set("gorm:auto_preload", true).Where(&Group{Name: name}).First(&group).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Group not found"})
			return
		}

		for _, v := range group.Schedule {
			if lower.Before(v.StartTime) && upper.After(v.StartTime) {
				data = append(data, mapData{
					Event: v.Event,
					Type:  v.Type,
					E:     v.Location.E,
					N:     v.Location.N,
				})
			}
		}
	}

	c.JSON(http.StatusOK, data)
}

func anlMapTimeYear(c *gin.Context) {
	db, err := gorm.Open("mysql", "root:password@/sd?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer db.Close()

	name := c.Query("name")
	size, _ := strconv.Atoi(c.Query("size"))

	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	upper := time.Now()
	lower := upper.AddDate(size, 0, 0)

	data := []mapData{}

	if name == "" {
		user := User{}
		if err := db.Set("gorm:auto_preload", true).Where(&User{Username: username}).First(&user).Error; err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"message": "User not found"})
			return
		}

		for _, v := range user.Schedule {
			if lower.Before(v.StartTime) && upper.After(v.StartTime) {
				data = append(data, mapData{
					Event: v.Event,
					Type:  v.Type,
					E:     v.Location.E,
					N:     v.Location.N,
				})
			}
		}
	} else {
		group := Group{}
		if err := db.Set("gorm:auto_preload", true).Where(&Group{Name: name}).First(&group).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Group not found"})
			return
		}

		for _, v := range group.Schedule {
			if lower.Before(v.StartTime) && upper.After(v.StartTime) {
				data = append(data, mapData{
					Event: v.Event,
					Type:  v.Type,
					E:     v.Location.E,
					N:     v.Location.N,
				})
			}
		}
	}

	c.JSON(http.StatusOK, data)
}
