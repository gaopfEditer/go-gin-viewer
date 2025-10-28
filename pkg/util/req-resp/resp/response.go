package resp

import (
	"cambridge-hit.com/gin-base/activateserver/pkg/util"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/logger"
	"cambridge-hit.com/gin-base/activateserver/resource"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type Err struct {
	Code resource.RspCode
	Msg  string
}

type Response struct {
	Code    resource.RspCode `json:"code,omitempty"`
	Message string           `json:"message,omitempty"`
	Data    interface{}      `json:"data,omitempty"`
}

type PageResponse struct {
	Total int64       `json:"total,omitempty"`
	List  interface{} `json:"list,omitempty"`
}

func Success(context *gin.Context, data ...interface{}) {
	var r Response
	r.Code = http.StatusOK

	if len(data) > 0 {
		r.Data = data[0]
	}

	message := resource.CODE_SUCCESS.Msg()
	translatedMessage := util.T(message, util.GetLanguage(context))
	r.Message = translatedMessage

	context.JSON(http.StatusOK, r)
	logger.Info("Success", zap.String("response", translatedMessage), zap.String("path", context.Request.URL.Path))
}

func SuccessL(context *gin.Context, message string, data ...interface{}) {
	var r Response
	r.Code = http.StatusOK

	if len(data) > 0 {
		r.Data = data[0]
	}

	translatedMessage := util.T(message, util.GetLanguage(context))
	r.Message = translatedMessage

	context.JSON(http.StatusOK, r)
	logger.Info("Success", zap.String("response", translatedMessage), zap.String("path", context.Request.URL.Path))
}

func Error(context *gin.Context, errorCode resource.RspCode, msg ...string) {
	var r Response
	var message string
	r.Code = errorCode
	if len(msg) > 0 {
		message = msg[0]
	} else {
		message = errorCode.Msg()
	}

	message = util.T(message, util.GetLanguage(context))
	r.Message = message
	context.JSON(http.StatusOK, r)
	logger.Error("Failed", zap.String("response", r.Code.Msg()), zap.String("path", context.Request.URL.Path))
}

func ErrorERR(context *gin.Context, e Err) {
	var r Response
	var message string
	r.Code = e.Code
	if e.Msg != "" {
		message = e.Msg
	} else {
		message = resource.ERR_SERVER_BUSY.Msg()
	}

	message = util.T(message, util.GetLanguage(context))
	r.Message = message
	context.JSON(http.StatusOK, r)
	logger.Error("Failed", zap.String("response", r.Code.Msg()), zap.String("path", context.Request.URL.Path))
}

func ErrorWithData(context *gin.Context, errorCode resource.RspCode, msg string, data ...interface{}) {
	var r Response
	var message string
	r.Code = errorCode
	if len(msg) > 0 {
		message = msg
	} else {
		message = errorCode.Msg()
	}
	if len(data) > 0 {
		r.Data = data[0]
	}
	message = util.T(message, util.GetLanguage(context))
	r.Message = message
	context.JSON(http.StatusOK, r)
	logger.Error("Failed", zap.String("response", r.Code.Msg()), zap.String("path", context.Request.URL.Path))
}
