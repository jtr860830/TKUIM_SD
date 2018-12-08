package handlers

import (
	"net/http"

	"github.com/jtr860830/LifePrint-Server/database"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

func GetGroupHdlr(c *gin.Context) {
	db := database.GetDB()
	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	user := database.User{}

	if err := db.Where(&database.User{Username: username}).Preload("Groups").First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not Found"})
		return
	}

	c.JSON(http.StatusOK, user.Groups)
}

func CreateGroupHdlr(c *gin.Context) {
	db := database.GetDB()
	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	name := c.PostForm("name")
	color := c.PostForm("color")
	sticker := c.PostForm("sticker")

	user := database.User{}

	if err := db.Where(&database.User{Username: username}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	if err := db.Model(&user).Association("Groups").Append(database.Group{
		Name:    name,
		Color:   color,
		Sticker: sticker,
	}).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid group name or database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}
