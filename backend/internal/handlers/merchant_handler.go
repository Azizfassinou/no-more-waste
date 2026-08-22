package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/FASSINOU/no-more-waste-api/internal/database"
	"github.com/FASSINOU/no-more-waste-api/internal/models"
	"github.com/FASSINOU/no-more-waste-api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v78"
	"golang.org/x/crypto/bcrypt"
)

func generateBarcode() string {
	now := time.Now()
	code := fmt.Sprintf("NMW-%s-%04d", now.Format("20060102150405"), rand.Intn(10000))
	return code
}

func RegisterMerchant(c *gin.Context) {
	var input struct {
		FirstName     string `json:"firstname" binding:"required"`
		LastName      string `json:"last_name" binding:"required"`
		Password      string `json:"password" binding:"required,min=8"`
		CompanyAdress string `json:"company_address" binding:"required"`
		CompanyName   string `json:"company_name" binding:"required"`
		SiretNumber   string `json:"siret_number" binding:"required"`
		Email         string `json:"email" binding:"required,email" `
		Phone         string `json:"phone" binding:"required"`
		Address       string `json:"address"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Veuillez renseigner toutes les informations de votre commerce correctement."})
		return
	}

	siretInfo, _ := services.ValidateSiret(input.SiretNumber)
	if siretInfo != nil && !siretInfo.IsValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SIRET Invalide : " + siretInfo.ErrorMessage})
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
		RenewalDate:    time.Now().AddDate(1, 0, 0),
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
	if err := tx.Where("user_id = ?", userID).First(&merchant).Error; err == nil {
		if err := tx.Delete(&merchant).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression du commerçant"})
			return
		}
	}

	var user models.User
	if err := tx.First(&user, userID).Error; err == nil {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Veuillez renseigner toutes les informations du produit correctement."})
		return
	}

	if input.OriginalPrice < 0 || input.DiscountPrice < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Le prix ne peut pas être inférieur à 0 €."})
		return
	}

	if input.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "La quantité doit être supérieure à 0."})
		return
	}

	if input.DiscountPrice > input.OriginalPrice {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Le prix réduit ne peut pas être supérieur au prix d'origine."})
		return
	}

	product := models.Product{
		Title:         input.Title,
		Description:   input.Description,
		Barcode:       generateBarcode(),
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

func SearchProductByBarcode(c *gin.Context) {
	barcode := c.Query("code")
	if barcode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Le paramètre 'code' est requis"})
		return
	}

	var product models.Product
	if err := database.DB.Preload("Merchant.User").Where("barcode = ?", barcode).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Aucun produit trouvé pour ce code-barre"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}

func SearchProducts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Le paramètre 'q' est requis"})
		return
	}

	var products []models.Product
	search := "%" + query + "%"
	if err := database.DB.Preload("Merchant.User").
		Where("title LIKE ? OR description LIKE ? OR barcode LIKE ?", search, search, search).
		Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la recherche"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": products})
}

func GetMerchantDashboardStats(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var merchant models.Merchant
	if err := database.DB.Where("user_id = ?", userID).First(&merchant).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profil commerçant non trouvé"})
		return
	}

	var products []models.Product
	database.DB.Where("merchant_id = ?", merchant.ID).Find(&products)

	var productIDs []uint
	totalRemainingUnits := 0
	for _, p := range products {
		productIDs = append(productIDs, p.ID)
		totalRemainingUnits += p.Quantity
	}

	var totalRevenue float64 = 0.0
	var totalSoldUnits int = 0
	var productStatsList []gin.H

	if len(productIDs) > 0 {
		var orders []models.Order
		database.DB.Where("product_id IN ? AND status = ?", productIDs, "paid").Find(&orders)

		productSales := make(map[uint]float64)
		productSoldQty := make(map[uint]int)

		for _, o := range orders {
			totalRevenue += o.TotalPrice
			productSales[o.ProductID] += o.TotalPrice
			productSoldQty[o.ProductID]++
			totalSoldUnits++
		}

		for _, p := range products {
			rev := productSales[p.ID]
			sold := productSoldQty[p.ID]
			productStatsList = append(productStatsList, gin.H{
				"id":             p.ID,
				"title":          p.Title,
				"description":    p.Description,
				"barcode":        p.Barcode,
				"original_price": p.OriginalPrice,
				"discount_price": p.DiscountPrice,
				"quantity":       p.Quantity,
				"expiry_date":    p.ExpiryDate,
				"is_available":   p.IsAvailable,
				"sold_units":     sold,
				"revenue":        rev,
			})
		}
	} else {
		productStatsList = []gin.H{}
	}

	c.JSON(http.StatusOK, gin.H{
		"cagnotte_revenue":         totalRevenue,
		"total_published_products": len(products),
		"total_sold_units":         totalSoldUnits,
		"total_remaining_stock":    totalRemainingUnits,
		"renewal_date":             merchant.RenewalDate,
		"is_approved":              merchant.IsApproved,
		"company_name":             merchant.CompanyName,
		"company_address":          merchant.CompanyAddress,
		"siret_number":             merchant.SiretNumber,
		"merchant_id":              merchant.ID,
		"products":                 productStatsList,
	})
}

func RenewMerchantSubscription(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var merchant models.Merchant
	if err := database.DB.Where("user_id = ?", userID).First(&merchant).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profil commerçant non trouvé"})
		return
	}

	var newRenewalDate time.Time
	if time.Now().After(merchant.RenewalDate) {
		newRenewalDate = time.Now().AddDate(1, 0, 0)
	} else {
		newRenewalDate = merchant.RenewalDate.AddDate(1, 0, 0)
	}

	merchant.RenewalDate = newRenewalDate
	if err := database.DB.Save(&merchant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour de l'abonnement"})
		return
	}

	lineItems := []*stripe.CheckoutSessionLineItemParams{
		{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String("eur"),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String("Abonnement Annuel Commerçant No-More-Waste"),
				},
				UnitAmount: stripe.Int64(1500),
			},
			Quantity: stripe.Int64(1),
		},
	}

	checkoutURL, err := services.CreateStripeCheckoutSession(lineItems, "/merchant-dashboard?payment=success", "/merchant-dashboard?payment=cancelled")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Erreur API Stripe : " + err.Error() + ". Veuillez vérifier votre clé STRIPE_SECRET_KEY dans le fichier backend/.env",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Redirection vers le paiement Stripe...",
		"renewal_date": merchant.RenewalDate,
		"url":          checkoutURL,
	})
}