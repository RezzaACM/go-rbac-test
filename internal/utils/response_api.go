package utils

import (
	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Status  int
	Message string
	Data    interface{}
}

func RespondJSON(w *gin.Context, status int, message string, payload interface{}) {
	var res ResponseData
	res.Status = status
	res.Message = message
	res.Data = payload

	w.JSON(status, res)
}
