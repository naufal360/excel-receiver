package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"excel-receiver/constant"
	"excel-receiver/dto/response"
	mockRepo "excel-receiver/mocks"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAppCreateServer(t *testing.T) {
	// Mock the logger and service using mockery generated mocks
	mockLogger := &mockRepo.ILogger{}
	mockToken := &mockRepo.TokenInterface{}
	mockService := &mockRepo.SendRequestInterface{}

	app := NewApp(mockLogger, mockService, mockToken)
	server, err := app.CreateServer(":8080")
	assert.NoError(t, err)
	assert.NotNil(t, server)
}

func TestAppCheckConnectivity(t *testing.T) {
	// Mock the logger and service using mockery generated mocks
	mockLogger := &mockRepo.ILogger{}
	mockToken := &mockRepo.TokenInterface{}
	mockService := &mockRepo.SendRequestInterface{}

	app := NewApp(mockLogger, mockService, mockToken)

	// Create a request to the /ping endpoint
	req, _ := http.NewRequest("GET", "/ping", nil)

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Create a Gin context
	r := gin.Default()
	r.GET("/ping", app.checkConnectivity)

	// Serve the request
	r.ServeHTTP(w, req)

	// Check the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func TestAppUploadFileRequest(t *testing.T) {
	tempDir := t.TempDir()
	testJson := map[string]string{
		"col1": "value1",
		"col2": "value2",
	}

	tests := []struct {
		name                string
		contentType         string
		filename            string
		keyName             string
		rowsCounter         int
		loggerProviderMock  func() (loggerProviderMock *mockRepo.ILogger)
		serviceProviderMock func() (serviceProviderMock *mockRepo.SendRequestInterface)
		tokenRepoMock       func() (tokenRepoMock *mockRepo.TokenInterface)
		expectedCode        string
	}{
		{
			name:        "valid request",
			contentType: "multipart/form-data",
			filename:    "test.xlsx",
			keyName:     "file",
			rowsCounter: 2,
			loggerProviderMock: func() (loggerProviderMock *mockRepo.ILogger) {
				loggerProviderMock = mockRepo.NewILogger(t)
				loggerProviderMock.On("InfoWithFields", mock.Anything, mock.Anything, mock.Anything)
				return
			},
			serviceProviderMock: func() (serviceProviderMock *mockRepo.SendRequestInterface) {
				serviceProviderMock = mockRepo.NewSendRequestInterface(t)
				serviceProviderMock.On("SendRequest", mock.Anything, mock.Anything, mock.Anything).Return("received", nil)
				return
			},
			tokenRepoMock: func() (tokenRepoMock *mockRepo.TokenInterface) {
				tokenRepoMock = mockRepo.NewTokenInterface(t)
				return
			},
			expectedCode: "100",
		},
		{
			name:        "internal server error",
			contentType: "multipart/form-data",
			filename:    "test.xlsx",
			keyName:     "",
			rowsCounter: 2,
			loggerProviderMock: func() (loggerProviderMock *mockRepo.ILogger) {
				loggerProviderMock = mockRepo.NewILogger(t)
				return
			},
			serviceProviderMock: func() (serviceProviderMock *mockRepo.SendRequestInterface) {
				serviceProviderMock = mockRepo.NewSendRequestInterface(t)
				return
			},
			tokenRepoMock: func() (tokenRepoMock *mockRepo.TokenInterface) {
				tokenRepoMock = mockRepo.NewTokenInterface(t)
				return
			},
			expectedCode: "500",
		},
		{
			name:        "invalid request",
			contentType: "application/json",
			filename:    "test.xlsx",
			keyName:     "file",
			rowsCounter: 2,
			loggerProviderMock: func() (loggerProviderMock *mockRepo.ILogger) {
				loggerProviderMock = mockRepo.NewILogger(t)
				loggerProviderMock.On("ErrorWithFields", mock.Anything, mock.Anything, mock.Anything)
				return
			},
			serviceProviderMock: func() (serviceProviderMock *mockRepo.SendRequestInterface) {
				serviceProviderMock = mockRepo.NewSendRequestInterface(t)
				return
			},
			tokenRepoMock: func() (tokenRepoMock *mockRepo.TokenInterface) {
				tokenRepoMock = mockRepo.NewTokenInterface(t)
				return
			},
			expectedCode: "101",
		},
		{
			name:        "file size more than 128kb",
			contentType: "multipart/form-data",
			filename:    "test.xlsx",
			keyName:     "file",
			rowsCounter: 6000, //135kb filesize
			loggerProviderMock: func() (loggerProviderMock *mockRepo.ILogger) {
				loggerProviderMock = mockRepo.NewILogger(t)
				loggerProviderMock.On("ErrorWithFields", mock.Anything, mock.Anything, mock.Anything)
				return
			},
			serviceProviderMock: func() (serviceProviderMock *mockRepo.SendRequestInterface) {
				serviceProviderMock = mockRepo.NewSendRequestInterface(t)
				return
			},
			tokenRepoMock: func() (tokenRepoMock *mockRepo.TokenInterface) {
				tokenRepoMock = mockRepo.NewTokenInterface(t)
				return
			},
			expectedCode: "104",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := tt.loggerProviderMock()
			mockService := tt.serviceProviderMock()
			mockToken := tt.tokenRepoMock()

			app := NewApp(mockLogger, mockService, mockToken)

			tempFile := filepath.Join(tempDir, tt.filename)
			createTestXLSXFile(t, tempFile, tt.rowsCounter)

			file, err := os.Open(tempFile)
			require.NoError(t, err)
			defer file.Close()

			bodyForm := new(bytes.Buffer)
			writer := multipart.NewWriter(bodyForm)

			fileWriter, err := writer.CreateFormFile(tt.keyName, filepath.Base(tempFile))
			assert.NoError(t, err)
			_, err = io.Copy(fileWriter, file)
			require.NoError(t, err)
			writer.Close()

			req, err := http.NewRequest(http.MethodPost, "/excel-upload", bodyForm)
			assert.NoError(t, err)

			if tt.contentType != "multipart/form-data" {
				jsonData, err := json.Marshal(testJson)
				assert.NoError(t, err)

				byteJson := []byte(jsonData)
				bodyJson := bytes.NewBuffer(byteJson)

				req, err = http.NewRequest(http.MethodPost, "/excel-upload", bodyJson)
				assert.NoError(t, err)
				req.Header.Add("Content-Type", tt.contentType)
			} else {
				req.Header.Add("Content-Type", writer.FormDataContentType())
			}

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			req = req.WithContext(context.WithValue(req.Context(), constant.RequestIDKey{}, "abc123"))
			ctx.Request = req

			app.uploadFileRequest(ctx)

			// Read response from recorder used during the request execution
			// Validate the response
			var response response.BaseResponse
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCode, response.Detail.Code)

			// Assert mock calls
			mockLogger.AssertExpectations(t)
			mockService.AssertExpectations(t)
			mockToken.AssertExpectations(t)
		})
	}
}

