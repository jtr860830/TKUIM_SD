package main

import (
	"log"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jtr860830/LifePrint-Server/database"
	"github.com/jtr860830/LifePrint-Server/handlers"
)

func main() {
	db := database.GetDB()
	defer db.Close()

	route := gin.Default()
	route.Use(cors.Default())

	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:      "test zone",
		Key:        []byte("secret-key"),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*database.Payload); ok {
				return jwt.MapClaims{
					"id":       v.UserID,
					"username": v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals database.Login
			if err := c.Bind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := loginVals.Username
			password := loginVals.Password

			user := database.User{}
			if err := db.Where(&database.User{Username: username}).Find(&user).Error; err != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			if user.Password != password {
				return nil, jwt.ErrFailedAuthentication
			}

			return &database.Payload{
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
	route.GET("/logout", handlers.LogoutHdlr)
	route.POST("/register", handlers.RegisterHdlr)

	account := route.Group("/user", authMiddleware.MiddlewareFunc())
	{
		account.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/user/profile")
		})

		account.GET("/profile", handlers.ProfileHdlr)
		account.PATCH("/profile", handlers.UdProfileHdlr)
		account.POST("/chpasswd", handlers.ChpasswdHdlr)

		account.GET("/countdown", handlers.CdHdlr)

		account.GET("/map", handlers.AnlMap)
		account.GET("/map/weeks", handlers.AnlMapTimeWeek)
		account.GET("/map/months", handlers.AnlMapTimeMonth)
		account.GET("/map/years", handlers.AnlMapTimeYear)

		account.GET("/analysis/1", handlers.GetUserAnalysisHdlr)
		account.GET("/analysis/2", handlers.GetUserAnalysis2Hdlr)

		account.GET("/friends", handlers.GetFriendHdlr)
		account.POST("/friends", handlers.AddFriendHdlr)
		account.DELETE("/friends", handlers.RmFriendHdlr)

		account.GET("/schedules", handlers.GetScheduleHdlr)
		account.POST("/schedules", handlers.AddScheduleHdlr)
		account.PATCH("/schedules", handlers.UdScheduleHdlr)
		account.DELETE("/schedules", handlers.RmScheduleHdlr)

		account.GET("/backups", handlers.GetBackupHdlr)
		account.POST("/backups", handlers.AddBackupHdlr)
		account.DELETE("/backups", handlers.RmBackupHdlr)

		group := account.Group("/group")
		{
			group.GET("/", handlers.GetGroupHdlr)
			group.POST("/", handlers.CreateGroupHdlr)

			group.GET("/member", handlers.GetMemberHdlr)
			group.POST("/member", handlers.AddMemberHdlr)
			group.DELETE("/member", handlers.RmMemberHdlr)

			group.GET("/schedules", handlers.GetGroupScheduleHdlr)
			group.POST("/schedules", handlers.AddGroupScheduleHdlr)
			group.POST("/schedules/all", handlers.AddAllGroupScheduleHdlr)
			group.DELETE("/schedules", handlers.RmGroupScheduleHdlr)

			group.GET("/analysis/1", handlers.GetGroupAnalysisHdlr)
			group.GET("/analysis/2", handlers.GetGroupAnalysis2Hdlr)
		}
	}

	route.Run(":8080")
}
