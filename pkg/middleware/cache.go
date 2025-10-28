package middleware

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"cambridge-hit.com/gin-base/activateserver/pkg/util/cache"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/logger"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/str"
	"cambridge-hit.com/gin-base/activateserver/resource"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Cache(fc *cache.FileCache) gin.HandlerFunc {
	return func(context *gin.Context) {
		pre := context.Request.URL.Path
		if !resource.Conf.Cache || !strings.Contains(pre, "imageview/wado/viewImage") {
			context.Next()
			return
		}
		key := pre + str.Md5To16(context.Request.RequestURI)

		respWrite := &responseCacheWriter{ResponseWriter: context.Writer, body: bytes.NewBuffer([]byte{})}
		context.Writer = respWrite

		if context.Request.Method == "GET" {
			// Get cache
			cacheData, err := fc.Get(key)
			if err == nil {
				fmt.Println("Cache hit")
				// Cache hit
				contentType := cacheData.(map[string]interface{})["contentType"].(string)
				body := cacheData.(map[string]interface{})["body"].([]byte)

				context.Writer.Header().Set("Content-Type", contentType)
				context.Writer.Write(body)
				context.Abort()
				return
			} else if !errors.Is(err, os.ErrNotExist) {
				logger.Error("reading cache file err", zap.Error(err))
			}

			// Cache miss
			context.Next()

			// Set cache with resp
			contentType := context.Writer.Header().Get("Content-Type")
			if contentType == "" {
				contentType = "application/json"
			}

			cacheData = map[string]interface{}{
				"contentType": contentType,
				"body":        respWrite.body.Bytes(),
			}

			if err := fc.Set(key, cacheData); err != nil {
				logger.Error("setting cache file err", zap.Error(err))
			}
			return
		} else {
			// Delete cache
			context.Next()
			if err := fc.DelByPrefix(pre); err != nil {
				logger.Error("cache deletion err", zap.Error(err))
			}
			return
		}
	}
}

// responseCacheWriter
type responseCacheWriter struct {
	gin.ResponseWriter

	body *bytes.Buffer
}

func (w *responseCacheWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *responseCacheWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
