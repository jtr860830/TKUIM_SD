package main // import "github.com/jtr860830/SD-Backend"

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func getGroupAnalysisHdlr(c *gin.Context) {
	db, err := gorm.Open("mysql", "root:password@/sd?charset=utf8&parseTime=True&loc=Local")
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

	data := []anlData{}

	for k, v := range group.Users {
		data = append(data, anlData{Username: (*v).Username, Cnt: 0})
		for _, c := range group.Schedule {
			if (*v).ID == c.SponsorID {
				data[k].Cnt++
			}
		}
	}

	c.JSON(http.StatusOK, data)
}
