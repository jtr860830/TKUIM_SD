package main // import "github.com/jtr860830/SD-Backend"

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {

	initDB()

	route := gin.Default()
	store := cookie.NewStore([]byte("secret-string"))
	route.Use(sessions.Sessions("session", store))

	route.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Home Page")
	})
	route.POST("/login", loginHandler)
	route.GET("/logout", logoutHandler)
	route.POST("/register", registerHandler)

	account := route.Group("/user", auth())
	{
		account.GET("/", func(c *gin.Context) {
			c.Redirect(301, "/user/profile")
		})
		account.GET("/profile", profileHandler)
	}

	route.Run(":8080")
}

func initDB() {
	db, err := gorm.Open("mysql", "root:password@/sd?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	if !db.HasTable(&User{}) {
		db.AutoMigrate(&User{}, &Group{}, &userSchedule{}, &groupSchedule{}, &backup{})
		db.Model(&userSchedule{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
		db.Model(&backup{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
		db.Model(&groupSchedule{}).AddForeignKey("group_id", "users(id)", "RESTRICT", "RESTRICT")
	}
}

func auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			c.String(http.StatusNotAcceptable, "You should not pass!")
			log.Println("A strangers attempted to log in!")
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func loginHandler(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")

	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Parameters can't be empty"})
		return
	}

	db, err := gorm.Open("mysql", "root:sd2018@/sd2018DB?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer db.Close()

	user := User{}
	if err := db.Where(&User{Username: username}).Find(&user); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
		return
	}
	if user.Password != password {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
		return
	}

	session.Set("user", username)
	err = session.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate session token"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user"})
	}
}

func logoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
	} else {
		log.Println(user)
		session.Delete("user")
		session.Save()
		c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
	}
}

func registerHandler(c *gin.Context) {
	db, err := gorm.Open("mysql", "root:sd2018@/sd2018DB?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer db.Close()

	username := c.PostForm("username")
	password := c.PostForm("Password")
	email := c.PostForm("email")
	birthday, _ := time.Parse("0000-00-00 00:00:00", c.PostForm("birthday"))

	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}

	var user = User{
		Username: username,
		Password: password,
		Email:    email,
		Birthday: birthday,
	}

	if err = db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func profileHandler(c *gin.Context) {
	db, err := gorm.Open("mysql", "root:sd2018@/sd2018DB?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer db.Close()

	session := sessions.Default(c)
	username := session.Get("user").(string)

	user := User{}
	if err := db.Where(&User{Username: username}).First(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK,
		gin.H{
			"username": user.Username,
			"email":    user.Email,
			"birthday": user.Birthday,
			"sticker":  user.Sticker,
		},
	)
}
