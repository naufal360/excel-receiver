package api

import (
	"errors"
	"excel-receiver/constant"
	"excel-receiver/dto/response"
	"excel-receiver/ierr"

	"net/http"
)

func mapResponse(err error, requestId, status string) (int, response.BaseResponse) {
	if err == nil {
		return http.StatusAccepted, response.BaseResponse{
			RequestID: requestId,
			Status:    status,
			Detail: response.Detail{
				Code:    string(constant.APISuccess),
				Message: "success",
			},
		}
	}

	var iErr *ierr.Error
	if errors.As(err, &iErr) {
		return iErr.GetHTTPCode(), response.BaseResponse{
			RequestID: requestId,
			Status:    status,
			Detail: response.Detail{
				Code:    iErr.GetCode(),
				Message: iErr.Error(),
			},
		}
	}

	return http.StatusInternalServerError, response.BaseResponse{
		RequestID: requestId,
		Status:    status,
		Detail: response.Detail{
			Code:    string(constant.APIInternalError),
			Message: "internal server error",
		},
	}
}
