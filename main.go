package main // import "github.com/jtr860830/SD-Backend"

import (
	"log"
	"net/http"
	"time"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {

	initDB()

	route := gin.Default()

	route.Use(cors.Default())

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
		account.PATCH("/profile", udProfileHdlr)
		account.POST("/chpasswd", chpasswdHdlr)

		account.GET("/countdown", cdHdlr)

		account.GET("/map", anlMap)
		account.GET("/map/weeks", anlMapTimeWeek)
		account.GET("/map/months", anlMapTimeMonth)
		account.GET("/map/years", anlMapTimeYear)

		account.GET("/analysis/1", getUserAnalysisHdlr)
		account.GET("/analysis/2", getUserAnalysis2Hdlr)

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
			group.POST("/schedules/all", addAllGroupScheduleHdlr)
			group.DELETE("/schedules", rmGroupScheduleHdlr)

			group.GET("/analysis/1", getGroupAnalysisHdlr)
			group.GET("/analysis/2", getGroupAnalysis2Hdlr)
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
		db.Set("gorm:table_options", "charset=utf8").AutoMigrate(&User{}, &Group{}, &userSchedule{}, &groupSchedule{}, &backup{}, &location{})
		db.Model(&userSchedule{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
		db.Model(&backup{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
		db.Model(&groupSchedule{}).AddForeignKey("group_id", "groups(id)", "RESTRICT", "RESTRICT")
	}
}
