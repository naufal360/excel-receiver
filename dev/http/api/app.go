package api

import (
	"excel-receiver/config"
	"excel-receiver/constant"
	"excel-receiver/http/api/ierr"
	"excel-receiver/middleware"
	"excel-receiver/provider"
	"excel-receiver/repository"
	"excel-receiver/service"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type App struct {
	log       provider.ILogger
	svc       service.SendRequestInterface
	tokenRepo repository.TokenInterface
}

func NewApp(log provider.ILogger, svc service.SendRequestInterface, tokenRepo repository.TokenInterface) *App {
	return &App{log: log, svc: svc, tokenRepo: tokenRepo}
}

func (a *App) CreateServer(address string) (*http.Server, error) {
	gin.SetMode(config.Configuration.Server.Mode)

	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(middleware.LoggingMiddleware(a.log))
	r.Use(middleware.Authorization(a.log, a.tokenRepo))
	r.POST(config.Configuration.Server.Endpoint.ExcelUpload, a.uploadFileRequest)

	r.GET("/ping", a.checkConnectivity)

	server := &http.Server{
		Addr:    address,
		Handler: r,
	}

	return server, nil
}

func (a *App) checkConnectivity(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (a *App) uploadFileRequest(ctx *gin.Context) {
	var (
		reqID           = ctx.Request.Context().Value(constant.RequestIDKey{}).(string)
		limitSize int64 = 128 * 1024 // 128kb
	)

	contentType := strings.Split(ctx.GetHeader("Content-Type"), ";")[0]
	if contentType != "multipart/form-data" {
		a.log.ErrorWithFields(
			provider.AppLog,
			map[string]interface{}{
				constant.ReqIDLog: reqID,
				"ERROR":           "invalid content-type",
			}, "failed to parse multipart form data request")

		ctx.JSON(ierr.MapResponse(ierr.NewF(constant.InvalidRequest, ""), reqID, "failed"))
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(ierr.MapResponse(ierr.NewF(constant.APIInternalError, ""), reqID, "failed"))
	}
	file := form.File["file"][0]
	extensionFile := filepath.Ext(file.Filename)

	if file.Size > limitSize {
		a.log.ErrorWithFields(
			provider.AppLog,
			map[string]interface{}{
				constant.ReqIDLog: reqID,
				"ERROR":           err,
			}, "failed to parse file size more than 128kb")

		ctx.JSON(ierr.MapResponse(ierr.NewF(constant.FileSizeLimit, ""), reqID, "failed"))
		return
	}

	status, err := a.svc.SendRequest(ctx, file, extensionFile)
	a.log.InfoWithFields(provider.AppLog, map[string]interface{}{
		constant.ReqIDLog: reqID,
		"ERROR":           err,
	}, "complete processing request")
	ctx.JSON(ierr.MapResponse(err, reqID, status))
}