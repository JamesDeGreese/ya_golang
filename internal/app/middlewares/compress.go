package middlewares

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type gzipResponseWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func Gzip() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Accept-Encoding") != "gzip" {
			return
		}

		if c.Request.Body != nil {
			gzreader, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.String(http.StatusInternalServerError, "")
				return
			}
			gzreader.Close()
			c.Request.Body = ioutil.NopCloser(gzreader)
		}

		gzwriter := gzip.NewWriter(c.Writer)

		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding")
		c.Writer = &gzipResponseWriter{c.Writer, gzwriter}

		defer func() {
			gzwriter.Close()
			c.Header("Content-Length", fmt.Sprint(c.Writer.Size()))
		}()
		c.Next()
	}
}
