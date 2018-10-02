package main // import "github.com/jtr860830/SD-Backend"

import (
	"net/http"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

func getGroupHdlr(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	user := User{}

	if err := db.Where(&User{Username: username}).Preload("Groups").First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not Found"})
		return
	}

	c.JSON(http.StatusOK, user.Groups)
}

func createGroupHdlr(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	name := c.PostForm("name")
	color := c.PostForm("color")
	sticker := c.PostForm("sticker")

	user := User{}

	if err := db.Where(&User{Username: username}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	if err := db.Model(&user).Association("Groups").Append(Group{
		Name:    name,
		Color:   color,
		Sticker: sticker,
	}).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid group name or database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}
