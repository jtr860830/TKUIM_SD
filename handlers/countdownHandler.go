package handlers

import (
	"net/http"
	"sort"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jtr860830/LifePrint-Server/database"
)

func CdHdlr(c *gin.Context) {
	db := database.GetDB()
	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	user := database.User{}

	if err := db.Set("gorm:auto_preload", true).Where(&database.User{Username: username}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "User not found"})
		return
	}

	lower := time.Now()
	upper := lower.Add(time.Hour * 24 * 7)

	cd := []database.CdItem{}

	for _, v := range user.Schedule {
		if t := v.StartTime; lower.Before(t) && upper.After(t) {
			cd = append(cd, database.CdItem{
				BelongsTo: "Personal",
				Event:     v.Event,
				StartTime: t,
				CD:        int(t.Sub(lower).Hours()) / 24,
			})
		}
	}

	for _, g := range user.Groups {
		db.Model(g).Related(&(*g).Schedule)
		for _, v := range (*g).Schedule {
			if t := v.StartTime; lower.Before(t) && upper.After(t) {
				cd = append(cd, database.CdItem{
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
