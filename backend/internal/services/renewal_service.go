package services

import (
	"log"
	"time"

	"github.com/FASSINOU/no-more-waste-api/internal/database"
	"github.com/FASSINOU/no-more-waste-api/internal/models"
)

func StartRenewalChecker() {
	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		for range ticker.C {
			checkRenewals()
		}
	}()
}

func checkRenewals() {
	var merchants []models.Merchant

	result := database.DB.Where("is_active = ?", true).Find(&merchants)
	if result.Error != nil {
		log.Printf("Erreur lors de la vérification des renouvellements : %v", result.Error)
		return
	}
	now := time.Now()
	sevenDaysBefore := now.Add(7 * 24 * time.Hour)

	for _, merchant := range merchants {
		if merchant.RenewalDate.Before(sevenDaysBefore) && merchant.RenewalDate.After(now) {
			log.Printf("Envoi d'un rappel à %s %s : Votre abonnement expire bientôt %s", merchant.FirstName, merchant.LastName, merchant.RenewalDate.Format("02/01/2006"))

			// Service d'e-mail à implémenter pour notifier l'expiration dans quelques jours
		}
		if merchant.RenewalDate.Before(now) {
			log.Printf("Le commerçant %s %s a un abonnement expiré", merchant.FirstName, merchant.LastName)
			merchant.IsActive = false
			database.DB.Save(&merchant)
			// Service d'e-mail à implémenter pour notifier l'expiration et la désactivation du compte
		}
	}
}
