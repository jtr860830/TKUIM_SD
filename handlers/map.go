package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/jtr860830/LifePrint-Server/database"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

func AnlMap(c *gin.Context) {
	db := database.GetDB()
	name := c.Query("name")

	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	data := []database.MapData{}

	if name == "" {
		user := database.User{}
		if err := db.Set("gorm:auto_preload", true).Where(&database.User{Username: username}).First(&user).Error; err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"message": "User not found"})
			return
		}

		for _, v := range user.Schedule {
			data = append(data, database.MapData{
				Event: v.Event,
				Type:  v.Type,
				E:     v.Location.E,
				N:     v.Location.N,
			})
		}
	} else {
		group := database.Group{}
		if err := db.Set("gorm:auto_preload", true).Where(&database.Group{Name: name}).First(&group).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Group not found"})
			return
		}

		for _, v := range group.Schedule {
			data = append(data, database.MapData{
				Event: v.Event,
				Type:  v.Type,
				E:     v.Location.E,
				N:     v.Location.N,
			})
		}
	}

	c.JSON(http.StatusOK, data)
}

func AnlMapTimeWeek(c *gin.Context) {
	db := database.GetDB()
	name := c.Query("name")
	size, _ := strconv.Atoi(c.Query("size"))

	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	upper := time.Now()
	lower := upper.AddDate(0, 0, -7*size)

	data := []database.MapData{}

	if name == "" {
		user := database.User{}
		if err := db.Set("gorm:auto_preload", true).Where(&database.User{Username: username}).First(&user).Error; err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"message": "User not found"})
			return
		}

		for _, v := range user.Schedule {
			if lower.Before(v.StartTime) && upper.After(v.StartTime) {
				data = append(data, database.MapData{
					Event: v.Event,
					Type:  v.Type,
					E:     v.Location.E,
					N:     v.Location.N,
				})
			}
		}
	} else {
		group := database.Group{}
		if err := db.Set("gorm:auto_preload", true).Where(&database.Group{Name: name}).First(&group).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Group not found"})
			return
		}

		for _, v := range group.Schedule {
			if lower.Before(v.StartTime) && upper.After(v.StartTime) {
				data = append(data, database.MapData{
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

func AnlMapTimeMonth(c *gin.Context) {
	db := database.GetDB()
	name := c.Query("name")
	size, _ := strconv.Atoi(c.Query("size"))

	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	upper := time.Now()
	lower := upper.AddDate(0, -size, 0)

	data := []database.MapData{}

	if name == "" {
		user := database.User{}
		if err := db.Set("gorm:auto_preload", true).Where(&database.User{Username: username}).First(&user).Error; err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"message": "User not found"})
			return
		}

		for _, v := range user.Schedule {
			if lower.Before(v.StartTime) && upper.After(v.StartTime) {
				data = append(data, database.MapData{
					Event: v.Event,
					Type:  v.Type,
					E:     v.Location.E,
					N:     v.Location.N,
				})
			}
		}
	} else {
		group := database.Group{}
		if err := db.Set("gorm:auto_preload", true).Where(&database.Group{Name: name}).First(&group).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Group not found"})
			return
		}

		for _, v := range group.Schedule {
			if lower.Before(v.StartTime) && upper.After(v.StartTime) {
				data = append(data, database.MapData{
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

func AnlMapTimeYear(c *gin.Context) {
	db := database.GetDB()
	name := c.Query("name")
	size, _ := strconv.Atoi(c.Query("size"))

	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	upper := time.Now()
	lower := upper.AddDate(-size, 0, 0)

	data := []database.MapData{}

	if name == "" {
		user := database.User{}
		if err := db.Set("gorm:auto_preload", true).Where(&database.User{Username: username}).First(&user).Error; err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"message": "User not found"})
			return
		}

		for _, v := range user.Schedule {
			if lower.Before(v.StartTime) && upper.After(v.StartTime) {
				data = append(data, database.MapData{
					Event: v.Event,
					Type:  v.Type,
					E:     v.Location.E,
					N:     v.Location.N,
				})
			}
		}
	} else {
		group := database.Group{}
		if err := db.Set("gorm:auto_preload", true).Where(&database.Group{Name: name}).First(&group).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Group not found"})
			return
		}

		for _, v := range group.Schedule {
			if lower.Before(v.StartTime) && upper.After(v.StartTime) {
				data = append(data, database.MapData{
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
