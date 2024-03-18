package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"excel-receiver/entity"
	mockRepo "excel-receiver/mocks"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthorizationMiddleware(t *testing.T) {
	tests := []struct {
		name               string
		header             string
		loggerProviderMock func() (loggerProviderMock *mockRepo.ILogger)
		tokenRepoMock      func() (tokenRepoMock *mockRepo.TokenInterface)
		expectedStatus     int
	}{
		{
			name:   "Valid Token",
			header: "Bearer validtoken",
			loggerProviderMock: func() (loggerProviderMock *mockRepo.ILogger) {
				loggerProviderMock = mockRepo.NewILogger(t)
				return
			},
			tokenRepoMock: func() (tokenRepoMock *mockRepo.TokenInterface) {
				tokenRepoMock = mockRepo.NewTokenInterface(t)
				tokenRepoMock.On("GetTokenAuthentication", mock.Anything, "validtoken").Return(
					&entity.Token{
						ID:        1,
						Token:     "validtoken",
						ExpiredAt: time.Now().Add(7 * time.Hour),
						CreatedAt: time.Now(),
					}, nil,
				)
				return
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "No token",
			header: "InvalidToken",
			loggerProviderMock: func() (loggerProviderMock *mockRepo.ILogger) {
				loggerProviderMock = mockRepo.NewILogger(t)
				loggerProviderMock.On("WithFields", mock.Anything, mock.Anything).Return(logrus.NewEntry(logrus.New()))
				return
			},
			tokenRepoMock: func() (tokenRepoMock *mockRepo.TokenInterface) {
				tokenRepoMock = mockRepo.NewTokenInterface(t)
				return
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:   "Not bearer",
			header: "Basic token",
			loggerProviderMock: func() (loggerProviderMock *mockRepo.ILogger) {
				loggerProviderMock = mockRepo.NewILogger(t)
				loggerProviderMock.On("WithFields", mock.Anything, mock.Anything).Return(logrus.NewEntry(logrus.New()))
				return
			},
			tokenRepoMock: func() (tokenRepoMock *mockRepo.TokenInterface) {
				tokenRepoMock = mockRepo.NewTokenInterface(t)
				return
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid bearer token",
			header:         "Bearer invalidtoken",
			expectedStatus: http.StatusUnauthorized,
			loggerProviderMock: func() (loggerProviderMock *mockRepo.ILogger) {
				loggerProviderMock = mockRepo.NewILogger(t)
				loggerProviderMock.On("ErrorWithFields", mock.Anything, mock.Anything, mock.Anything)
				return
			},
			tokenRepoMock: func() (tokenRepoMock *mockRepo.TokenInterface) {
				tokenRepoMock = mockRepo.NewTokenInterface(t)
				tokenRepoMock.On("GetTokenAuthentication", mock.Anything, "invalidtoken").Return(nil, errors.New("test"))
				return
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := gin.HandlerFunc(func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
			})

			mockLogger := tt.loggerProviderMock()
			mockTokenRepo := tt.tokenRepoMock()

			router := gin.New()
			router.Use(Authorization(mockLogger, mockTokenRepo))
			router.GET("/test", handler)

			req, _ := http.NewRequest("GET", "/test", nil)
			req.Header.Set("Authorization", tt.header)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestLoggingMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		requestPath    string
		expectedStatus int
	}{
		{
			name:           "Logging Middleware",
			requestPath:    "/path",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := &mockRepo.ILogger{}
			mockLogger.On("InfoWithFields", mock.Anything, mock.Anything, mock.Anything)

			loggingMiddleware := LoggingMiddleware(mockLogger)

			req, err := http.NewRequest("GET", tt.requestPath, nil)
			assert.NoError(t, err)

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Set("req-id", "sampleRequestID")
			ctx.Request = req

			loggingMiddleware(ctx)

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockLogger.AssertExpectations(t)
		})
	}
}
