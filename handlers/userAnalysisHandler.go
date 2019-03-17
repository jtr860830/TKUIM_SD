package handlers

import (
	"net/http"
	"sort"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"LifePrint-Server/database"
)

func GetUserAnalysisHdlr(c *gin.Context) {
	db := database.GetDB()
	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	user := database.User{}

	if err := db.Where(&database.User{Username: username}).Preload("Groups").First(&user).Error; err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "User not found"})
		return
	}

	data := []database.UsrAnlData{}

	for _, v := range user.Groups {
		data = append(data, database.UsrAnlData{Groupname: (*v).Name, Cnt: db.Model(v).Association("Schedule").Count()})
	}

	c.JSON(http.StatusOK, data)
}

func GetUserAnalysis2Hdlr(c *gin.Context) {
	db := database.GetDB()
	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	user := database.User{}

	if err := db.Where(&database.User{Username: username}).Preload("Groups").First(&user).Error; err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "User not found"})
		return
	}

	data := []database.UsrAnl2Data{}

	for k, g := range user.Groups {
		db.Model(g).Related(&(*g).Schedule)
		data = append(data, database.UsrAnl2Data{Groupname: (*g).Name, Cnt: 0})
		for _, v := range (*g).Schedule {
			data[k].Cnt += v.EndTime.Sub(v.StartTime).Seconds()
		}
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].Cnt > data[j].Cnt
	})

	c.JSON(http.StatusOK, data)
}
