package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"excel-receiver/config"
	"excel-receiver/constant"
	mockRepo "excel-receiver/mocks"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSendRequest(t *testing.T) {
	tempDir := t.TempDir()
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name                string
		contentType         string
		filename            string
		extension           string
		keyName             string
		rowsCounter         int
		isEmptyMandatory    bool
		loggerProviderMock  func() (loggerProviderMock *mockRepo.ILogger)
		requestProviderMock func() (requestProviderMock *mockRepo.RequestInterface)
		artemisRepoMock     func() (artemisRepoMock *mockRepo.QueueInterface)
		expectedStatus      string
	}{
		{
			name:             "success XLSX request",
			contentType:      "multipart/form-data",
			filename:         "test",
			extension:        ".xlsx",
			keyName:          "file",
			rowsCounter:      4,
			isEmptyMandatory: false,
			loggerProviderMock: func() (loggerProviderMock *mockRepo.ILogger) {
				loggerProviderMock = mockRepo.NewILogger(t)
				return
			},
			requestProviderMock: func() (requestProviderMock *mockRepo.RequestInterface) {
				requestProviderMock = mockRepo.NewRequestInterface(t)
				requestProviderMock.On("CreateRequest", mock.Anything, mock.Anything).Return("received", nil)
				return
			},
			artemisRepoMock: func() (artemisRepoMock *mockRepo.QueueInterface) {
				artemisRepoMock = mockRepo.NewQueueInterface(t)
				artemisRepoMock.On("ProduceQueue", mock.Anything, mock.Anything).Return(nil)
				return
			},
			expectedStatus: "received",
		},
		{
			name:             "success CSV request",
			contentType:      "multipart/form-data",
			filename:         "test",
			extension:        ".csv",
			keyName:          "file",
			rowsCounter:      2,
			isEmptyMandatory: false,
			loggerProviderMock: func() (loggerProviderMock *mockRepo.ILogger) {
				loggerProviderMock = mockRepo.NewILogger(t)
				return
			},
			requestProviderMock: func() (requestProviderMock *mockRepo.RequestInterface) {
				requestProviderMock = mockRepo.NewRequestInterface(t)
				requestProviderMock.On("CreateRequest", mock.Anything, mock.Anything).Return("received", nil)
				return
			},
			artemisRepoMock: func() (artemisRepoMock *mockRepo.QueueInterface) {
				artemisRepoMock = mockRepo.NewQueueInterface(t)
				artemisRepoMock.On("ProduceQueue", mock.Anything, mock.Anything).Return(nil)
				return
			},
			expectedStatus: "received",
		},
		{
			name:             "failed invalid extension request",
			contentType:      "multipart/form-data",
			filename:         "test",
			extension:        ".txt",
			keyName:          "file",
			rowsCounter:      2,
			isEmptyMandatory: false,
			loggerProviderMock: func() (loggerProviderMock *mockRepo.ILogger) {
				loggerProviderMock = mockRepo.NewILogger(t)
				return
			},
			requestProviderMock: func() (requestProviderMock *mockRepo.RequestInterface) {
				requestProviderMock = mockRepo.NewRequestInterface(t)
				return
			},
			artemisRepoMock: func() (artemisRepoMock *mockRepo.QueueInterface) {
				artemisRepoMock = mockRepo.NewQueueInterface(t)
				return
			},
			expectedStatus: "failed",
		},
		{
			name:             "failed csv invalid data rows",
			contentType:      "multipart/form-data",
			filename:         "test",
			extension:        ".csv",
			keyName:          "file",
			rowsCounter:      0,
			isEmptyMandatory: false,
			loggerProviderMock: func() (loggerProviderMock *mockRepo.ILogger) {
				loggerProviderMock = mockRepo.NewILogger(t)
				loggerProviderMock.On("ErrorWithFields", mock.Anything, mock.Anything, mock.Anything)
				return
			},
			requestProviderMock: func() (requestProviderMock *mockRepo.RequestInterface) {
				requestProviderMock = mockRepo.NewRequestInterface(t)
				return
			},
			artemisRepoMock: func() (artemisRepoMock *mockRepo.QueueInterface) {
				artemisRepoMock = mockRepo.NewQueueInterface(t)
				return
			},
			expectedStatus: "failed",
		},
		{
			name:             "failed xlsx invalid data rows",
			contentType:      "multipart/form-data",
			filename:         "test",
			extension:        ".xlsx",
			keyName:          "file",
			rowsCounter:      0,
			isEmptyMandatory: false,
			loggerProviderMock: func() (loggerProviderMock *mockRepo.ILogger) {
				loggerProviderMock = mockRepo.NewILogger(t)
				loggerProviderMock.On("ErrorWithFields", mock.Anything, mock.Anything, mock.Anything)
				return
			},
			requestProviderMock: func() (requestProviderMock *mockRepo.RequestInterface) {
				requestProviderMock = mockRepo.NewRequestInterface(t)
				return
			},
			artemisRepoMock: func() (artemisRepoMock *mockRepo.QueueInterface) {
				artemisRepoMock = mockRepo.NewQueueInterface(t)
				return
			},
			expectedStatus: "failed",
		},
		{
			name:             "failed csv empty row mandatory",
			contentType:      "multipart/form-data",
			filename:         "test",
			extension:        ".csv",
			keyName:          "file",
			rowsCounter:      1,
			isEmptyMandatory: true,
			loggerProviderMock: func() (loggerProviderMock *mockRepo.ILogger) {
				loggerProviderMock = mockRepo.NewILogger(t)
				return
			},
			requestProviderMock: func() (requestProviderMock *mockRepo.RequestInterface) {
				requestProviderMock = mockRepo.NewRequestInterface(t)
				return
			},
			artemisRepoMock: func() (artemisRepoMock *mockRepo.QueueInterface) {
				artemisRepoMock = mockRepo.NewQueueInterface(t)
				return
			},
			expectedStatus: "failed",
		},
	}

	config.LoadConfig("../")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := tt.loggerProviderMock()
			mockRequest := tt.requestProviderMock()
			mockQueue := tt.artemisRepoMock()

			filePath := fmt.Sprintf("%s%s", tt.filename, tt.extension)

			svc := NewSendRequestService(mockLogger, mockQueue, mockRequest)

			tempFile := filepath.Join(tempDir, filePath)

			switch tt.extension {
			case ".xlsx":
				createTestXLSXFile(t, tempFile, tt.rowsCounter, tt.isEmptyMandatory)
			case ".csv":
				createCsvFile(t, tempFile, tt.rowsCounter, tt.isEmptyMandatory)
			default:
				createTxtFile(t, tempFile, tt.rowsCounter)
			}

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
				req.Header.Add("Content-Type", tt.contentType)
			} else {
				req.Header.Add("Content-Type", writer.FormDataContentType())
			}

			w := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(w)
			req = req.WithContext(context.WithValue(req.Context(), constant.RequestIDKey{}, "abc123"))
			ctx.Request = req

			fileHeader, err := ctx.FormFile("file")
			assert.NoError(t, err)

			status, err := svc.SendRequest(ctx, fileHeader, tt.extension)

			if err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedStatus, status)

			// Assert mock calls
			mockLogger.AssertExpectations(t)
			mockRequest.AssertExpectations(t)
			mockQueue.AssertExpectations(t)
		})
	}
}

