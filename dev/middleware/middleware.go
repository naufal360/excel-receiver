package middleware

import (
	"context"
	"excel-receiver/constant"
	"excel-receiver/dto/response"
	"excel-receiver/provider"
	"excel-receiver/repository"
	"excel-receiver/util"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Authorization(logger provider.ILogger, tokenRepo repository.TokenInterface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqID := util.GenerateReqID()
		ctx.Set("request-id", reqID)
		ReqIDLog := constant.ReqIDLog

		token := strings.Split(ctx.GetHeader("Authorization"), " ")

		unauthorizedResp := response.BaseResponse{
			RequestID: reqID,
			Status:    "failed",
			Detail: response.Detail{
				Code:    string(constant.APIUnauthorized),
				Message: "unauthorized",
			},
		}

		if len(token) != 2 {
			logger.WithFields(provider.AppLog, logrus.Fields{ReqIDLog: reqID}).Error("middleware : Invalid bearer token")
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				unauthorizedResp,
			)
			return
		}

		if strings.ToLower(token[0]) != "bearer" {
			logger.WithFields(provider.AppLog, logrus.Fields{ReqIDLog: reqID}).Error("middleware : Invalid bearer token")
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				unauthorizedResp,
			)
			return
		}

		bearerToken := token[1]

		dataToken, err := tokenRepo.GetTokenAuthentication(ctx, bearerToken)
		if err != nil {
			logger.ErrorWithFields(provider.AppLog,
				map[string]interface{}{
					"ERROR":  err,
					ReqIDLog: reqID,
					"TOKEN":  token,
				},
				"middleware: failed get user authentication")
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				unauthorizedResp,
			)
			return
		}

		nCtx := context.WithValue(ctx.Request.Context(), constant.RequestIDKey{}, reqID)
		ctx.Request = ctx.Request.WithContext(nCtx)

		ctx.Set("tokenAuth", dataToken)
		ctx.Next()
	}
}

func LoggingMiddleware(logger provider.ILogger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqID := util.GenerateReqID()
		ReqIDLog := constant.ReqIDLog
		nCtx := context.WithValue(ctx.Request.Context(), constant.RequestIDKey{}, reqID)
		ctx.Request = ctx.Request.WithContext(nCtx)

		// Starting time
		startTime := time.Now()

		// Processing request
		ctx.Next()

		// End Time
		endTime := time.Now()

		// execution time
		latencyTime := endTime.Sub(startTime)

		// Request method
		reqMethod := ctx.Request.Method

		// Request route
		reqUri := ctx.Request.RequestURI

		// status code
		statusCode := ctx.Writer.Status()

		// Request IP
		clientIP := ctx.ClientIP()

		logger.InfoWithFields(provider.AppLog, map[string]interface{}{
			"METHOD":    reqMethod,
			"URI":       reqUri,
			"STATUS":    statusCode,
			"LATENCY":   latencyTime,
			"CLIENT_IP": clientIP,
			ReqIDLog:    reqID,
		}, "HTTP REQUEST")

		ctx.Next()
	}
}
