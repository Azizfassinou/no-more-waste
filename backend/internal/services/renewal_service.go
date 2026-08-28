package services

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
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

func sendEmail(to string, subject string, body string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USERNAME")
	smtpPass := os.Getenv("SMTP_PASSWORD")

	if smtpHost == "" || smtpPort == "" {
		log.Printf("[SIMULATION EMAIL] À: %s | Sujet: %s | Corps: %s", to, subject, body)
		return nil
	}
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body))

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{to}, msg)
	return err
}

func checkRenewals() {
	var merchants []models.Merchant

	result := database.DB.Preload("User").Where("is_active = ? AND is_approved = ?", true, true).Find(&merchants)
	if result.Error != nil {
		log.Printf("Erreur lors de la vérification des renouvellements : %v", result.Error)
		return
	}
	now := time.Now()
	sevenDaysBefore := now.Add(7 * 24 * time.Hour)

	for _, merchant := range merchants {
		if merchant.RenewalDate.Before(sevenDaysBefore) && merchant.RenewalDate.After(now) {
			subject := "No More Waste - Votre adhésion expire bientôt"
			body := fmt.Sprintf("Bonjour %s,\n\nVotre adhésion annuelle à No More Waste expire le %s.\nPensez à la renouveler pour continuer à profiter de nos services.\n\nL'équipe No More Waste", merchant.User.FirstName, merchant.RenewalDate.Format("02/01/2006"))

			err := sendEmail(merchant.User.Email, subject, body)
			if err != nil {
				log.Printf("Erreur d'envoi d'email de rappel à %s: %v", merchant.User.Email, err)
			} else {
				log.Printf("Rappel envoyé à %s", merchant.User.Email)
			}
		}
		if merchant.RenewalDate.Before(now) {
			log.Printf("Le commerçant %s %s a un abonnement expiré", merchant.User.FirstName, merchant.User.LastName)
			merchant.IsActive = false
			database.DB.Save(&merchant)

			subject := "No More Waste - Adhésion expirée"
			body := fmt.Sprintf("Bonjour %s,\n\nVotre adhésion à No More Waste est arrivée à expiration.\nVotre compte a été temporairement désactivé.\n\nContactez-nous pour la renouveler.\n\nL'équipe No More Waste", merchant.User.FirstName)

			err := sendEmail(merchant.User.Email, subject, body)
			if err != nil {
				log.Printf("Erreur d'envoi d'email de désactivation à %s: %v", merchant.User.Email, err)
			} else {
				log.Printf("Email de désactivation envoyé à %s", merchant.User.Email)
			}
		}
	}
}
