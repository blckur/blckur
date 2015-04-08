package handlers

import (
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/static"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
	"path/filepath"
)

type staticHandler struct {
	Static gin.HandlerFunc
	source string
	store  *static.Store
	expire string
}

func (s *staticHandler) proxy(c *gin.Context) {
	resp, err := http.Get(s.source + c.Params.ByName("path"))
	if err != nil {
		c.Fail(500, err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Fail(500, err)
		return
	}

	c.Writer.Header().Add("Cache-Control",
		"no-cache, no-store, must-revalidate")
	c.Writer.Header().Add("Pragma", "no-cache")

	c.Data(200, resp.Header.Get("Content-Type"), body)
}

func (s *staticHandler) local(c *gin.Context) {
	catch := true
	path := s.source + c.Params.ByName("path")

	file, ok := s.store.HashFiles[path]
	if !ok {
		file, ok = s.store.Files[path]
		if !ok {
			c.AbortWithStatus(404)
			return
		}
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

func (s *staticHandler) Index(c *gin.Context) {
	path := s.source + "/index.html"

	file, ok := s.store.Files[path]
	if !ok {
		c.AbortWithStatus(404)
		return
	}

	c.Writer.Header().Add("Cache-Control",
	"no-cache, no-store, must-revalidate")
	c.Writer.Header().Add("Pragma", "no-cache")

	if strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
		c.Writer.Header().Add("Content-Encoding", "gzip")
		c.Data(200, file.Type, file.GzipData)
	} else {
		c.Data(200, file.Type, file.Data)
	}
}

func newStaticHandler(source string) (handler *staticHandler) {
	source = filepath.Dir(source)

	handler = &staticHandler{
		source: source,
	}

	if source[:4] == "http" {
		handler.Static = handler.proxy
	} else {
		store, err := static.NewStore(source)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("handlers: Static store error")
			panic(err)
		}

		handler.store = store
		handler.Static = handler.local
	}

	return
}
