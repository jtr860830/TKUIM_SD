package main // import "github.com/jtr860830/SD-Backend"

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	route := gin.Default()
	store := cookie.NewStore([]byte("secret-string"))
	route.Use(sessions.Sessions("session", store))

	route.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Home Page")
	})
	route.POST("/login", loginHandler)
	route.GET("/logout", logoutHandler)
	route.POST("/register", registerHandler)

	account := route.Group("/user")
	{
		account.GET("/", func(c *gin.Context) {
			c.Redirect(301, "/user/profile")
		})
		account.GET("/profile", profileHandler)
	}
	account.Use(auth())

	route.Run(":8080")
}

func auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			c.String(http.StatusNotAcceptable, "You should not pass!")
			log.Println("A strangers attempted to log in!")
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
	if username == "hello" && password == "itsme" {
		session.Set("user", username)
		err := session.Save()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate session token"})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user"})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
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

}

func profileHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	c.JSON(http.StatusOK, user)
}
