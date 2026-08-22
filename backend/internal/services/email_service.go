package services

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strconv"
	"strings"
)

func SendEmail(to string, subject string, htmlBody string) error {
	host := strings.TrimSpace(os.Getenv("SMTP_HOST"))
	portStr := strings.TrimSpace(os.Getenv("SMTP_PORT"))
	username := strings.TrimSpace(os.Getenv("SMTP_USERNAME"))
	password := strings.TrimSpace(os.Getenv("SMTP_PASSWORD"))
	from := strings.TrimSpace(os.Getenv("SMTP_FROM"))

	if from == "" {
		from = "No More Waste <contact@nomorewaste.fr>"
	}

	if host == "" || username == "" || password == "" {
		err := errors.New("configuration SMTP absente ou incomplète (SMTP_HOST, SMTP_USERNAME ou SMTP_PASSWORD non renseignés dans backend/.env)")
		log.Printf(" [SMTP ERROR] Impossible d'envoyer l'e-mail à %s [Sujet: %s] : %v", to, subject, err)
		return err
	}

	port := 587
	if p, err := strconv.Atoi(portStr); err == nil && p > 0 {
		port = p
	}

	header := make(map[string]string)
	header["From"] = from
	header["To"] = to
	header["Subject"] = "=?UTF-8?B?" + encodeRFC2047(subject) + "?="
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + htmlBody

	auth := smtp.PlainAuth("", username, password, host)
	addr := fmt.Sprintf("%s:%d", host, port)

	if port == 465 {
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         host,
		}
		conn, err := tls.Dial("tcp", addr, tlsconfig)
		if err != nil {
			log.Printf(" [SMTP ERROR] Échec de la connexion TLS avec le serveur %s:%d : %v", host, port, err)
			return fmt.Errorf("erreur connexion TLS SMTP (%s:%d): %w", host, port, err)
		}
		client, err := smtp.NewClient(conn, host)
		if err != nil {
			log.Printf(" [SMTP ERROR] Échec de l'initialisation du client SMTP pour %s : %v", host, err)
			return fmt.Errorf("erreur client SMTP: %w", err)
		}
		if err = client.Auth(auth); err != nil {
			log.Printf(" [SMTP ERROR] Échec de l'authentification SMTP (vérifiez SMTP_USERNAME et SMTP_PASSWORD) : %v", err)
			return fmt.Errorf("échec authentification SMTP: %w", err)
		}
		if err = client.Mail(from); err != nil {
			log.Printf(" [SMTP ERROR] Rejet de l'adresse expéditeur '%s' par le serveur SMTP : %v", from, err)
			return fmt.Errorf("expéditeur SMTP rejeté: %w", err)
		}
		if err = client.Rcpt(to); err != nil {
			log.Printf(" [SMTP ERROR] Rejet du destinataire '%s' par le serveur SMTP : %v", to, err)
			return fmt.Errorf("destinataire SMTP rejeté: %w", err)
		}
		w, err := client.Data()
		if err != nil {
			log.Printf(" [SMTP ERROR] Erreur lors du transfert des données de l'e-mail : %v", err)
			return fmt.Errorf("erreur données SMTP: %w", err)
		}
		_, err = w.Write([]byte(message))
		if err != nil {
			log.Printf(" [SMTP ERROR] Erreur lors de l'écriture du contenu de l'e-mail : %v", err)
			return fmt.Errorf("erreur écriture message SMTP: %w", err)
		}
		w.Close()
		client.Quit()
	} else {
		err := smtp.SendMail(addr, auth, username, []string{to}, []byte(message))
		if err != nil {
			log.Printf(" [SMTP ERROR] Échec d'envoi de l'e-mail à %s via le serveur %s:%d : %v. Veuillez vérifier votre hôte, port et mot de passe d'application dans backend/.env", to, host, port, err)
			return fmt.Errorf("échec d'envoi SMTP (%s:%d): %w", host, port, err)
		}
	}

	log.Printf(" [SMTP SUCCESS] E-mail transmis avec succès à %s [Sujet: %s]", to, subject)
	return nil
}

func SendAsyncEmail(to string, subject string, htmlBody string) {
	go func() {
		_ = SendEmail(to, subject, htmlBody)
	}()
}

