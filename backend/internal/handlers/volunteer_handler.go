package handlers

import (
	"net/http"
	"time"

	"github.com/FASSINOU/no-more-waste-api/internal/database"
	"github.com/FASSINOU/no-more-waste-api/internal/models"
	"github.com/gin-gonic/gin"
)

func getOrCreateVolunteerProfile(userID uint) (*models.Volunteer, error) {
	var volunteer models.Volunteer
	err := database.DB.Where("user_id = ?", userID).First(&volunteer).Error
	if err != nil {
		volunteer = models.Volunteer{
			UserID:       userID,
			ZoneArea:     "Toutes zones",
			Availability: "Flexible",
			Vehicle:      false,
		}
		if createErr := database.DB.Create(&volunteer).Error; createErr != nil {
			return nil, createErr
		}
	}
	return &volunteer, nil
}

func GetMyVolunteerRounds(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	volunteer, _ := getOrCreateVolunteerProfile(userID)
	if volunteer == nil {
		c.JSON(http.StatusOK, gin.H{"volunteer": nil, "rounds": []models.DistributionRound{}})
		return
	}

	var rounds []models.DistributionRound
	if err := database.DB.Preload("Deliveries").Where("volunteer_id = ?", volunteer.ID).Order("date desc").Find(&rounds).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des tournées"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"volunteer": volunteer,
		"rounds":    rounds,
	})
}

func GetMyVolunteerMissions(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	volunteer, _ := getOrCreateVolunteerProfile(userID)
	if volunteer == nil {
		c.JSON(http.StatusOK, gin.H{"missions": []models.CollectionMission{}})
		return
	}

	var missions []models.CollectionMission
	if err := database.DB.Preload("Merchant").Where("volunteer_id = ?", volunteer.ID).Order("pickup_date desc").Find(&missions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des missions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"missions": missions,
	})
}

func GetMyVolunteerServices(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	volunteer, _ := getOrCreateVolunteerProfile(userID)
	if volunteer == nil {
		c.JSON(http.StatusOK, gin.H{"services": []models.Service{}})
		return
	}

	var services []models.Service
	if err := database.DB.Preload("Registrations.Client").Where("volunteer_id = ?", volunteer.ID).Order("date desc").Find(&services).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération de vos services"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"services": services,
	})
}

func CreateVolunteerService(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	volunteer, err := getOrCreateVolunteerProfile(userID)
	if err != nil || volunteer == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Impossible de charger votre profil bénévole"})
		return
	}

	var input struct {
		Title           string    `json:"title" binding:"required"`
		Description     string    `json:"description" binding:"required"`
		Category        string    `json:"category" binding:"required"`
		MaxParticipants int       `json:"max_participants" binding:"required"`
		Date            time.Time `json:"date" binding:"required"`
		Location        string    `json:"location" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données incomplètes : " + err.Error()})
		return
	}

	service := models.Service{
		Title:           input.Title,
		Description:     input.Description,
		Category:        input.Category,
		MaxParticipants: input.MaxParticipants,
		Date:            input.Date,
		Location:        input.Location,
		Status:          "open",
		VolunteerID:     volunteer.ID,
	}

	if err := database.DB.Create(&service).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de la formation"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Formation / Atelier créé avec succès !",
		"data":    service,
	})
}

func UpdateVolunteerProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	volunteer, _ := getOrCreateVolunteerProfile(userID)

	var input struct {
		Firstname    string `json:"firstname"`
		Lastname     string `json:"lastname"`
		Phone        string `json:"phone"`
		Address      string `json:"address"`
		ZoneArea     string `json:"zone_area"`
		Availability string `json:"availability"`
		Vehicle      bool   `json:"vehicle"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Firstname != "" {
		user.FirstName = input.Firstname
	}
	if input.Lastname != "" {
		user.LastName = input.Lastname
	}
	if input.Phone != "" {
		user.Phone = input.Phone
	}
	if input.Address != "" {
		user.Address = input.Address
	}
	database.DB.Save(&user)

	if volunteer != nil {
		if input.ZoneArea != "" {
			volunteer.ZoneArea = input.ZoneArea
		}
		if input.Availability != "" {
			volunteer.Availability = input.Availability
		}
		volunteer.Vehicle = input.Vehicle
		database.DB.Save(volunteer)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Profil bénévole mis à jour avec succès !",
		"user":      user,
		"volunteer": volunteer,
	})
}

func UpdateVolunteerMissionStatus(c *gin.Context) {
	missionID := c.Param("id")
	userID := c.MustGet("user_id").(uint)

	volunteer, _ := getOrCreateVolunteerProfile(userID)
	if volunteer == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Accès refusé"})
		return
	}

	var input struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var mission models.CollectionMission
	if err := database.DB.Where("id = ? AND volunteer_id = ?", missionID, volunteer.ID).First(&mission).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mission non trouvée ou non assignée à ce bénévole"})
		return
	}

	mission.Status = input.Status
	if err := database.DB.Save(&mission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Statut de la mission mis à jour",
		"data":    mission,
	})
}

func UpdateVolunteerDeliveryStatus(c *gin.Context) {
	deliveryID := c.Param("delivery_id")

	var input struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var delivery models.Delivery
	if err := database.DB.First(&delivery, deliveryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Livraison non trouvée"})
		return
	}

	delivery.Status = input.Status
	if err := database.DB.Save(&delivery).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour de la livraison"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Statut de livraison mis à jour avec succès",
		"data":    delivery,
	})
}