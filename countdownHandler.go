package main // import "github.com/jtr860830/SD-Backend"

import (
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func cdHdlr(c *gin.Context) {
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

	if err := db.Set("gorm:auto_preload", true).Where(&User{Username: username}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "User not found"})
		return
	}

	lower := time.Now()
	upper := lower.Add(time.Hour * 24 * 7)

	cd := []cdItem{}

	for _, v := range user.Schedule {
		if t := v.StartTime; lower.After(t) && upper.Before(t) {
			cd = append(cd, cdItem{
				BelongsTo: "Personal",
				Event:     v.Event,
				StartTime: t,
				CD:        int(t.Sub(lower).Hours()) / 24,
			})
		}
	}

	for _, g := range user.Groups {
		for _, v := range (*g).Schedule {
			if t := v.StartTime; lower.After(t) && upper.Before(t) {
				cd = append(cd, cdItem{
					BelongsTo: (*g).Name,
					Event:     v.Event,
					StartTime: t,
					CD:        int(t.Sub(lower).Hours()) / 24,
				})
			}
		}
	}

	sort.Slice(cd, func(i, j int) bool {
		return cd[i].CD < cd[j].CD
	})

	c.JSON(http.StatusOK, cd)
}
