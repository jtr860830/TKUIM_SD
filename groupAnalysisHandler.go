package main // import "github.com/jtr860830/SD-Backend"

import (
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func getGroupAnalysisHdlr(c *gin.Context) {
	db, err := gorm.Open(DBMS, DBLC)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer db.Close()

	name := c.Query("name")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Querystring name can't be empty"})
		return
	}

	group := Group{}

	if err := db.Where(&Group{Name: name}).Preload("Schedule").Preload("Users").First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Group not found"})
		return
	}

	data := []gpAnlData{}

	for k, v := range group.Users {
		data = append(data, gpAnlData{Username: (*v).Username, Cnt: 0})
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

func getGroupAnalysis2Hdlr(c *gin.Context) {
	db, err := gorm.Open(DBMS, DBLC)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer db.Close()

	name := c.Query("name")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Querystring name can't be empty"})
		return
	}

	group := Group{}

	if err := db.Where(&Group{Name: name}).Preload("Schedule").First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Group not found"})
		return
	}

	data := [12]gpAnl2Data{}

	for _, v := range group.Schedule {
		if v.StartTime.Year() == time.Now().Year() {
			data[v.StartTime.Month()-1].Cnt++
		}
	}

	c.JSON(http.StatusOK, data)
}
