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
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Commerçant créé avec succès",
		"data":    input,
	})
}

func UpdateMerchant(c *gin.Context) {
	id := c.Param("id")
	var merchant models.Merchant
	if err := database.DB.Where("id = ?", id).First(&merchant).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Commerçant non trouvé"})
		return
	}

	var input models.Merchant
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	merchant.Name = input.Name
	merchant.Email = input.Email
	merchant.Phone = input.Phone
	merchant.Address = input.Address
	merchant.RenewalDate = input.RenewalDate
	merchant.IsActive = input.IsActive

	if err := database.DB.Save(&merchant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Commerçant mis à jour avec succès",
		"data":    merchant,
	})
}

func DeleteMerchant(c *gin.Context) {
	id := c.Param("id")

	var merchant models.Merchant
	if err := database.DB.First(&merchant, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	database.DB.Delete(&merchant)
	c.JSON(http.StatusOK, gin.H{
		"message": "Commerçant supprimé avec succès",
	})
}