func TestSendRequest_checkMandatory(t *testing.T) {
	tests := []struct {
		name        string
		reqID       string
		titleCol    []string
		row         []string
		expectedErr bool
	}{
		{
			name:  "invalid row data",
			reqID: "abc123",
			titleCol: []string{
				"uniqid",
				"title",
				"description",
				"condition",
				"price",
				"color",
				"size",
				"age_group",
				"material",
				"weight_kg",
			},
			row: []string{
				"koko_xl",
				"Sample Title",
				"Sample Description",
				"New",
				"100000",
				"Red",
				"M",
				"Adult",
				"Cotton",
				"abcd", // weight != float
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rows [][]string
			rows = append(rows, tt.titleCol)
			rows = append(rows, tt.row)

			mockLogger := mockRepo.NewILogger(t)
			mockRequest := mockRepo.NewRequestInterface(t)
			mockQueue := mockRepo.NewQueueInterface(t)

			svc := NewSendRequestService(mockLogger, mockQueue, mockRequest)
			_, status, err := svc.checkMandatory(tt.reqID, rows)

			// Assertions
			if tt.expectedErr {
				assert.Equal(t, "failed", status)
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.expectedErr, err != nil)
			}
			// Assert mock calls
			mockLogger.AssertExpectations(t)
			mockRequest.AssertExpectations(t)
			mockQueue.AssertExpectations(t)
		})
	}
}

func TestSendRequest_processRow(t *testing.T) {
	tests := []struct {
		name        string
		reqID       string
		titleCol    []string
		row         []string
		expectedErr bool
	}{
		{
			name:  "failed price parse float",
			reqID: "abc123",
			titleCol: []string{
				"uniqid",
				"description",
				"condition",
				"price",
				"color",
				"size",
				"age_group",
				"material",
				"weight_kg",
			},
			row: []string{
				"koko_xl",
				"Sample Description",
				"New",
				"abc", // price != float
				"Red",
				"M",
				"Adult",
				"Cotton",
				"4",
			},
			expectedErr: true,
		},
		{
			name:  "failed weight parse float",
			reqID: "koko_xl",
			titleCol: []string{
				"uniqid",
				"description",
				"condition",
				"price",
				"color",
				"size",
				"age_group",
				"material",
				"weight_kg",
			},
			row: []string{
				"koko_xl",
				"Sample Description",
				"New",
				"100000",
				"Red",
				"M",
				"Adult",
				"Cotton",
				"abcd", // weight != float
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockLogger := mockRepo.NewILogger(t)
			mockRequest := mockRepo.NewRequestInterface(t)
			mockQueue := mockRepo.NewQueueInterface(t)

			svc := NewSendRequestService(mockLogger, mockQueue, mockRequest)
			_, err := svc.processRow(tt.reqID, tt.titleCol, tt.row)

			// Assertions
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.expectedErr, err != nil)
			}
			// Assert mock calls
			mockLogger.AssertExpectations(t)
			mockRequest.AssertExpectations(t)
			mockQueue.AssertExpectations(t)
		})
	}
}

