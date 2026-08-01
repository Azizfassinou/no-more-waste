package handlers

import (
	"net/http"
	"time"

	"github.com/FASSINOU/no-more-waste-api/internal/database"
	"github.com/FASSINOU/no-more-waste-api/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetPendingMerchants(c *gin.Context) {
	var merchants []models.Merchant

	if err := database.DB.Preload("User").Where("is_approved = ?", false).Find(&merchants).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des commerçants en attente de validation"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": merchants,
	})
}

func ApproveMerchant(c *gin.Context) {
	id := c.Param("id")
	var merchant models.Merchant
	if err := database.DB.First(&merchant, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Commerçant non trouvé"})
		return
	}
	merchant.IsApproved = true
	if err := database.DB.Save(&merchant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'approbation du commerçant"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Commerçant approuvé avec succès",
		"data":    merchant,
	})
}

func RegisterVolunteer(c *gin.Context) {
	var input struct {
		FirstName    string `json:"firstname" binding:"required"`
		LastName     string `json:"lastname" binding:"required"`
		Email        string `json:"email" binding:"required, email"`
		Password     string `json:"password" binding:"required,min=8"`
		ZoneArea     string `json:"zone_area" binding:"required"`
		Availability string `json:"availability" binding:"required"`
		Vehicle      bool   `json:"vehicle" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	user := models.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  string(hashedPassword),
		Role:      "volunteer",
		IsActive:  true,
	}

	tx := database.DB.Begin()
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de l'utilisateur"})
		return
	}

	volunteer := models.Volunteer{
		UserID:       user.ID,
		ZoneArea:     input.ZoneArea,
		Availability: input.Availability,
		Vehicle:      input.Vehicle,
	}

	if err := tx.Create(&volunteer).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du bénévole"})
		return
	}
	tx.Commit()
	c.JSON(http.StatusCreated, gin.H{
		"message": "Votre compte a été créé avec succès",
	})

}

func CreateCollectionMission(c *gin.Context) {
	var input struct {
		MerchantID  uint      `json:"merchant_id" binding:"required"`
		VolunteerID uint      `json:"volunteer_id" binding:"required"`
		PickupDate  time.Time `json:"pickup_date" binding:"required"`
		ProductID   uint      `json:"product_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}
