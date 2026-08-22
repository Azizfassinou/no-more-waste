package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/FASSINOU/no-more-waste-api/internal/database"
	"github.com/FASSINOU/no-more-waste-api/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/go-pdf/fpdf"
)

func CreateDistributionRound(c *gin.Context) {
	var input struct {
		Date        time.Time `json:"date" binding:"required"`
		VolunteerID uint      `json:"volunteer_id" binding:"required"`
		Notes       string    `json:"notes"`
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

	round := models.DistributionRound{
		Date:        input.Date,
		VolunteerID: input.VolunteerID,
		Notes:       input.Notes,
		Status:      "planned",
	}
	if err := database.DB.Create(&round).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de la tournée"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Tournée de distribution créée avec succès",
		"data":    round,
	})
}

func GetAllDistributionRounds(c *gin.Context) {
	var rounds []models.DistributionRound
	if err := database.DB.Preload("Volunteer.User").Preload("Deliveries.Product").Find(&rounds).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des tournées"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rounds})
}

func GetDistributionRoundByID(c *gin.Context) {
	id := c.Param("id")
	var round models.DistributionRound
	if err := database.DB.Preload("Volunteer.User").Preload("Deliveries.Product").First(&round, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournée non trouvée"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": round})
}

func UpdateDistributionRoundStatus(c *gin.Context) {
	id := c.Param("id")
	var round models.DistributionRound
	if err := database.DB.First(&round, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournée non trouvée"})
		return
	}

	var input struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	round.Status = input.Status
	if err := database.DB.Save(&round).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Statut de la tournée mis à jour",
		"data":    round,
	})
}

func DeleteDistributionRound(c *gin.Context) {
	id := c.Param("id")
	var round models.DistributionRound
	if err := database.DB.First(&round, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournée non trouvée"})
		return
	}

	tx := database.DB.Begin()
	if err := tx.Where("distribution_round_id = ?", round.ID).Delete(&models.Delivery{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression des livraisons"})
		return
	}
	if err := tx.Delete(&round).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression de la tournée"})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Tournée supprimée avec succès"})
}

