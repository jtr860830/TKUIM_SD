package main // import "github.com/jtr860830/SD-Backend"

import (
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func getBackupHdlr(c *gin.Context) {
	db, err := gorm.Open(DBMS, DBLC)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer db.Close()

	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	user := User{}
	bpData := []backup{}

	if err := db.Where(&User{Username: username}).Preload("Backup").First(&user).Error; err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "User not found"})
		return
	}

	bpData = user.Backup
	sort.Slice(bpData, func(i, j int) bool {
		return bpData[i].Importance > bpData[j].Importance
	})

	c.JSON(http.StatusOK, bpData)
}

func addBackupHdlr(c *gin.Context) {
	db, err := gorm.Open(DBMS, DBLC)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer db.Close()

	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	title := c.PostForm("title")
	info := c.PostForm("info")
	importance, _ := strconv.Atoi(c.PostForm("importance"))

	if strings.Trim(title, " ") == "" || strings.Trim(info, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Title and content can't be empty"})
		return
	}

	user := User{}

	if err := db.Where(&User{Username: username}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "User not found"})
		return
	}

	if err := db.Model(&user).Association("Backup").Append(backup{
		Title:      title,
		Info:       info,
		Importance: importance,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func rmBackupHdlr(c *gin.Context) {
	db, err := gorm.Open(DBMS, DBLC)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer db.Close()

	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	title := c.PostForm("title")
	info := c.PostForm("info")
	importance, _ := strconv.Atoi(c.PostForm("importance"))

	user := User{}
	bp := backup{}

	if err := db.Where(&User{Username: username}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	if err := db.Where(&backup{
		UserID:     user.ID,
		Title:      title,
		Info:       info,
		Importance: importance,
	}).First(&bp).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "You don't have this backup"})
		return
	}

	if err := db.Delete(&bp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}
