package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/FASSINOU/no-more-waste-api/internal/database"
	"github.com/FASSINOU/no-more-waste-api/internal/models"
	"github.com/FASSINOU/no-more-waste-api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Veuillez renseigner une adresse e-mail et un mot de passe valides."})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou mot de passe incorrect"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou mot de passe incorrect"})
		return
	}
	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Compte désactivé"})
		return
	}

	secretKey := []byte(os.Getenv("JWT_SECRET"))
	if len(secretKey) == 0 {
		secretKey = []byte("no_more_waste_default_secret_key_2026")
	}
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la génération du token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

func ForgotPassword(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Veuillez fournir une adresse email valide"})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "Si l'adresse existe, un code de réinitialisation vous a été envoyé.",
		})
		return
	}

	code := "123456"
	expires := time.Now().Add(15 * time.Minute)

	user.ResetCode = code
	user.ResetExpiresAt = &expires
	database.DB.Save(&user)

	services.SendResetPasswordCode(user.Email, code)

	c.JSON(http.StatusOK, gin.H{
		"message": "Un code de réinitialisation vous a été transmis par email.",
		"code":    code,
	})
}

func ResetPassword(c *gin.Context) {
	var input struct {
		Email       string `json:"email" binding:"required,email"`
		Code        string `json:"code" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides (minimum 6 caractères pour le mot de passe)"})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	if user.ResetCode == "" || user.ResetCode != input.Code {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code de réinitialisation incorrect"})
		return
	}

	if user.ResetExpiresAt != nil && time.Now().After(*user.ResetExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Le code de réinitialisation a expiré"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	user.ResetCode = ""
	user.ResetExpiresAt = nil
	database.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"message": "Mot de passe réinitialisé avec succès ! Vous pouvez maintenant vous connecter.",
	})
}