func SendResetPasswordCode(toEmail string, code string) {
	subject := "Code de réinitialisation - No More Waste"
	body := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; padding: 20px; background-color: #f8fafc; color: #1e293b;">
			<div style="max-width: 500px; margin: 0 auto; background: #ffffff; padding: 30px; border-radius: 12px; border: 1px solid #e2e8f0;">
				<h2 style="color: #2d6a4f; margin-top: 0;">No More Waste </h2>
				<p>Bonjour,</p>
				<p>Vous avez demandé la réinitialisation de votre mot de passe.</p>
				<p>Voici votre code de sécurité à 6 chiffres :</p>
				<div style="background: #f1f5f9; padding: 15px; text-align: center; border-radius: 8px; font-size: 24px; font-weight: bold; letter-spacing: 5px; color: #2d6a4f; margin: 20px 0;">
					%s
				</div>
				<p style="font-size: 13px; color: #64748b;">Ce code expire dans 15 minutes. Si vous n'avez pas demandé cette réinitialisation, veuillez ignorer ce message.</p>
			</div>
		</div>
	`, code)
	SendAsyncEmail(toEmail, subject, body)
}

func SendOrderConfirmationEmail(toEmail string, orderID uint, totalAmount float64) {
	subject := fmt.Sprintf("Confirmation de commande N° %d - No More Waste", orderID)
	body := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; padding: 20px; background-color: #f8fafc; color: #1e293b;">
			<div style="max-width: 550px; margin: 0 auto; background: #ffffff; padding: 30px; border-radius: 12px; border: 1px solid #e2e8f0;">
				<h2 style="color: #2d6a4f; margin-top: 0;">Commande Validée ! </h2>
				<p>Merci pour votre engagement anti-gaspillage !</p>
				<p>Votre commande <strong>N° %d</strong> a bien été enregistrée et réglée avec succès.</p>
				<table style="width: 100%%; margin: 20px 0; border-collapse: collapse;">
					<tr>
						<td style="padding: 10px; border-bottom: 1px solid #e2e8f0;"><strong>Montant Total Payé :</strong></td>
						<td style="padding: 10px; border-bottom: 1px solid #e2e8f0; text-align: right; color: #2d6a4f; font-weight: bold;">%.2f €</td>
					</tr>
				</table>
				<p>Nos équipes bénévoles préparent votre tournée de livraison.</p>
				<p style="font-size: 13px; color: #64748b; margin-top: 30px;">L'équipe No More Waste</p>
			</div>
		</div>
	`, orderID, totalAmount)
	SendAsyncEmail(toEmail, subject, body)
}

func SendMerchantApprovedEmail(toEmail string, companyName string) {
	subject := "Bienvenue ! Votre compte commerçant No More Waste est validé "
	body := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; padding: 20px; background-color: #f8fafc; color: #1e293b;">
			<div style="max-width: 550px; margin: 0 auto; background: #ffffff; padding: 30px; border-radius: 12px; border: 1px solid #e2e8f0;">
				<h2 style="color: #2d6a4f; margin-top: 0;">Félicitations %s !</h2>
				<p>Votre demande de compte commerçant sur <strong>No More Waste</strong> a été vérifiée et approuvée par notre équipe Staff.</p>
				<p>Vous pouvez dès maintenant vous connecter et publier vos paniers d'invendus alimentaires :</p>
				<div style="text-align: center; margin: 25px 0;">
					<a href="http://localhost:5173" style="background: #2d6a4f; color: #ffffff; padding: 12px 25px; text-decoration: none; border-radius: 8px; font-weight: bold; display: inline-block;">Accéder à mon espace commerçant</a>
				</div>
			</div>
		</div>
	`, companyName)
	SendAsyncEmail(toEmail, subject, body)
}

func stripHTML(input string) string {
	var result strings.Builder
	inTag := false
	for _, r := range input {
		if r == '<' {
			inTag = true
		} else if r == '>' {
			inTag = false
		} else if !inTag {
			result.WriteRune(r)
		}
	}
	return strings.TrimSpace(result.String())
}

func encodeRFC2047(s string) string {
	return s
}