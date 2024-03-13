package service

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type SendRequestInterface interface {
	SendRequest(ctx *gin.Context, file *multipart.FileHeader, extensionFile string) (status string, err error)
}