func AddDeliveryToRound(c *gin.Context) {
	roundID := c.Param("id")

	var round models.DistributionRound
	if err := database.DB.First(&round, roundID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournée non trouvée"})
		return
	}

	var input struct {
		RecipientName    string `json:"recipient_name" binding:"required"`
		RecipientAddress string `json:"recipient_address" binding:"required"`
		RecipientType    string `json:"recipient_type" binding:"required"`
		ProductID        uint   `json:"product_id" binding:"required"`
		Quantity         int    `json:"quantity" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	if err := database.DB.First(&product, input.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produit non trouvé"})
		return
	}

	delivery := models.Delivery{
		RecipientName:       input.RecipientName,
		RecipientAddress:    input.RecipientAddress,
		RecipientType:       input.RecipientType,
		ProductID:           input.ProductID,
		Quantity:            input.Quantity,
		Status:              "pending",
		DistributionRoundID: round.ID,
	}
	if err := database.DB.Create(&delivery).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'ajout de la livraison"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Livraison ajoutée à la tournée",
		"data":    delivery,
	})
}

func UpdateDeliveryStatus(c *gin.Context) {
	deliveryID := c.Param("delivery_id")

	var delivery models.Delivery
	if err := database.DB.First(&delivery, deliveryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Livraison non trouvée"})
		return
	}

	var input struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	delivery.Status = input.Status
	if err := database.DB.Save(&delivery).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Statut de la livraison mis à jour",
		"data":    delivery,
	})
}

func GenerateDistributionPDF(c *gin.Context) {
	id := c.Param("id")

	var round models.DistributionRound
	if err := database.DB.
		Preload("Volunteer.User").
		Preload("Deliveries.Product.Merchant.User").
		First(&round, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournée non trouvée"})
		return
	}

	pdf := fpdf.New("P", "mm", "A4", "")
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.SetAutoPageBreak(true, 15)
	pdf.AddPage()

	pdf.SetFont("Helvetica", "B", 20)
	pdf.CellFormat(190, 12, tr("NO MORE WASTE"), "", 1, "C", false, 0, "")
	pdf.SetFont("Helvetica", "", 10)
	pdf.CellFormat(190, 6, tr("Association humanitaire de lutte contre le gaspillage"), "", 1, "C", false, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Helvetica", "B", 14)
	pdf.CellFormat(190, 8, tr(fmt.Sprintf("Récapitulatif de tournée #%d", round.ID)), "", 1, "L", false, 0, "")
	pdf.Ln(4)

	pdf.SetFont("Helvetica", "", 11)
	pdf.CellFormat(190, 6, tr(fmt.Sprintf("Date : %s", round.Date.Format("02/01/2006"))), "", 1, "L", false, 0, "")
	pdf.CellFormat(190, 6, tr(fmt.Sprintf("Statut : %s", round.Status)), "", 1, "L", false, 0, "")
	pdf.CellFormat(190, 6, tr(fmt.Sprintf("Bénévole : %s %s", round.Volunteer.User.FirstName, round.Volunteer.User.LastName)), "", 1, "L", false, 0, "")
	if round.Notes != "" {
		pdf.CellFormat(190, 6, tr(fmt.Sprintf("Notes : %s", round.Notes)), "", 1, "L", false, 0, "")
	}
	pdf.Ln(8)

	pdf.SetFont("Helvetica", "B", 12)
	pdf.CellFormat(190, 8, tr(fmt.Sprintf("Livraisons (%d)", len(round.Deliveries))), "", 1, "L", false, 0, "")
	pdf.Ln(2)

	pdf.SetFillColor(60, 60, 60)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Helvetica", "B", 9)
	pdf.CellFormat(8, 8, "#", "1", 0, "C", true, 0, "")
	pdf.CellFormat(35, 8, tr("Destinataire"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(22, 8, tr("Type"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(52, 8, tr("Adresse"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(45, 8, tr("Produit"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(10, 8, tr("Qté"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(18, 8, tr("Statut"), "1", 1, "C", true, 0, "")

	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Helvetica", "", 8)

	for i, d := range round.Deliveries {
		pdf.CellFormat(8, 7, fmt.Sprintf("%d", i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(35, 7, tr(truncateText(d.RecipientName, 22)), "1", 0, "L", false, 0, "")
		pdf.CellFormat(22, 7, tr(truncateText(d.RecipientType, 14)), "1", 0, "C", false, 0, "")
		pdf.CellFormat(52, 7, tr(truncateText(d.RecipientAddress, 32)), "1", 0, "L", false, 0, "")
		pdf.CellFormat(45, 7, tr(truncateText(d.Product.Title, 28)), "1", 0, "L", false, 0, "")
		pdf.CellFormat(10, 7, fmt.Sprintf("%d", d.Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(18, 7, tr(truncateText(d.Status, 10)), "1", 1, "C", false, 0, "")
	}

	pdf.Ln(15)
	pdf.SetFont("Helvetica", "", 9)
	pdf.CellFormat(190, 6, tr(fmt.Sprintf("Document généré le %s", time.Now().Format("02/01/2006 à 15:04"))), "", 1, "R", false, 0, "")
	pdf.CellFormat(190, 6, tr("NO MORE WASTE - Tous droits réservés"), "", 1, "R", false, 0, "")

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=tournee_%d.pdf", round.ID))

	if err := pdf.Output(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la generation du PDF"})
		return
	}
}

func GetAllDeliveries(c *gin.Context) {
	var deliveries []models.Delivery
	if err := database.DB.Preload("Product").Preload("DistributionRound.Volunteer.User").Order("created_at desc").Find(&deliveries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des livraisons"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": deliveries})
}

func AssignDeliveryToRound(c *gin.Context) {
	deliveryID := c.Param("delivery_id")
	var input struct {
		DistributionRoundID uint `json:"distribution_round_id" binding:"required"`
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

	var round models.DistributionRound
	if err := database.DB.First(&round, input.DistributionRoundID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tournée non trouvée"})
		return
	}

	delivery.DistributionRoundID = round.ID
	if err := database.DB.Save(&delivery).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'assignation de la livraison"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Livraison assignée à la tournée avec succès",
		"data":    delivery,
	})
}

func truncateText(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) > maxLen {
		return string(runes[:maxLen-3]) + "..."
	}
	return s
}