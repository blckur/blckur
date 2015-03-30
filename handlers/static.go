package handlers

import (
	"github.com/blckur/blckur/static"
	"github.com/gin-gonic/gin"
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

type staticHandler struct {
	source string
	store *static.Store
	expire string
}

func (s *staticHandler) Proxy(c *gin.Context) {
	resp, err := http.Get(s.source + c.Params.ByName("path"))
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

func (s *staticHandler) Static(c *gin.Context) {
	catch := true

	path := c.Params.ByName("path")
	if path == "" || path == "/" {
		path = "/index.html"
	}
	path = s.source + path

	file, ok := s.store.HashFiles[path]
	if !ok {
		file, ok = s.store.Files[path]
		if !ok {
			c.AbortWithStatus(404)
			return
		}
//		} else {
//			logrus.WithFields(logrus.Fields{
//				"path": path,
//			}).Warning("handlers: Non hash static lookup")
//		}
		catch = false
	}

	if catch {
		c.Writer.Header().Add("Cache-Control", "public, max-age=63113900")
		c.Writer.Header().Add("Expires", s.expire)
	} else {
		c.Writer.Header().Add("Cache-Control",
			"no-cache, no-store, must-revalidate")
		c.Writer.Header().Add("Pragma", "no-cache")
	}

	if strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
		c.Writer.Header().Add("Content-Encoding", "gzip")
		c.Data(200, file.Type, file.GzipData)
	} else {
		c.Data(200, file.Type, file.Data)
	}
}

func newStaticHandler(source string) (handlerFunc gin.HandlerFunc) {
	handler := &staticHandler{
		source: source,
	}

	if source[:4] == "http" {
		handlerFunc = handler.Proxy
	} else {
		store, err := static.NewStore(source)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("handlers: Static store error")
			panic(err)
		}

		handler.store = store
		handlerFunc = handler.Static
	}

	return
}
