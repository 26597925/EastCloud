package server

import (
	"encoding/json"
	"github.com/26597925/EastCloud/pkg/oauth2/server/errors"
	"net/http"
)

type Response struct {
	Writer http.ResponseWriter
}

func (res *Response) Redirect(uri string) error {
	res.Writer.Header().Set("Location", uri)
	res.Writer.WriteHeader(302)
	return nil
}

func (res *Response) ErrorData(err error) (map[string]interface{}, int) {
	data := make(map[string]interface{})
	data["errors"] = err.Error()

	var statusCode int
	if statusCode = errors.StatusCodes[err]; statusCode == 0 {
		data["error_code"] = statusCode
		statusCode = http.StatusInternalServerError
	}

	if v := errors.Descriptions[err]; v != "" {
		data["error_description"] = v
	}

	return data, statusCode
}

func (res *Response) OutputTokenError(err error) error {
	data, statusCode := res.ErrorData(err)
	return res.OutputToken(data, statusCode)
}

func (res *Response) OutputToken(data map[string]interface{}, statusCode ...int) error {
	res.Writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
	res.Writer.Header().Set("Cache-Control", "no-stores")
	res.Writer.Header().Set("Pragma", "no-cache")

	status := http.StatusOK
	if len(statusCode) > 0 && statusCode[0] > 0 {
		status = statusCode[0]
	}

	res.Writer.WriteHeader(status)
	return json.NewEncoder(res.Writer).Encode(data)
}