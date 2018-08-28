package main // import "github.com/jtr860830/SD-Backend"

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type payload struct {
	UserID   uint
	Username string
}

func main() {

	initDB()

	route := gin.Default()

	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:      "test zone",
		Key:        []byte("secret-key"),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*payload); ok {
				return jwt.MapClaims{
					"id":       v.UserID,
					"username": v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.Bind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := loginVals.Username
			password := loginVals.Password

			db, err := gorm.Open("mysql", "root:password@/sd?charset=utf8&parseTime=True&loc=Local")
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
				return nil, jwt.ErrFailedAuthentication
			}
			defer db.Close()

			user := User{}
			if err := db.Where(&User{Username: username}).Find(&user).Error; err != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			if user.Password != password {
				return nil, jwt.ErrFailedAuthentication
			}

			return &payload{
				UserID:   user.ID,
				Username: user.Username,
			}, nil
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},

		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		TimeFunc:    time.Now,
	}

	route.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	route.POST("/login", authMiddleware.LoginHandler)
	route.GET("/logout", logoutHdlr)
	route.POST("/register", registerHdlr)

	account := route.Group("/user", authMiddleware.MiddlewareFunc())
	{
		account.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/user/profile")
		})

		account.GET("/profile", profileHdlr)

		account.POST("/chpasswd", chpasswdHdlr)

		account.GET("/friends", getFriendHdlr)
		account.POST("/friends", addFriendHdlr)
		account.DELETE("/friends", rmFriendHdlr)

		account.GET("/schedules", getScheduleHdlr)
		account.POST("/schedules", addScheduleHdlr)
		account.PATCH("/schedules", udScheduleHdlr)
		account.DELETE("/schedules", rmScheduleHdlr)

		account.GET("/backups", getBackupHdlr)
		account.POST("/backups", addBackupHdlr)
		account.DELETE("/backups", rmBackupHdlr)

		group := account.Group("/group")
		{
			group.GET("/", getGroupHdlr)
			group.POST("/", createGroupHdlr)

			group.GET("/member", getMemberHdlr)
			group.POST("/member", addMemberHdlr)
			group.DELETE("/member", rmMemberHdlr)

			group.GET("/schedules", getGroupScheduleHdlr)
			group.POST("/schedules", addGroupScheduleHdlr)
			group.DELETE("/schedules", rmGroupScheduleHdlr)
		}
	}

	route.Run(":8080")
}

