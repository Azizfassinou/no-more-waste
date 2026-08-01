package handlers

import (
	"net/http"

	"github.com/FASSINOU/no-more-waste-api/internal/database"
	"github.com/FASSINOU/no-more-waste-api/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CreateStaffProfile(c *gin.Context) {
	var input struct {
		FirstName  string  `json:"firstname" binding:"required"`
		LastName   string  `json:"lastname" binding:"required"`
		Email      string  `json:"email" binding:"required, email"`
		Password   string  `json:"password" binding:"required,min=8"`
		Department string  `json:"department"`
		JobTitle   string  `json:"job_title"`
		Salary     float64 `json:"salary"`
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
		Role:      "staff",
		IsActive:  true,
	}

	tx := database.DB.Begin()
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du profil du personnel"})
		return
	}

	staff := models.Staff{
		Department: input.Department,
		JobTitle:   input.JobTitle,
		Salary:     input.Salary,
		UserID:     user.ID,
	}

	if err := tx.Create(&staff).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du profil du personnel"})
		return
	}
	tx.Commit()
	c.JSON(http.StatusCreated, gin.H{
		"message": "Profil du personnel créé avec succès",
		"data":    staff,
	})
}

func GetUsers(c *gin.Context) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des utilisateurs"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}

func UpdateUserRole(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := database.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur introuvable"})
		return
	}
	var input struct {
		Role     string `json:"role"`
		IsActive bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.Role != "" {
		user.Role = input.Role
	}
	if input.IsActive {
		user.IsActive = input.IsActive
	}
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour du rôle de l'utilisateur"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Rôle de l'utilisateur mis à jour avec succès",
		"data":    user,
	})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := database.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur introuvable"})
		return
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression de l'utilisateur"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Utilisateur supprimé avec succès",
	})
}
