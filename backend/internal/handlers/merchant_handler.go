package handlers

import (
	"net/http"

	"github.com/FASSINOU/no-more-waste-api/internal/database"
	"github.com/FASSINOU/no-more-waste-api/internal/models"
	"github.com/gin-gonic/gin"
)

func GetMerchants(c *gin.Context) {
	var merchants []models.Merchant

	result := database.DB.Find(&merchants)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, merchants)
}

func CreateMerchant(c *gin.Context) {
	var input models.Merchant
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := database.DB.Create(&input)
	if result.Error != nil {
		println("ERREUR SQL GORM: ", result.Error.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Marchant créé avec succès",
		"data":    input,
	})
}
