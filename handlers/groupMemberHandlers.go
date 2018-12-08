package handlers

import (
	"net/http"

	"github.com/jtr860830/LifePrint-Server/database"

	"github.com/gin-gonic/gin"
)

func GetMemberHdlr(c *gin.Context) {
	db := database.GetDB()
	name := c.Query("name")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Querystring name can't be empty"})
		return
	}

	group := database.Group{}

	if err := db.Where(&database.Group{Name: name}).Preload("Users").First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Group not found"})
		return
	}

	c.JSON(http.StatusOK, group.Users)
}

func AddMemberHdlr(c *gin.Context) {
	db := database.GetDB()
	username := c.PostForm("username")
	name := c.PostForm("name")

	user := database.User{}
	group := database.Group{}

	if err := db.Where(&database.User{Username: username}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	if err := db.Where(&database.Group{Name: name}).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Group not found"})
		return
	}

	if err := db.Model(&user).Association("Groups").Append(group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func RmMemberHdlr(c *gin.Context) {
	db := database.GetDB()
	username := c.PostForm("username")
	name := c.PostForm("name")

	user := database.User{}
	group := database.Group{}

	if err := db.Where(&database.User{Username: username}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	if err := db.Where(&database.Group{Name: name}).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Group not found"})
		return
	}

	if err := db.Model(&user).Association("Groups").Delete(group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}
