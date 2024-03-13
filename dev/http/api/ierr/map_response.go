package ierr

import (
	"errors"
	"excel-receiver/constant"
	"excel-receiver/dto/response"

	"net/http"
)

func MapResponse(err error, requestId, status string) (int, response.BaseResponse) {
	if err == nil {
		return http.StatusAccepted, response.BaseResponse{
			// Code:     string(constant.APISuccess),
			// Message:  "success",
			// Products: data,
			RequestID: requestId,
			Status:    status,
			Detail: response.Detail{
				Code:    string(constant.APISuccess),
				Message: "success",
			},
		}
	}

	var iErr *Error
	if errors.As(err, &iErr) {
		return iErr.GetHTTPCode(), response.BaseResponse{
			// Code:    iErr.GetCode(),
			// Message: iErr.Error(),
			RequestID: requestId,
			Status:    status,
			Detail: response.Detail{
				Code:    iErr.GetCode(),
				Message: iErr.Error(),
			},
		}
	}

	return http.StatusInternalServerError, response.BaseResponse{
		// Code:    string(constant.APIInternalError),
		// Message: "internal server error",
		RequestID: requestId,
		Status:    status,
		Detail: response.Detail{
			Code:    string(constant.APIInternalError),
			Message: "internal server error",
		},
	}
}
