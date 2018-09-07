package main // import "github.com/jtr860830/SD-Backend"

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func getMemberHdlr(c *gin.Context) {
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

	if err := db.Where(&Group{Name: name}).Preload("Users").First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Group not found"})
		return
	}

	c.JSON(http.StatusOK, group.Users)
}

func addMemberHdlr(c *gin.Context) {
	db, err := gorm.Open("mysql", "root:password@/sd?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer db.Close()

	username := c.PostForm("username")
	name := c.PostForm("name")

	user := User{}
	group := Group{}

	if err := db.Where(&User{Username: username}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	if err := db.Where(&Group{Name: name}).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Group not found"})
		return
	}

	if err := db.Model(&user).Association("Groups").Append(group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func rmMemberHdlr(c *gin.Context) {
	db, err := gorm.Open("mysql", "root:password@/sd?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer db.Close()

	username := c.PostForm("username")
	name := c.PostForm("name")

	user := User{}
	group := Group{}

	if err := db.Where(&User{Username: username}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	if err := db.Where(&Group{Name: name}).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Group not found"})
		return
	}

	if err := db.Model(&user).Association("Groups").Delete(group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}
