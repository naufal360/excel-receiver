package api

import (
	"errors"
	"excel-receiver/constant"
	"excel-receiver/dto/response"
	"excel-receiver/ierr"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapResponse(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		requestID      string
		status         string
		expectedStatus int
		expectedResp   response.BaseResponse
	}{
		{
			name:           "No error",
			err:            nil,
			requestID:      "abc123",
			status:         "success",
			expectedStatus: http.StatusAccepted,
			expectedResp: response.BaseResponse{
				RequestID: "abc123",
				Status:    "success",
				Detail: response.Detail{
					Code:    string(constant.APISuccess),
					Message: "success",
				},
			},
		},
		{
			name:           "Error is ierr.Error",
			err:            ierr.NewF("101", "invalid request"),
			requestID:      "def456",
			status:         "failed",
			expectedStatus: http.StatusBadRequest,
			expectedResp: response.BaseResponse{
				RequestID: "def456",
				Status:    "failed",
				Detail: response.Detail{
					Code:    "101",
					Message: "invalid request",
				},
			},
		},
		{
			name:           "Error is generic error",
			err:            errors.New("some error"),
			requestID:      "ghi789",
			status:         "failed",
			expectedStatus: http.StatusInternalServerError,
			expectedResp: response.BaseResponse{
				RequestID: "ghi789",
				Status:    "failed",
				Detail: response.Detail{
					Code:    string(constant.APIInternalError),
					Message: "internal server error",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, resp := mapResponse(tt.err, tt.requestID, tt.status)
			assert.Equal(t, tt.expectedStatus, status)
			assert.Equal(t, tt.expectedResp, resp)
		})
	}
}
