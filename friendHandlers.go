package main // import "github.com/jtr860830/SD-Backend"

import (
	"net/http"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

func getFriendHdlr(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	user := User{}
	if err := db.Where(&User{Username: username}).Preload("Friend").First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	if user.Friend == nil {
		c.JSON(http.StatusOK, gin.H{"message": "You don't have any friends"})
		return
	}

	c.JSON(http.StatusOK, user.Friend)
}

func addFriendHdlr(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)
	friendname := c.PostForm("username")

	user := User{}
	friend := User{}

	if err := db.Where(&User{Username: username}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	if err := db.Where(&User{Username: friendname}).First(&friend).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Your friend is not exist"})
		return
	}

	if err := db.Model(&user).Association("Friend").Append(friend).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func rmFriendHdlr(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)
	friendname := c.PostForm("username")

	user := User{}
	friend := User{}

	if err := db.Where(&User{Username: username}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	if err := db.Where(&User{Username: friendname}).First(&friend).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Your friend is not exist"})
		return
	}

	if err := db.Model(&user).Association("Friend").Delete(friend).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}
