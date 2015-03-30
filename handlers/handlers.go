// Web api endpoint handlers.
package handlers

import (
	"github.com/blckur/blckur/static"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/session"
	"github.com/gin-gonic/gin"
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
	"strings"
)

var (
	staticStore *static.Store
	staticExpire = time.Now().Add(17531 * time.Hour).UTC().Format(time.RFC1123)
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
			panic(err)
		}

		sess, err := cook.GetSession(db)
		switch err.(type) {
		case nil:
		case *session.NotFoundError:
			sess = nil
			err = nil
		default:
			panic(err)
		}

		if required && sess == nil {
			c.AbortWithStatus(401)
		}

		c.Set("session", sess)
	}
}

func Proxy(c *gin.Context) {
	resp, err := http.Get(
		"http://localhost:8080" + c.Params.ByName("path"))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	c.Writer.Header().Add("Cache-Control",
		"no-cache, no-store, must-revalidate")
	c.Writer.Header().Add("Pragma", "no-cache")

	c.Data(200, resp.Header.Get("Content-Type"), body)
}

func Static(c *gin.Context) {
	path := c.Params.ByName("path")
	if path == "" || path == "/" {
		path = "/index.html"
	}

	path = "dart/build/web" + path


	file, ok := staticStore.HashFiles[path]
	if !ok {
		file, ok = staticStore.Files[path]
		if !ok {
			c.AbortWithStatus(404)
			return
		} else {
			// TODO Remove
			logrus.WithFields(logrus.Fields{
				"path": path,
			}).Warning("handlers: Non hash static lookup")
		}

		c.Writer.Header().Add("Cache-Control",
			"no-cache, no-store, must-revalidate")
		c.Writer.Header().Add("Pragma", "no-cache")
	} else {
		c.Writer.Header().Add("Cache-Control", "public, max-age=63113900")
		c.Writer.Header().Add("Expires", staticExpire)
	}

	if strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
		c.Writer.Header().Add("Content-Encoding", "gzip")
		c.Data(200, file.Type, file.GzipData)
	} else {
		c.Data(200, file.Type, file.Data)
	}
}

// Register all endpoint handlers
func Register(engine *gin.Engine) {
	var err error
	staticStore, err = static.NewStore("dart/build/web")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("handlers: Static store error")
		panic(err)
	}

	engine.Use(Limiter)
	engine.Use(gin.Recovery())

	dbGroup := engine.Group("")
	dbGroup.Use(Database)

	authGroup := dbGroup.Group("")
	authGroup.Use(Session(true))

	sessGroup := dbGroup.Group("")
	sessGroup.Use(Session(false))

	dbGroup.POST("/signup", signupPost)
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

	engine.GET("/s/*path", Proxy)
}
