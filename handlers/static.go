package handlers

import (
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/session"
	"github.com/blckur/blckur/static"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

const (
	Local = 0
	Proxy = 1
)

type staticHandler struct {
	Type   int
	Static gin.HandlerFunc
	Index  gin.HandlerFunc
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

	if c.Params.ByName("path") == "/" || c.Params.ByName("path") == "" {
		sess := c.MustGet("session").(*session.Session)

		var page string
		if sess != nil {
			page = "feed"
		} else {
			page = "login"
		}

		body = []byte(strings.Replace(string(body), "default-page", page, 1))
	}

	c.Writer.Header().Add("Cache-Control",
		"no-cache, no-store, must-revalidate")
	c.Writer.Header().Add("Pragma", "no-cache")

	c.Data(200, resp.Header.Get("Content-Type"), body)
}

func (s *staticHandler) proxyIndex(c *gin.Context) {
	c.Redirect(302, "/s/")
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

func (s *staticHandler) localIndex(c *gin.Context) {
	sess := c.MustGet("session").(*session.Session)
	path := s.source

	if sess != nil {
		path += "/index-feed.html"
	} else {
		path += "/index-login.html"
	}

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
	handler = &staticHandler{}

	if source[:4] == "http" {
		handler.Type = Proxy
		handler.Static = handler.proxy
		handler.Index = handler.proxyIndex
	} else {
		source = filepath.Dir(source)
		store, err := static.NewStore(source)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("handlers: Static store error")
			panic(err)
		}

		handler.store = store
		handler.Type = Local
		handler.Static = handler.local
		handler.Index = handler.localIndex
	}

	handler.source = source

	return
}
