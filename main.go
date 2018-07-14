package main // import "github.com/jtr860830/SD-Backend"

import (
	//"log"
	"net/http"
	//"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	route := gin.Default()
	store := cookie.NewStore([]byte("secret-string"))
	route.Use(sessions.Sessions("mysession", store))

	route.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Home Page")
	})
	route.POST("/login", loginHandler)
	route.GET("/logout", logoutHandler)
	route.POST("/register", registerHandler)

	route.Run(":8080")
}

func loginHandler(c *gin.Context) {

}

func logoutHandler(c *gin.Context) {

}

func registerHandler(c *gin.Context) {

}
