package handlers

import (
	"net/http"
	"sort"
	"strconv"
	"strings"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"LifePrint-Server/database"
)

func GetBackupHdlr(c *gin.Context) {
	db := database.GetDB()
	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	user := database.User{}
	bpData := []database.Backup{}

	if err := db.Where(&database.User{Username: username}).Preload("Backup").First(&user).Error; err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "User not found"})
		return
	}

	bpData = user.Backup
	sort.Slice(bpData, func(i, j int) bool {
		return bpData[i].Importance > bpData[j].Importance
	})

	c.JSON(http.StatusOK, bpData)
}

func AddBackupHdlr(c *gin.Context) {
	db := database.GetDB()
	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	title := c.PostForm("title")
	info := c.PostForm("info")
	importance, _ := strconv.Atoi(c.PostForm("importance"))

	if strings.Trim(title, " ") == "" || strings.Trim(info, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Title and content can't be empty"})
		return
	}

	user := database.User{}

	if err := db.Where(&database.User{Username: username}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "User not found"})
		return
	}

	if err := db.Model(&user).Association("Backup").Append(database.Backup{
		Title:      title,
		Info:       info,
		Importance: importance,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func RmBackupHdlr(c *gin.Context) {
	db := database.GetDB()
	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	title := c.PostForm("title")
	info := c.PostForm("info")
	importance, _ := strconv.Atoi(c.PostForm("importance"))

	user := database.User{}
	bp := database.Backup{}

	if err := db.Where(&database.User{Username: username}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	if err := db.Where(&database.Backup{
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
