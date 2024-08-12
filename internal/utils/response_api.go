package utils

import (
	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResponseWithMeta struct {
	ResponseData
	Meta interface{} `json:"meta"`
}

// RespondJSON sends a JSON response with a status code, message, and payload.
//
// Parameters:
// - w: gin.Context to write the response
// - status: int representing the HTTP status code
// - message: string containing the message to be included in the response
// - payload: interface{} representing the data to be sent in the response
func RespondJSON(w *gin.Context, status int, message string, payload interface{}) {
	var res ResponseData
	res.Status = status
	res.Message = message
	res.Data = payload

	w.JSON(status, res)
}

// ResponseJSONWithMeta sends a JSON response with a meta field.
//
// Parameters:
// - w: gin.Context to write the response
// - status: int representing the HTTP status code
// - message: string containing the message to be included in the response
// - payload: interface{} representing the data to be sent in the response
// - meta: interface{} containing metadata about the response
func ResponseJSONWithMeta(w *gin.Context, status int, message string, payload interface{}, meta interface{}) {
	var res ResponseWithMeta
	res.Status = status
	res.Message = message
	res.Data = payload
	res.Meta = meta

	w.JSON(status, res)
}
