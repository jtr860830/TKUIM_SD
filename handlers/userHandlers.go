package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/jtr860830/LifePrint-Server/database"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

func RegisterHdlr(c *gin.Context) {
	db := database.GetDB()
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")
	birthday, _ := time.Parse(time.RFC3339, c.PostForm("birthday"))

	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Parameters can't be empty"})
		return
	}

	var user = database.User{
		Username: username,
		Password: password,
		Email:    email,
		Birthday: birthday,
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"massage": "Can't use this username or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func ProfileHdlr(c *gin.Context) {
	db := database.GetDB()
	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	user := database.User{}
	if err := db.Where(&database.User{Username: username}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UdProfileHdlr(c *gin.Context) {
	db := database.GetDB()
	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	user := database.User{}
	if err := db.Where(&database.User{Username: username}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	email := c.PostForm("email")
	birthday, _ := time.Parse(time.RFC3339, c.PostForm("birthday"))

	if err := db.Model(&user).Updates(database.User{Email: email, Birthday: birthday}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func ChpasswdHdlr(c *gin.Context) {
	db := database.GetDB()
	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	opasswd := c.PostForm("orpassword")
	cpasswd := c.PostForm("chpassword")

	user := database.User{}

	if err := db.Where(&database.User{Username: username}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	if opasswd != user.Password {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Original password not match"})
		return
	}

	if err := db.Model(&user).Update("password", cpasswd).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func LogoutHdlr(c *gin.Context) {

}