func TestSendRequest_parseAndAssignFloat(t *testing.T) {
	tests := []struct {
		name          string
		col           string
		expectedValue float64
		expectedErr   bool
	}{
		{
			name:          "success parse float",
			col:           "8",
			expectedValue: 8,
			expectedErr:   false,
		},
		{
			name:          "failed parse float",
			col:           "abc",
			expectedValue: 8,
			expectedErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var target float64

			mockLogger := mockRepo.NewILogger(t)
			mockRequest := mockRepo.NewRequestInterface(t)
			mockQueue := mockRepo.NewQueueInterface(t)
			svc := NewSendRequestService(mockLogger, mockQueue, mockRequest)

			err := svc.parseAndAssignFloat(tt.col, &target)

			// Assertions
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.expectedErr, err != nil)
				assert.Equal(t, tt.expectedValue, target)
			}
			// Assert mock calls
			mockLogger.AssertExpectations(t)
			mockRequest.AssertExpectations(t)
			mockQueue.AssertExpectations(t)
		})
	}
}

// // Function to generate fake Excel file content
func createTestXLSXFile(t *testing.T, filepath string, loopCount int, isEmpty bool) {
	excelFile := excelize.NewFile()
	excelFile.NewSheet("Worksheet")
	excelFile.DeleteSheet("Sheet1")

	colBase := map[int]string{
		0: "uniqid",
		1: "title",
		2: "description",
		3: "condition",
		4: "price",
		5: "color",
		6: "size",
		7: "age_group",
		8: "material",
		9: "weight_kg",
	}
	colValue := map[int]string{
		0: "koko_xl",
		1: "Baju Koko China",
		2: "Baju Koko China Baru",
		3: "new",
		4: "200000",
		5: "blue",
		6: "XL",
		7: "adult",
		8: "cotton",
		9: "8",
	}

	var idxCol int
	for char := 'A'; char <= 'J'; char++ {
		cell := fmt.Sprintf("%c%d", char, 1)
		excelFile.SetCellValue("Worksheet", cell, colBase[idxCol])
		nextChar := char + 1
		if nextChar > 'H' {
			nextChar = 'A'
		}
		idxCol++
	}
	rowIdx := 2
	for i := 0; i < loopCount; i++ {
		idxCol = 0
		for char := 'A'; char <= 'J'; char++ {
			cell := fmt.Sprintf("%c%d", char, rowIdx)
			excelFile.SetCellValue("Worksheet", cell, colValue[idxCol])
			nextChar := char + 1
			if nextChar > 'H' {
				nextChar = 'A'
			}
			idxCol++
		}
		rowIdx++
	}
	err := excelFile.SaveAs(filepath)
	require.NoError(t, err)
}

func createCsvFile(t *testing.T, filePath string, loopCount int, isEmpty bool) {
	column := [][]string{
		{"uniqid", "title", "description", "condition", "price", "color", "size", "age_group", "material", "weight_kg"},
	}
	value := [][]string{
		{"koko_xl", "Baju Koko China", "Baju Koko China Baru", "new", "200000", "blue", "XL", "adult", "cooton", "8"},
	}
	emptyMandatoryColumn := [][]string{
		{"uniqid", "title", "condition", "price", "color", "size", "age_group", "material", "weight_kg"},
	}
	emptyMandatoryValue := [][]string{
		{"koko_xl", "Baju Koko China", "new", "200000", "blue", "XL", "adult", "cooton", "8"},
	}

	// Membuka file untuk penulisan. Jika file tidak ada, maka akan dibuat baru.
	file, err := os.Create(filePath)
	assert.NoError(t, err)
	defer file.Close()

	// Inisialisasi writer CSV
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Menulis data ke file CSV
	for i := 0; i <= loopCount; i++ {
		if i < 1 {
			if isEmpty {
				err = writer.WriteAll(emptyMandatoryColumn)
			} else {
				err = writer.WriteAll(column)
			}
		} else {
			if isEmpty {
				err = writer.WriteAll(emptyMandatoryValue)
			} else {
				err = writer.WriteAll(value)
			}
		}
		assert.NoError(t, err)
	}
}

func createTxtFile(t *testing.T, filePath string, loopCount int) {
	content := "Ini adalah contoh teks yang akan ditulis ke file."
	// Membuka file untuk penulisan. Jika file tidak ada, maka akan dibuat baru.
	file, err := os.Create(filePath)
	assert.NoError(t, err)

	defer file.Close()

	// Tulis konten ke file
	for i := 0; i < loopCount; i++ {
		_, err = file.WriteString(content)
		assert.NoError(t, err)
	}

}
