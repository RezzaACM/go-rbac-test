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

func RespondJSON(w *gin.Context, status int, message string, payload interface{}) {
	var res ResponseData
	res.Status = status
	res.Message = message
	res.Data = payload

	w.JSON(status, res)
}

func ResponseJSONWithMeta(w *gin.Context, status int, message string, payload interface{}, meta interface{}) {
	var res ResponseWithMeta
	res.Status = status
	res.Message = message
	res.Data = payload
	res.Meta = meta

	w.JSON(status, res)
}
