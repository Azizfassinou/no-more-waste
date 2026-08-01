package handlers

import (
	"net/http"
	"time"

	"github.com/FASSINOU/no-more-waste-api/internal/database"
	"github.com/FASSINOU/no-more-waste-api/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterMerchant(c *gin.Context) {
	var input struct {
		FirstName     string `json:"firstname" binding:"required"`
		LastName      string `json:"last_name" binding:"required"`
		Password      string `json:"password" binding:"required,min=8"`
		CompanyAdress string `json:"company_address" binding:"required"`
		CompanyName   string `json:"company_name" binding:"required"`
		SiretNumber   string `json:"siret_number" binding:"required"`
		Email         string `json:"email" binding:"required, email" `
		Phone         string `json:"phone" binding:"required"`
		Address       string `json:"address"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	user := models.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Address:   input.Address,
		Email:     input.Email,
		Phone:     input.Phone,
		Password:  string(hashedPassword),
		Role:      "merchant",
		IsActive:  true,
	}
	tx := database.DB.Begin()

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de l'utilisateur"})
		return
	}

	merchant := models.Merchant{
		CompanyName:    input.CompanyName,
		SiretNumber:    input.SiretNumber,
		CompanyAddress: input.CompanyAdress,
		UserID:         user.ID,
		IsApproved:     false,
	}
	if err := tx.Create(&merchant).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du commerçant"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Votre demande de création de compte a été soumise avec succès",
	})
}

func UpdateMyMerchantProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var merchant models.Merchant
	if err := database.DB.Preload("User").Where("user_id = ?", userID).First(&merchant).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Commerçant non trouvé"})
		return
	}

	var input struct {
		FirstName     string `json:"firstname" binding:"required"`
		LastName      string `json:"last_name" binding:"required"`
		Phone         string `json:"phone" binding:"required"`
		Address       string `json:"address" binding:"required"`
		CompanyName   string `json:"company_name" binding:"required"`
		CompanyAdress string `json:"company_address" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.Phone != "" {
		merchant.User.Phone = input.Phone
	}
	if input.Address != "" {
		merchant.User.Address = input.Address
	}
	if input.CompanyAdress != "" {
		merchant.CompanyAddress = input.CompanyAdress
	}
	if input.CompanyName != "" {
		merchant.CompanyName = input.CompanyName
	}

	if err := database.DB.Save(&merchant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profil Commerçant mis à jour avec succès",
		"data":    merchant,
	})
}

func DeleteMyMerchantProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	tx := database.DB.Begin()

	var merchant models.Merchant
	if err := tx.Where("user_id = ?", userID).First(&merchant).Error; err != nil {
		if err := tx.Delete(&merchant).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression du commerçant"})
			return
		}
	}

	var user models.User
	if err := tx.First(&user, userID).Error; err != nil {
		if err := tx.Delete(&user).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression de l'utilisateur"})
			return
		}
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "Votre compte de commerçant a été supprimé avec succès",
	})
}

func CreateProduct(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var merchant models.Merchant
	if err := database.DB.Where("user_id = ?", userID).First(&merchant).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Profil de commerçant non trouvé"})
		return
	}

	if !merchant.IsApproved {
		c.JSON(http.StatusForbidden, gin.H{"error": "Votre profil de commerçant n'est pas approuvé"})
		return
	}

	var input struct {
		Title         string    `json:"title" binding:"required"`
		Description   string    `json:"description" binding:"required"`
		OriginalPrice float64   `json:"original_price" binding:"required"`
		DiscountPrice float64   `json:"discount_price" binding:"required"`
		Quantity      int       `json:"quantity" binding:"required"`
		ExpiryDate    time.Time `json:"expiry_date" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := models.Product{
		Title:         input.Title,
		Description:   input.Description,
		OriginalPrice: input.OriginalPrice,
		DiscountPrice: input.DiscountPrice,
		Quantity:      input.Quantity,
		ExpiryDate:    input.ExpiryDate,
		IsAvailable:   true,
		MerchantID:    merchant.ID,
	}
	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du Produit"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Panier d'invidus créé et publié avec succès",
		"data":    product,
	})
}
