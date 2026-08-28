package handlers

import (
	"net/http"

	"github.com/FASSINOU/no-more-waste-api/internal/database"
	"github.com/FASSINOU/no-more-waste-api/internal/models"
	"github.com/gin-gonic/gin"
)

func CreateSkill(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Category string `json:"category" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	skill := models.Skill{
		Name:     input.Name,
		Category: input.Category,
	}
	if err := database.DB.Create(&skill).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de la compétence"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Compétence créée avec succès",
		"data":    skill,
	})
}

func GetAllSkills(c *gin.Context) {
	var skills []models.Skill
	if err := database.DB.Find(&skills).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des compétences"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": skills})
}

func DeleteSkill(c *gin.Context) {
	id := c.Param("id")
	var skill models.Skill
	if err := database.DB.First(&skill, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Compétence non trouvée"})
		return
	}

	database.DB.Model(&skill).Association("Volunteers").Clear()

	if err := database.DB.Delete(&skill).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression de la compétence"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Compétence supprimée avec succès"})
}

func AssignSkillsToVolunteer(c *gin.Context) {
	volunteerID := c.Param("id")

	var volunteer models.Volunteer
	if err := database.DB.First(&volunteer, volunteerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bénévole non trouvé"})
		return
	}

	var input struct {
		SkillIDs []uint `json:"skill_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var skills []models.Skill
	if err := database.DB.Where("id IN ?", input.SkillIDs).Find(&skills).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des compétences"})
		return
	}
	if len(skills) != len(input.SkillIDs) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Certaines compétences n'existent pas"})
		return
	}

	if err := database.DB.Model(&volunteer).Association("Skills").Replace(skills); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'assignation des compétences"})
		return
	}

	if err := database.DB.Preload("Skills").Preload("User").First(&volunteer, volunteerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bénévole non trouvé"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Compétences assignées avec succès",
		"data":    volunteer,
	})
}

func GetVolunteerWithSkills(c *gin.Context) {
	volunteerID := c.Param("id")

	var volunteer models.Volunteer
	if err := database.DB.Preload("User").Preload("Skills").First(&volunteer, volunteerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bénévole non trouvé"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": volunteer})
}

func GetAllVolunteers(c *gin.Context) {
	var volunteers []models.Volunteer
	if err := database.DB.Preload("User").Preload("Skills").Find(&volunteers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des bénévoles"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": volunteers})
}

func SearchVolunteersBySkill(c *gin.Context) {
	skillName := c.Query("skill")
	if skillName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Le paramètre 'skill' est requis"})
		return
	}

	var skill models.Skill
	if err := database.DB.Where("name = ?", skillName).First(&skill).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Compétence non trouvée"})
		return
	}

	var volunteers []models.Volunteer
	if err := database.DB.Model(&skill).Preload("User").Preload("Skills").Association("Volunteers").Find(&volunteers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la recherche"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": volunteers})
}
