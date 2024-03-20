package testing

import (
	"bytes"
	"encoding/json"
	"excel-receiver/config"
	"excel-receiver/dto/response"
	"excel-receiver/testing/test_util"
	"fmt"
	"mime/multipart"
	"net/http"
	"testing"
	"time"

	"github.com/go-stomp/stomp/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(testSvc))
}

func (suite *testSvc) TestAPI() {
	t := suite.T()

	type want struct {
		err                 bool
		code                int
		response            response.BaseResponse
		expectedLengthQueue int
	}

	tests := []struct {
		name      string
		header    map[string]string
		path      string
		filePath  string
		typeMedia test_util.MultipartRequest
		want      want
		reqString string
	}{
		{
			name:     "No Auth",
			path:     "/excel-upload",
			filePath: "./storage/final_test.csv",
			typeMedia: test_util.MultipartRequest{
				Field: "type",
				Value: "file",
			},
			want: want{
				err:  true,
				code: http.StatusUnauthorized,
				response: response.BaseResponse{
					RequestID: "abc123",
					Status:    "failed",
					Detail: response.Detail{
						Code:    "401",
						Message: "unauthorized",
					},
				},
			},
		},
		{
			name: "Basic Auth",
			header: map[string]string{
				"Authorization": "Basic invalid",
			},
			path:     "/excel-upload",
			filePath: "./storage/final_test.csv",
			typeMedia: test_util.MultipartRequest{
				Field: "type",
				Value: "file",
			},
			want: want{
				err:  true,
				code: http.StatusUnauthorized,
				response: response.BaseResponse{
					RequestID: "abc123",
					Status:    "failed",
					Detail: response.Detail{
						Code:    "401",
						Message: "unauthorized",
					},
				},
			},
		},
		{
			name: "Token Not Found",
			header: map[string]string{
				"Authorization": "Bearer notfound",
			},
			path:     "/excel-upload",
			filePath: "./storage/final_test.csv",
			typeMedia: test_util.MultipartRequest{
				Field: "type",
				Value: "file",
			},
			want: want{
				err:  true,
				code: http.StatusUnauthorized,
				response: response.BaseResponse{
					RequestID: "abc123",
					Status:    "failed",
					Detail: response.Detail{
						Code:    "401",
						Message: "unauthorized",
					},
				},
			},
		},
		{
			name: "Session Token Expired",
			header: map[string]string{
				"Authorization": "Bearer expired_token",
			},
			path:     "/excel-upload",
			filePath: "./storage/final_test.csv",
			typeMedia: test_util.MultipartRequest{
				Field: "type",
				Value: "file",
			},
			want: want{
				err:  true,
				code: http.StatusUnauthorized,
				response: response.BaseResponse{
					RequestID: "abc123",
					Status:    "failed",
					Detail: response.Detail{
						Code:    "401",
						Message: "unauthorized",
					},
				},
			},
		},
		{
			name: "Success Receive Request XLSX 1 Data",
			header: map[string]string{
				"Authorization": "Bearer token_test",
			},
			path:     "/excel-upload",
			filePath: "./storage/final_test.xlsx",
			typeMedia: test_util.MultipartRequest{
				Field: "type",
				Value: "file",
			},
			want: want{
				err:  false,
				code: http.StatusAccepted,
				response: response.BaseResponse{
					RequestID: "abc123",
					Status:    "received",
					Detail: response.Detail{
						Code:    "100",
						Message: "success",
					},
				},
				expectedLengthQueue: 1,
			},
		},
		{
			name: "Success Receive Request CSV 1 Data",
			header: map[string]string{
				"Authorization": "Bearer token_test",
			},
			path:     "/excel-upload",
			filePath: "./storage/final_test.csv",
			typeMedia: test_util.MultipartRequest{
				Field: "type",
				Value: "file",
			},
			want: want{
				err:  false,
				code: http.StatusAccepted,
				response: response.BaseResponse{
					RequestID: "abc123",
					Status:    "received",
					Detail: response.Detail{
						Code:    "100",
						Message: "success",
					},
				},
				expectedLengthQueue: 1,
			},
		},
		{
			name: "Success Receive Request XLSX 3 Data",
			header: map[string]string{
				"Authorization": "Bearer token_test",
			},
			path:     "/excel-upload",
			filePath: "./storage/final_test_valid.xlsx",
			typeMedia: test_util.MultipartRequest{
				Field: "type",
				Value: "file",
			},
			want: want{
				err:  false,
				code: http.StatusAccepted,
				response: response.BaseResponse{
					RequestID: "abc123",
					Status:    "received",
					Detail: response.Detail{
						Code:    "100",
						Message: "success",
					},
				},
				expectedLengthQueue: 3,
			},
		},
		{
			name: "Success Receive Request CSV 3 Data",
			header: map[string]string{
				"Authorization": "Bearer token_test",
			},
			path:     "/excel-upload",
			filePath: "./storage/final_test_valid.csv",
			typeMedia: test_util.MultipartRequest{
				Field: "type",
				Value: "file",
			},
			want: want{
				err:  false,
				code: http.StatusAccepted,
				response: response.BaseResponse{
					RequestID: "abc123",
					Status:    "received",
					Detail: response.Detail{
						Code:    "100",
						Message: "success",
					},
				},
				expectedLengthQueue: 3,
			},
		},
		{
			name: "Success Receive Request CSV Only Mandatory Row",
			header: map[string]string{
				"Authorization": "Bearer token_test",
			},
			path:     "/excel-upload",
			filePath: "./storage/final_test_only_mandatory.csv",
			typeMedia: test_util.MultipartRequest{
				Field: "type",
				Value: "file",
			},
			want: want{
				err:  false,
				code: http.StatusAccepted,
				response: response.BaseResponse{
					RequestID: "abc123",
					Status:    "received",
					Detail: response.Detail{
						Code:    "100",
						Message: "success",
					},
				},
				expectedLengthQueue: 1,
			},
		},
		{
			name: "Success Receive Request XLSX Only Mandatory Row",
			header: map[string]string{
				"Authorization": "Bearer token_test",
			},
			path:     "/excel-upload",
			filePath: "./storage/final_test_only_mandatory.xlsx",
			typeMedia: test_util.MultipartRequest{
				Field: "type",
				Value: "file",
			},
			want: want{
				err:  false,
				code: http.StatusAccepted,
				response: response.BaseResponse{
					RequestID: "abc123",
					Status:    "received",
					Detail: response.Detail{
						Code:    "100",
						Message: "success",
					},
				},
				expectedLengthQueue: 1,
			},
		},
		{
			name: "Success Receive Request CSV Exist Column Nil Value",
			header: map[string]string{
				"Authorization": "Bearer token_test",
			},
			path:     "/excel-upload",
			filePath: "./storage/final_test_value_nil.csv",
			typeMedia: test_util.MultipartRequest{
				Field: "type",
				Value: "file",
			},
			want: want{
				err:  false,
				code: http.StatusAccepted,
				response: response.BaseResponse{
					RequestID: "abc123",
					Status:    "received",
					Detail: response.Detail{
						Code:    "100",
						Message: "success",
					},
				},
				expectedLengthQueue: 1,
			},
		},
		{
			name: "Success Receive Request XLSX Exist Column Nil Value",
			header: map[string]string{
				"Authorization": "Bearer token_test",
			},
			path:     "/excel-upload",
			filePath: "./storage/final_test_value_nil.xlsx",
			typeMedia: test_util.MultipartRequest{
				Field: "type",
				Value: "file",
			},
			want: want{
				err:  false,
				code: http.StatusAccepted,
				response: response.BaseResponse{
					RequestID: "abc123",
					Status:    "received",
					Detail: response.Detail{
						Code:    "100",
						Message: "success",
					},
				},
				expectedLengthQueue: 1,
			},
		},
		{
			name: "Invalid Content-type", //json
			header: map[string]string{
				"Authorization": "Bearer token_test",
			},
			path: "/excel-upload",
			want: want{
				err:  true,
				code: http.StatusBadRequest,
				response: response.BaseResponse{
					RequestID: "abc123",
					Status:    "failed",
					Detail: response.Detail{
						Code:    "101",
						Message: "invalid request",
					},
				},
			},
		},
		{
			name: "Invalid Extension File",
			header: map[string]string{
				"Authorization": "Bearer token_test",
			},
			path:     "/excel-upload",
			filePath: "./storage/final_test.txt",
			typeMedia: test_util.MultipartRequest{
				Field: "type",
				Value: "file",
			},
			want: want{
				err:  true,
				code: http.StatusBadRequest,
				response: response.BaseResponse{
					RequestID: "abc123",
					Status:    "failed",
					Detail: response.Detail{
						Code:    "101",
						Message: "invalid request",
					},
				},
			},
		},
		{
			name: "Invalid File Size",
			header: map[string]string{
				"Authorization": "Bearer token_test",
			},
			path:     "/excel-upload",
			filePath: "./storage/final_test_exceed_size.xlsx",
			typeMedia: test_util.MultipartRequest{
				Field: "type",
				Value: "file",
			},
			want: want{
				err:  true,
				code: http.StatusBadRequest,
				response: response.BaseResponse{
					RequestID: "abc123",
					Status:    "failed",
					Detail: response.Detail{
						Code:    "104",
						Message: "file size more than 128kb",
					},
				},
			},
		},
		{
			name: "Invalid Row Data Less Than 1",
			header: map[string]string{
				"Authorization": "Bearer token_test",
			},
			path:     "/excel-upload",
			filePath: "./storage/final_test_no_data.csv",
			typeMedia: test_util.MultipartRequest{
				Field: "type",
				Value: "file",
			},
			want: want{
				err:  true,
				code: http.StatusBadRequest,
				response: response.BaseResponse{
					RequestID: "abc123",
					Status:    "failed",
					Detail: response.Detail{
						Code:    "102",
						Message: "invalid row data",
					},
				},
			},
		},
		{
			name: "Invalid Row Data More Than 10",
			header: map[string]string{
				"Authorization": "Bearer token_test",
			},
			path:     "/excel-upload",
			filePath: "./storage/final_test_more_than_10.csv",
			typeMedia: test_util.MultipartRequest{
				Field: "type",
				Value: "file",
			},
			want: want{
				err:  true,
				code: http.StatusBadRequest,
				response: response.BaseResponse{
					RequestID: "abc123",
					Status:    "failed",
					Detail: response.Detail{
						Code:    "102",
						Message: "invalid row data",
					},
				},
			},
		},
		{
			name: "Invalid Row Mandatory Uniqid",
			header: map[string]string{
				"Authorization": "Bearer token_test",
			},
			path:     "/excel-upload",
			filePath: "./storage/final_test_empty_uniqid.csv",
			typeMedia: test_util.MultipartRequest{
				Field: "type",
				Value: "file",
			},
			want: want{
				err:  true,
				code: http.StatusBadRequest,
				response: response.BaseResponse{
					RequestID: "abc123",
					Status:    "failed",
					Detail: response.Detail{
						Code:    "103",
						Message: "empty row mandatory column uniqid",
					},
				},
			},
		},
		{
			name: "Invalid Row Mandatory Title",
			header: map[string]string{
				"Authorization": "Bearer token_test",
			},
			path:     "/excel-upload",
			filePath: "./storage/final_test_empty_title.csv",
			typeMedia: test_util.MultipartRequest{
				Field: "type",
				Value: "file",
			},
			want: want{
				err:  true,
				code: http.StatusBadRequest,
				response: response.BaseResponse{
					RequestID: "abc123",
					Status:    "failed",
					Detail: response.Detail{
						Code:    "103",
						Message: "empty row mandatory column title",
					},
				},
			},
		},
		{
			name: "Invalid Row Mandatory Description",
			header: map[string]string{
				"Authorization": "Bearer token_test",
			},
			path:     "/excel-upload",
			filePath: "./storage/final_test_empty_description.csv",
			typeMedia: test_util.MultipartRequest{
				Field: "type",
				Value: "file",
			},
			want: want{
				err:  true,
				code: http.StatusBadRequest,
				response: response.BaseResponse{
					RequestID: "abc123",
					Status:    "failed",
					Detail: response.Detail{
						Code:    "103",
						Message: "empty row mandatory column description",
					},
				},
			},
		},
		{
			name: "Invalid Row Mandatory Condition",
			header: map[string]string{
				"Authorization": "Bearer token_test",
			},
			path:     "/excel-upload",
			filePath: "./storage/final_test_empty_condition.csv",
			typeMedia: test_util.MultipartRequest{
				Field: "type",
				Value: "file",
			},
			want: want{
				err:  true,
				code: http.StatusBadRequest,
				response: response.BaseResponse{
					RequestID: "abc123",
					Status:    "failed",
					Detail: response.Detail{
						Code:    "103",
						Message: "empty row mandatory column condition",
					},
				},
			},
		},
		{
			name: "Invalid Row Mandatory Price",
			header: map[string]string{
				"Authorization": "Bearer token_test",
			},
			path:     "/excel-upload",
			filePath: "./storage/final_test_empty_price.csv",
			typeMedia: test_util.MultipartRequest{
				Field: "type",
				Value: "file",
			},
			want: want{
				err:  true,
				code: http.StatusBadRequest,
				response: response.BaseResponse{
					RequestID: "abc123",
					Status:    "failed",
					Detail: response.Detail{
						Code:    "103",
						Message: "empty row mandatory column price",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.name)
			t.Log(tt.typeMedia)
			res, statusCode, err := suite.requestSend(t, tt.header, tt.path, tt.filePath, tt.typeMedia)
			assert.NoError(t, err)
			if v, ok := res.(response.BaseResponse); ok {
				assert.Equal(t, tt.want.response.Status, v.Status)
				assert.Equal(t, tt.want.response.Detail, v.Detail)
				assert.Equal(t, tt.want.code, statusCode)
			} else {
				assert.Fail(t, fmt.Sprintf("%s: failed response is not expected", tt.name))
			}
			if !tt.want.err {
				t.Log("total expected queue: ", tt.want.expectedLengthQueue)
				suite.validateQueue(t, tt.want.expectedLengthQueue)
			}
		})
	}
}

func (s *testSvc) requestSend(t *testing.T, headers map[string]string, path string, filepath string, otherData test_util.MultipartRequest) (interface{}, int, error) {
	var b bytes.Buffer
	var w *multipart.Writer

	if filepath != "" {
		b, w = test_util.CreateMultipartFormData(t, "file", filepath, otherData)
	} else {
		data := test_util.JSONRequest{
			Key:  "file",
			Type: "file",
			Src:  "/root/cobacoba",
		}
		jsonData, _ := json.Marshal(data)
		b = *bytes.NewBuffer(jsonData)
	}

	url := fmt.Sprintf("http://localhost:5050%s", path)
	req, err := http.NewRequest(http.MethodPost, url, &b)
	assert.NoError(t, err)
	if filepath != "" {
		req.Header.Set("Content-Type", w.FormDataContentType())
	} else {
		req.Header.Set("Content-Type", "application/json")
	}
	for key, val := range headers {
		req.Header.Set(key, val)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	respBody := json.NewDecoder(resp.Body)

	var responseMap map[string]interface{}
	err = respBody.Decode(&responseMap)
	if err != nil {
		return nil, 0, err
	}

	fmt.Printf("RESP = %v\n", responseMap)

	respBytes, _ := json.Marshal(responseMap)

	var response response.BaseResponse
	err = json.Unmarshal(respBytes, &response)
	if err != nil {
		return nil, 0, err
	}

	return response, resp.StatusCode, nil
}

func (s *testSvc) validateQueue(t *testing.T, expectedLengthQueue int) {
	subs, err := s.artemis.Subscribe(config.Configuration.Artemis.Address, stomp.AckAuto)
	assert.NoError(t, err)
	assert.NotNil(t, subs)
	totalQueue := 0

	stop := make(chan struct{})
	// Mulai goroutine untuk membaca pesan dari saluran
	go func() {
		for {
			select {
			case m, ok := <-subs.C:
				if !ok {
					return
				}
				t.Log(string(m.Body))
				totalQueue++
			case <-stop:
				return
			}
		}
	}()

	<-time.After(time.Second * 1)
	close(stop)

	err = subs.Unsubscribe()
	assert.NoError(t, err)
}
