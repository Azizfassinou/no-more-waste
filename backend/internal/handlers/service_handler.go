package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"time"

	"github.com/FASSINOU/no-more-waste-api/internal/database"
	"github.com/FASSINOU/no-more-waste-api/internal/models"
	"github.com/gin-gonic/gin"
)

func CreateService(c *gin.Context) {
	var input struct {
		Title           string    `json:"title" binding:"required"`
		Description     string    `json:"description" binding:"required"`
		Category        string    `json:"category" binding:"required"`
		MaxParticipants int       `json:"max_participants" binding:"required"`
		Date            time.Time `json:"date" binding:"required"`
		Location        string    `json:"location" binding:"required"`
		VolunteerID     uint      `json:"volunteer_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var volunteer models.Volunteer
	if err := database.DB.First(&volunteer, input.VolunteerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bénévole non trouvé"})
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
		VolunteerID:     input.VolunteerID,
	}
	if err := database.DB.Create(&service).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du service"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Service créé avec succès",
		"data":    service,
	})
}

func GetAllServices(c *gin.Context) {
	var services []models.Service
	if err := database.DB.Preload("Volunteer.User").Preload("Registrations").Find(&services).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des services"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": services})
}

func GetAvailableServices(c *gin.Context) {
	var services []models.Service
	now := time.Now()
	if err := database.DB.Preload("Volunteer.User").
		Where("status = ? AND date > ?", "open", now).
		Find(&services).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des services"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": services})
}

func GetServiceByID(c *gin.Context) {
	id := c.Param("id")
	var service models.Service
	if err := database.DB.Preload("Volunteer.User").Preload("Registrations.User").First(&service, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service non trouvé"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": service})
}

func UpdateService(c *gin.Context) {
	id := c.Param("id")
	var service models.Service
	if err := database.DB.First(&service, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service non trouvé"})
		return
	}

	var input struct {
		Title           string    `json:"title"`
		Description     string    `json:"description"`
		Category        string    `json:"category"`
		MaxParticipants int       `json:"max_participants"`
		Date            time.Time `json:"date"`
		Location        string    `json:"location"`
		Status          string    `json:"status"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Title != "" {
		service.Title = input.Title
	}
	if input.Description != "" {
		service.Description = input.Description
	}
	if input.Category != "" {
		service.Category = input.Category
	}
	if input.MaxParticipants > 0 {
		service.MaxParticipants = input.MaxParticipants
	}
	if !input.Date.IsZero() {
		service.Date = input.Date
	}
	if input.Location != "" {
		service.Location = input.Location
	}
	if input.Status != "" {
		service.Status = input.Status
	}

	if err := database.DB.Save(&service).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour du service"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Service mis à jour avec succès",
		"data":    service,
	})
}

func DeleteService(c *gin.Context) {
	id := c.Param("id")
	var service models.Service
	if err := database.DB.First(&service, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service non trouvé"})
		return
	}

	tx := database.DB.Begin()
	if err := tx.Where("service_id = ?", service.ID).Delete(&models.ServiceRegistration{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression des inscriptions"})
		return
	}
	if err := tx.Delete(&service).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression du service"})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Service supprimé avec succès"})
}

func RegisterToService(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	serviceID := c.Param("id")

	var service models.Service
	if err := database.DB.Preload("Registrations").First(&service, serviceID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service non trouvé"})
		return
	}

	if service.Status != "open" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ce service n'est plus ouvert aux inscriptions"})
		return
	}

	var existing models.ServiceRegistration
	if err := database.DB.Where("service_id = ? AND user_id = ? AND status = ?", service.ID, userID, "registered").First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Vous êtes déjà inscrit à ce service"})
		return
	}

	currentCount := 0
	for _, r := range service.Registrations {
		if r.Status == "registered" {
			currentCount++
		}
	}
	if currentCount >= service.MaxParticipants {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ce service est complet"})
		return
	}

	registration := models.ServiceRegistration{
		ServiceID: service.ID,
		UserID:    userID,
		Status:    "registered",
	}
	if err := database.DB.Create(&registration).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'inscription"})
		return
	}

	if currentCount+1 >= service.MaxParticipants {
		service.Status = "full"
		database.DB.Save(&service)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Inscription au service réussie",
		"data":    registration,
	})
}

func UnregisterFromService(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	serviceID := c.Param("id")

	var registration models.ServiceRegistration
	if err := database.DB.Where("service_id = ? AND user_id = ? AND status = ?", serviceID, userID, "registered").First(&registration).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inscription non trouvée"})
		return
	}

	registration.Status = "cancelled"
	if err := database.DB.Save(&registration).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la désinscription"})
		return
	}

	var service models.Service
	if err := database.DB.First(&service, serviceID).Error; err == nil {
		if service.Status == "full" {
			service.Status = "open"
			database.DB.Save(&service)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Désinscription réussie"})
}

func GetMyServiceRegistrations(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var registrations []models.ServiceRegistration
	if err := database.DB.Preload("Service.Volunteer.User").
		Where("user_id = ? AND status = ?", userID, "registered").
		Find(&registrations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des inscriptions"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": registrations})
}

func ExportServicesCSV(c *gin.Context) {
	var services []models.Service
	if err := database.DB.Preload("Volunteer.User").Preload("Registrations.User").Find(&services).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des données"})
		return
	}

	c.Writer.Header().Set("Content-Type", "text/csv")
	c.Writer.Header().Set("Content-Disposition", "attachment;filename=services_export.csv")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	headers := []string{"ID Service", "Titre", "Date", "Catégorie", "Bénévole (Animateur)", "Statut Service", "Participant Prénom", "Participant Nom", "Participant Email", "Statut Inscription"}
	if err := writer.Write(headers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la génération du CSV"})
		return
	}

	for _, service := range services {
		dateStr := service.Date.Format("02/01/2006 15:04")
		volunteerName := ""
		if service.Volunteer.User.ID != 0 {
			volunteerName = service.Volunteer.User.FirstName + " " + service.Volunteer.User.LastName
		}

		if len(service.Registrations) == 0 {
			row := []string{
				fmt.Sprintf("%d", service.ID),
				service.Title,
				dateStr,
				service.Category,
				volunteerName,
				service.Status,
				"Aucun", "Aucun", "Aucun", "-",
			}
			writer.Write(row)
		} else {
			for _, reg := range service.Registrations {
				row := []string{
					fmt.Sprintf("%d", service.ID),
					service.Title,
					dateStr,
					service.Category,
					volunteerName,
					service.Status,
					reg.User.FirstName,
					reg.User.LastName,
					reg.User.Email,
					reg.Status,
				}
				writer.Write(row)
			}
		}
	}
}