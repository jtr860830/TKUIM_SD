package main // import "github.com/jtr860830/SD-Backend"

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func registerHdlr(c *gin.Context) {
	db, err := gorm.Open("mysql", "root:password@/sd?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer db.Close()

	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")
	birthday, _ := time.Parse(time.RFC3339, c.PostForm("birthday"))

	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Parameters can't be empty"})
		return
	}

	var user = User{
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

func profileHdlr(c *gin.Context) {
	db, err := gorm.Open("mysql", "root:password@/sd?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer db.Close()

	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	user := User{}
	if err := db.Where(&User{Username: username}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func udProfileHdlr(c *gin.Context) {
	db, err := gorm.Open("mysql", "root:password@/sd?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer db.Close()

	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	user := User{}
	if err := db.Where(&User{Username: username}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	email := c.PostForm("email")
	birthday, _ := time.Parse(time.RFC3339, c.PostForm("birthday"))

	if err := db.Model(&user).Updates(User{Email: email, Birthday: birthday}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func chpasswdHdlr(c *gin.Context) {
	db, err := gorm.Open("mysql", "root:password@/sd?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	db.Close()

	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	opasswd := c.PostForm("orpassword")
	cpasswd := c.PostForm("chpassword")

	user := User{}

	if err := db.Where(&User{Username: username}).First(&user).Error; err != nil {
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

func logoutHdlr(c *gin.Context) {

}
