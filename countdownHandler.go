package main // import "github.com/jtr860830/SD-Backend"

import (
	"log"
	"net/http"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func cdHdlr(c *gin.Context) {
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

	if err := db.Set("gorm:auto_preload", true).Where(&User{Username: username}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "User not found"})
		return
	}

	//lower := time.Now()
	//upper := lower.Add(time.Hour * 24 * 7)

}
