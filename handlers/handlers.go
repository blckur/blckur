// Web api endpoint handlers.
package handlers

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/session"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Limit size of request body
func Limiter(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 4096)
}

// Get database from session
func Database(c *gin.Context) {
	db := database.GetDatabase()
	c.Set("db", db)
	c.Next()
	db.Close()
}

// Create new session and store in cookies
func Session(required bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := c.MustGet("db").(*database.Database)

		cook, err := session.GetCookie(c)
		if err != nil {
			c.Fail(500, err)
			return
		}

		sess, err := cook.GetSession(db)
		switch err.(type) {
		case nil:
		case *session.NotFoundError:
			sess = nil
			err = nil
		default:
			c.Fail(500, err)
			return
		}

		if required && sess == nil {
			c.AbortWithStatus(401)
			return
		}

		c.Set("session", sess)
	}
}

// Register all endpoint handlers
func Register(engine *gin.Engine, source string) {
	engine.Use(Limiter)
	engine.Use(gin.Recovery())

	dbGroup := engine.Group("")
	dbGroup.Use(Database)

	authGroup := dbGroup.Group("")
	authGroup.Use(Session(true))

	sessGroup := dbGroup.Group("")
	sessGroup.Use(Session(false))

	dbGroup.POST("/signup", signupPost)
	dbGroup.GET("/reset/:key", resetGet)
	dbGroup.PUT("/reset", resetPut)
	sessGroup.POST("/login", loginPost)
	authGroup.DELETE("/session", sessionDelete)

	authGroup.GET("/account_types", accountsTypesGet)
	authGroup.GET("/filter_types/:type", filterTypesGet)

	authGroup.GET("/accounts", accountsGet)
	authGroup.PUT("/accounts/:id", accountsPut)
	authGroup.DELETE("/accounts/:id", accountsDel)
	authGroup.POST("/accounts/:type", accountsPost)

	dbGroup.GET("/callback/:type", callbackGet)

	authGroup.GET("/user", userGet)
	authGroup.PUT("/user", userPut)

	authGroup.GET("/events", eventGet)

	authGroup.GET("/notifications", notificationGet)
	authGroup.PUT("/notifications/:id", notificationPut)
	authGroup.DELETE("/notifications/:id", notificationDel)

	staticHand := newStaticHandler(source)
	engine.GET("/", staticHand.Index)
	engine.GET("/s/*path", staticHand.Static)
}