// Function to generate fake Excel file content
func createTestXLSXFile(t *testing.T, filepath string, loopCount int) {
	excelFile := excelize.NewFile()
	excelFile.NewSheet("Worksheet")
	excelFile.DeleteSheet("Sheet1")

	for i := 0; i < loopCount; i++ {
		for char := 'A'; char <= 'J'; char++ {
			cell := fmt.Sprintf("%c%d", char, i+1)
			if i < 1 {
				excelFile.SetCellValue("Worksheet", cell, "uniqid")
				excelFile.SetCellValue("Worksheet", cell, "title")
				excelFile.SetCellValue("Worksheet", cell, "description")
				excelFile.SetCellValue("Worksheet", cell, "condition")
				excelFile.SetCellValue("Worksheet", cell, "price")
				excelFile.SetCellValue("Worksheet", cell, "color")
				excelFile.SetCellValue("Worksheet", cell, "size")
				excelFile.SetCellValue("Worksheet", cell, "age_group")
				excelFile.SetCellValue("Worksheet", cell, "materiall")
				excelFile.SetCellValue("Worksheet", cell, "weight_kg")
			} else {
				excelFile.SetCellValue("Worksheet", cell, "koko_xl")
				excelFile.SetCellValue("Worksheet", cell, "Baju Koko China")
				excelFile.SetCellValue("Worksheet", cell, "Baju Koko China Baru")
				excelFile.SetCellValue("Worksheet", cell, "new")
				excelFile.SetCellValue("Worksheet", cell, "200000")
				excelFile.SetCellValue("Worksheet", cell, "blue")
				excelFile.SetCellValue("Worksheet", cell, "XL")
				excelFile.SetCellValue("Worksheet", cell, "cotton")
				excelFile.SetCellValue("Worksheet", cell, "8")
			}
			nextChar := char + 1
			if nextChar > 'H' {
				nextChar = 'A'
			}
		}
	}
	err := excelFile.SaveAs(filepath)
	require.NoError(t, err)
}
