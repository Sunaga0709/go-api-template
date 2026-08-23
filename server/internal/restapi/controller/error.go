package controller

import (
	"errors"
	"net/http"

	"github.com/Sunaga0709/go-api-template/internal/restapi/middleware"
)

func requestBodyRequiredError() middleware.ErrorDetail {
	return middleware.NewErrorDetail(http.StatusBadRequest, errors.New("request body is required"))
}

func unknownError(err error) middleware.ErrorDetail {
	return middleware.NewErrorDetail(http.StatusInternalServerError, err)
}
