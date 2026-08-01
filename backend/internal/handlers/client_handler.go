package handlers

import (
	"net/http"
	"time"

	"github.com/FASSINOU/no-more-waste-api/internal/database"
	"github.com/FASSINOU/no-more-waste-api/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterClient(c *gin.Context) {
	var input struct {
		FirstName string `json:"firstname" binding:"required"`
		LastName  string `json:"lastname" binding:"required"`
		Email     string `json:"email" binding:"required, email"`
		Password  string `json:"password" binding:"required,min=8"`
		Address   string `json:"address" binding:"required"`
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
		Address:   input.Address,
		Password:  string(hashedPassword),
		Role:      "client",
		IsActive:  true,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de l'utilisateur"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Votre compte a été créé avec succès",
	})
}

func GetAvailableProductsForClients(c *gin.Context) {
	var products []models.Product
	now := time.Now()

	if err := database.DB.Preload("Merchant").Where("is_available = ? AND quantity > ? AND expiry_date > ?", true, 0, now).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des produits"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": products,
	})
}

func CreateOrder(c *gin.Context) {
	clientID, _ := c.Get("user_id")
	var input struct {
		ProductID uint `json:"product_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var product models.Product
	if err := database.DB.Where("id = ?", input.ProductID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produit non trouvé"})
		return
	}

	if !product.IsAvailable || product.Quantity <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produit non disponible"})
		return
	}
	tx := database.DB.Begin()

	product.Quantity--
	if product.Quantity == 0 {
		product.IsAvailable = false
	}
	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour du produit"})
		return
	}

	order := models.Order{
		ClientID:   clientID.(uint),
		ProductID:  product.ID,
		TotalPrice: product.DiscountPrice,
		Status:     "pending",
	}
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de la commande"})
		return
	}
	tx.Commit()
	c.JSON(http.StatusCreated, gin.H{"message": "Panier en attente de réservation"})
}

func GetMyProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}
func UpdateProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	var input struct {
		FirstName string `json:"firstname"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
		Address   string `json:"address"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.FirstName != "" {
		user.FirstName = input.FirstName
	}
	if input.LastName != "" {
		user.LastName = input.LastName
	}
	if input.Email != "" {
		user.Email = input.Email
	}
	if input.Phone != "" {
		user.Phone = input.Phone
	}
	if input.Address != "" {
		user.Address = input.Address
	}

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour du profil"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profil mis à jour avec succès",
	})
}
