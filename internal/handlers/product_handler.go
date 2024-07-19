package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProducts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get products"})
}

func CreateProduct(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Create Product"})
}
