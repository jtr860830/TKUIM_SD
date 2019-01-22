package handlers

import (
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jtr860830/LifePrint-Server/database"
)

func GetGroupAnalysisHdlr(c *gin.Context) {
	db := database.GetDB()
	name := c.Query("name")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Querystring name can't be empty"})
		return
	}

	group := database.Group{}

	if err := db.Where(&database.Group{Name: name}).Preload("Schedule").Preload("Users").First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Group not found"})
		return
	}

	data := []database.GpAnlData{}

	for k, v := range group.Users {
		data = append(data, database.GpAnlData{Username: (*v).Username, Cnt: 0})
		for _, c := range group.Schedule {
			if (*v).ID == c.SponsorID {
				data[k].Cnt++
			}
		}
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].Cnt > data[j].Cnt
	})

	c.JSON(http.StatusOK, data)
}

func GetGroupAnalysis2Hdlr(c *gin.Context) {
	db := database.GetDB()
	name := c.Query("name")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Querystring name can't be empty"})
		return
	}

	group := database.Group{}

	if err := db.Where(&database.Group{Name: name}).Preload("Schedule").First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Group not found"})
		return
	}

	data := [12]database.GpAnl2Data{}

	for _, v := range group.Schedule {
		if v.StartTime.Year() == time.Now().Year() {
			data[v.StartTime.Month()-1].Cnt++
		}
	}

	c.JSON(http.StatusOK, data)
}
