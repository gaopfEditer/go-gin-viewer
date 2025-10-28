package ginstatic

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServeFileSystem interface {
	http.FileSystem
	Exists(prefix string, path string) bool
}

func ServeRoot(urlPrefix, root string) gin.HandlerFunc {
	return Serve(urlPrefix, LocalFile(root, false))
}

// Serve returns a middleware controller that serves resource files in the given directory.
func Serve(urlPrefix string, fs ServeFileSystem) gin.HandlerFunc {
	return ServeCached(urlPrefix, fs, 0)
}

// ServeCached returns a middleware controller that similar as Serve
// but with the Cache-Control Header set as passed in the cacheAge parameter
func ServeCached(urlPrefix string, fs ServeFileSystem, cacheAge uint) gin.HandlerFunc {
	fileserver := http.FileServer(fs)
	if urlPrefix != "" {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}
	return func(context *gin.Context) {
		if fs.Exists(urlPrefix, context.Request.URL.Path) {
			if cacheAge != 0 {
				context.Writer.Header().Add("Cache-Control", fmt.Sprintf("max-age=%d", cacheAge))
			}
			fileserver.ServeHTTP(context.Writer, context.Request)
			context.Abort()
		}
	}
}
