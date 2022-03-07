package middlewares

import (
	"compress/gzip"

	"github.com/gin-gonic/gin"
)

type gzipResponseWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func Gzip() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Content-Encoding") == "gzip" {
			if c.Request.Body != nil {
				gzreader, err := gzip.NewReader(c.Request.Body)
				if err != nil {
					return
				}
				gzreader.Close()
				c.Request.Body = gzreader
			}
		}

		//if c.Request.Header.Get("Accept-Encoding") == "gzip" {
		//	gzwriter := gzip.NewWriter(c.Writer)
		//
		//	c.Header("Content-Encoding", "gzip")
		//	c.Writer = &gzipResponseWriter{c.Writer, gzwriter}
		//
		//	defer func() {
		//		gzwriter.Close()
		//	}()
		//}

		c.Next()
	}
}