func initDB() {
	db, err := gorm.Open("mysql", "root:password@/sd?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
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

func logoutHdlr(c *gin.Context) {

}

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

func getFriendHdlr(c *gin.Context) {
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
	db, err := gorm.Open("mysql", "root:password@/sd?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer db.Close()

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
	db, err := gorm.Open("mysql", "root:password@/sd?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer db.Close()

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

func getScheduleHdlr(c *gin.Context) {
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

	if err := db.Where(&User{Username: username}).Preload("Schedule").First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user.Schedule)
}

func addScheduleHdlr(c *gin.Context) {
	db, err := gorm.Open("mysql", "root:password@/sd?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer db.Close()

	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	event := c.PostForm("event")
	startTime, _ := time.Parse(time.RFC3339, c.PostForm("start"))
	endTime, _ := time.Parse(time.RFC3339, c.PostForm("end"))
	location := c.PostForm("location")
	color := c.PostForm("color")
	note := c.PostForm("note")

	if strings.Trim(event, " ") == "" || strings.Trim(startTime.String(), " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Event and time can't be empty"})
		return
	}

	user := User{}

	if err := db.Where(&User{Username: username}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	if err := db.Model(&user).Association("Schedule").Append(userSchedule{
		Event:     event,
		StartTime: startTime,
		EndTime:   endTime,
		Location:  location,
		Color:     color,
		Note:      note,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func udScheduleHdlr(c *gin.Context) {

}

func rmScheduleHdlr(c *gin.Context) {
	db, err := gorm.Open("mysql", "root:password@/sd?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer db.Close()

	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	event := c.PostForm("event")
	startTime, _ := time.Parse(time.RFC3339, c.PostForm("start"))
	endTime, _ := time.Parse(time.RFC3339, c.PostForm("end"))
	location := c.PostForm("location")
	color := c.PostForm("color")
	note := c.PostForm("note")

	user := User{}
	schedule := userSchedule{}

	if err := db.Where(&User{Username: username}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	if err := db.Where(&userSchedule{
		UserID:    user.ID,
		Event:     event,
		StartTime: startTime,
		EndTime:   endTime,
		Location:  location,
		Color:     color,
		Note:      note,
	}).First(&schedule).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "You don't have this schedule"})
		return
	}

	if err := db.Delete(&schedule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func getBackupHdlr(c *gin.Context) {
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

	if err := db.Where(&User{Username: username}).Preload("Backup").First(&user).Error; err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user.Backup)
}

func addBackupHdlr(c *gin.Context) {
	db, err := gorm.Open("mysql", "root:password@/sd?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer db.Close()

	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	title := c.PostForm("title")
	info := c.PostForm("info")
	importance, _ := strconv.Atoi(c.PostForm("importance"))

	if strings.Trim(title, " ") == "" || strings.Trim(info, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Title and content can't be empty"})
		return
	}

	user := User{}

	if err := db.Where(&User{Username: username}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "User not found"})
		return
	}

	if err := db.Model(&user).Association("Backup").Append(backup{
		Title:      title,
		Info:       info,
		Importance: importance,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func rmBackupHdlr(c *gin.Context) {
	db, err := gorm.Open("mysql", "root:password@/sd?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer db.Close()

	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	title := c.PostForm("title")
	info := c.PostForm("info")
	importance, _ := strconv.Atoi(c.PostForm("importance"))

	user := User{}
	bp := backup{}

	if err := db.Where(&User{Username: username}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	if err := db.Where(&backup{
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

func getGroupHdlr(c *gin.Context) {
	db, err := gorm.Open("mysql", "root:password@/sd?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

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
	db, err := gorm.Open("mysql", "root:password@/sd?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer db.Close()

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

func getMemberHdlr(c *gin.Context) {
	db, err := gorm.Open("mysql", "root:password@/sd?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer db.Close()

	name := c.PostForm("name")

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

func getGroupScheduleHdlr(c *gin.Context) {
	db, err := gorm.Open("mysql", "root:password@/sd?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer db.Close()

	name := c.PostForm("name")

	group := Group{}

	if err := db.Where(&Group{Name: name}).Preload("Schedule").First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Group not found"})
		return
	}

	c.JSON(http.StatusOK, group.Schedule)
}

func addGroupScheduleHdlr(c *gin.Context) {
	db, err := gorm.Open("mysql", "root:password@/sd?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer db.Close()

	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	name := c.PostForm("name")
	event := c.PostForm("event")
	startTime, _ := time.Parse(time.RFC3339, c.PostForm("start"))
	endTime, _ := time.Parse(time.RFC3339, c.PostForm("end"))
	location := c.PostForm("location")
	color := c.PostForm("color")
	note := c.PostForm("note")

	if strings.Trim(event, " ") == "" || strings.Trim(startTime.String(), " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Event and time can't be empty"})
		return
	}

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

	if err := db.Model(&group).Association("Schedule").Append(groupSchedule{
		SponsorID: user.ID,
		Event:     event,
		StartTime: startTime,
		EndTime:   endTime,
		Location:  location,
		Color:     color,
		Note:      note,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func rmGroupScheduleHdlr(c *gin.Context) {
	db, err := gorm.Open("mysql", "root:password@/sd?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	defer db.Close()

	claims := jwt.ExtractClaims(c)
	username := claims["username"].(string)

	name := c.PostForm("name")
	event := c.PostForm("event")
	startTime, _ := time.Parse(time.RFC3339, c.PostForm("start"))
	endTime, _ := time.Parse(time.RFC3339, c.PostForm("end"))
	location := c.PostForm("location")
	color := c.PostForm("color")
	note := c.PostForm("note")

	user := User{}
	group := Group{}
	schedule := groupSchedule{}

	if err := db.Where(&User{Username: username}).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	if err := db.Where(&Group{Name: name}).First(&group).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Group not found"})
		return
	}

	if err := db.Where(&groupSchedule{
		SponsorID: user.ID,
		Event:     event,
		StartTime: startTime,
		EndTime:   endTime,
		Location:  location,
		Color:     color,
		Note:      note,
	}).First(&schedule).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "You don't have this schedule"})
		return
	}

	if err := db.Delete(&schedule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}
