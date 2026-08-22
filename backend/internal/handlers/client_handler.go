package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/FASSINOU/no-more-waste-api/internal/database"
	"github.com/FASSINOU/no-more-waste-api/internal/models"
	"github.com/FASSINOU/no-more-waste-api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v78"
	"golang.org/x/crypto/bcrypt"
)

func RegisterClient(c *gin.Context) {
	var input struct {
		FirstName string `json:"firstname" binding:"required"`
		LastName  string `json:"lastname" binding:"required"`
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required,min=8"`
		Address   string `json:"address" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Veuillez renseigner tous les champs obligatoires correctement."})
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
	minExpiry := time.Now().Add(24 * time.Hour)

	if err := database.DB.Preload("Merchant").Where("is_available = ? AND quantity > ? AND expiry_date > ?", true, 0, minExpiry).Find(&products).Error; err != nil {
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
		ProductID uint `json:"product_id"`
		Quantity  int  `json:"quantity"`
		Items     []struct {
			ProductID uint `json:"product_id"`
			Quantity  int  `json:"quantity"`
		} `json:"items"`
		DeliveryAddress string `json:"delivery_address"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format du panier invalide"})
		return
	}

	if len(input.Items) == 0 && input.ProductID > 0 {
		q := input.Quantity
		if q <= 0 {
			q = 1
		}
		input.Items = append(input.Items, struct {
			ProductID uint `json:"product_id"`
			Quantity  int  `json:"quantity"`
		}{ProductID: input.ProductID, Quantity: q})
	}

	if len(input.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Le panier est vide"})
		return
	}

	tx := database.DB.Begin()
	var totalOrderPrice float64 = 0
	type itemDetail struct {
		Product  models.Product
		Quantity int
	}
	var details []itemDetail

	for _, item := range input.Items {
		var product models.Product
		if err := tx.Where("id = ?", item.ProductID).First(&product).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Produit non trouvé"})
			return
		}

		qty := item.Quantity
		if qty <= 0 {
			qty = 1
		}

		if !product.IsAvailable || product.Quantity < qty {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Stock insuffisant pour le produit : " + product.Title})
			return
		}

		if product.ExpiryDate.Before(time.Now().Add(24 * time.Hour)) {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Le produit '" + product.Title + "' expire dans moins de 24h et ne peut plus être acheté."})
			return
		}

		totalOrderPrice += product.DiscountPrice * float64(qty)
		details = append(details, itemDetail{Product: product, Quantity: qty})
	}

	if totalOrderPrice < 10.0 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Le montant minimum d'achat pour valider la commande est de 10,00 € (Montant actuel du panier : %.2f €).", totalOrderPrice),
		})
		return
	}

	var clientUser models.User
	tx.First(&clientUser, clientID.(uint))
	recipientAddr := input.DeliveryAddress
	if recipientAddr == "" {
		recipientAddr = clientUser.Address
	}

	var lastOrderID uint = 0
	for _, det := range details {
		det.Product.Quantity -= det.Quantity
		if det.Product.Quantity <= 0 {
			det.Product.IsAvailable = false
		}
		if err := tx.Save(&det.Product).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour du stock"})
			return
		}

		order := models.Order{
			ClientID:   clientID.(uint),
			ProductID:  det.Product.ID,
			TotalPrice: det.Product.DiscountPrice * float64(det.Quantity),
			Status:     "paid",
		}
		if err := tx.Create(&order).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de la commande"})
			return
		}
		lastOrderID = order.ID

		delivery := models.Delivery{
			RecipientName:    clientUser.FirstName + " " + clientUser.LastName,
			RecipientType:    "client",
			RecipientAddress: recipientAddr,
			ProductID:        det.Product.ID,
			Quantity:         det.Quantity,
			Status:           "pending",
		}
		tx.Create(&delivery)
	}

	tx.Commit()

	if lastOrderID > 0 {
		services.SendOrderConfirmationEmail(clientUser.Email, lastOrderID, totalOrderPrice)
	}

	var lineItems []*stripe.CheckoutSessionLineItemParams
	for _, det := range details {
		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String("eur"),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String(det.Product.Title),
				},
				UnitAmount: stripe.Int64(int64(det.Product.DiscountPrice * 100)),
			},
			Quantity: stripe.Int64(int64(det.Quantity)),
		})
	}

	checkoutURL, err := services.CreateStripeCheckoutSession(lineItems, "/client-dashboard?payment=success", "/client-dashboard?payment=cancelled")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Erreur API Stripe : " + err.Error() + ". Veuillez vérifier votre clé STRIPE_SECRET_KEY dans le fichier backend/.env",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Redirection vers le paiement Stripe...",
		"url":     checkoutURL,
	})
}

func GetMyProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	var merchant models.Merchant
	if user.Role == "merchant" {
		database.DB.Where("user_id = ?", user.ID).First(&merchant)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     user,
		"merchant": merchant,
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

func GetMyOrders(c *gin.Context) {
	clientID := c.MustGet("user_id").(uint)

	var orders []models.Order
	if err := database.DB.Preload("Product.Merchant").Where("client_id = ?", clientID).Order("created_at desc").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des commandes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": orders,
	})
}