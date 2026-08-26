package main

import (
	"log"

	"github.com/FASSINOU/no-more-waste-api/internal/database"
	"github.com/FASSINOU/no-more-waste-api/internal/models"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("../../.env")

	log.Println("=== Démarrage du Seeder de Production — No More Waste ===")
	database.InitDB()
	db := database.DB

	defaultPassword := "Password123!"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Erreur lors du hachage du mot de passe : %v", err)
	}
	hashedStr := string(hashedPassword)

	adminUser := models.User{
		FirstName: "Admin",
		LastName:  "Système",
		Email:     "admin@nomorewaste.com",
		Password:  hashedStr,
		Role:      "admin",
		Phone:     "0102030405",
		Address:   "1 Place de la République, 75011 Paris",
		IsActive:  true,
	}
	if err := db.Where(models.User{Email: adminUser.Email}).FirstOrCreate(&adminUser).Error; err != nil {
		log.Printf("Erreur création Admin : %v", err)
	} else {
		log.Printf("[OK] Compte Admin : %s (Mot de passe : %s)", adminUser.Email, defaultPassword)
	}

	staffUser := models.User{
		FirstName: "Marc",
		LastName:  "Dubois",
		Email:     "staff@nomorewaste.com",
		Password:  hashedStr,
		Role:      "staff",
		Phone:     "0612345678",
		Address:   "10 Rue de la Paix, 75002 Paris",
		IsActive:  true,
	}
	if err := db.Where(models.User{Email: staffUser.Email}).FirstOrCreate(&staffUser).Error; err != nil {
		log.Printf("Erreur création Staff User : %v", err)
	} else {
		log.Printf("[OK] Compte Staff : %s (Mot de passe : %s)", staffUser.Email, defaultPassword)
		staffProfile := models.Staff{
			UserID:     staffUser.ID,
			Department: "Logistique & Opérations",
			JobTitle:   "Responsable de Collecte",
			Salary:     3200.00,
		}
		db.Where(models.Staff{UserID: staffUser.ID}).FirstOrCreate(&staffProfile)
	}

	skillsList := []models.Skill{
		{Name: "Permis B / Transport", Category: "Transport"},
		{Name: "Distribution Alimentaire", Category: "Logistique"},
		{Name: "Sensibilisation Zéro-Déchet", Category: "Formation"},
		{Name: "Animation d'Atelier", Category: "Animation"},
		{Name: "Gestion des Stocks", Category: "Manutention"},
	}

	seededSkills := make([]models.Skill, 0)
	for _, sk := range skillsList {
		var s models.Skill
		if err := db.Where("name = ?", sk.Name).FirstOrCreate(&s, sk).Error; err == nil {
			seededSkills = append(seededSkills, s)
			log.Printf("[OK] Compétence : %s (%s)", s.Name, s.Category)
		}
	}

	volunteersData := []struct {
		FirstName    string
		LastName     string
		Email        string
		Phone        string
		Address      string
		ZoneArea     string
		Availability string
		Vehicle      bool
		SkillIndices []int
	}{
		{
			FirstName:    "Sophie",
			LastName:     "Martin",
			Email:        "benevole1@nomorewaste.com",
			Phone:        "0699887766",
			Address:      "25 Rue Oberkampf, 75011 Paris",
			ZoneArea:     "Paris-Nord",
			Availability: "Semaine (Matin)",
			Vehicle:      true,
			SkillIndices: []int{0, 2},
		},
		{
			FirstName:    "Lucas",
			LastName:     "Bernard",
			Email:        "benevole2@nomorewaste.com",
			Phone:        "0655443322",
			Address:      "88 Boulevard Saint-Germain, 75005 Paris",
			ZoneArea:     "Paris-Sud",
			Availability: "Week-end",
			Vehicle:      false,
			SkillIndices: []int{1, 3},
		},
	}

	for _, vData := range volunteersData {
		u := models.User{
			FirstName: vData.FirstName,
			LastName:  vData.LastName,
			Email:     vData.Email,
			Password:  hashedStr,
			Role:      "volunteer",
			Phone:     vData.Phone,
			Address:   vData.Address,
			IsActive:  true,
		}
		if err := db.Where(models.User{Email: u.Email}).FirstOrCreate(&u).Error; err == nil {
			log.Printf("[OK] Compte Bénévole User : %s (Mot de passe : %s)", u.Email, defaultPassword)
			vProfile := models.Volunteer{
				UserID:       u.ID,
				ZoneArea:     vData.ZoneArea,
				Availability: vData.Availability,
				Vehicle:      vData.Vehicle,
			}
			if err := db.Where(models.Volunteer{UserID: u.ID}).FirstOrCreate(&vProfile).Error; err == nil {
				assignedSkills := make([]models.Skill, 0)
				for _, idx := range vData.SkillIndices {
					if idx < len(seededSkills) {
						assignedSkills = append(assignedSkills, seededSkills[idx])
					}
				}
				db.Model(&vProfile).Association("Skills").Replace(assignedSkills)
			}
		}
	}

	clientsData := []struct {
		FirstName string
		LastName  string
		Email     string
		Phone     string
		Address   string
	}{
		{
			FirstName: "Claire",
			LastName:  "Petit",
			Email:     "client1@nomorewaste.com",
			Phone:     "0711223344",
			Address:   "14 Rue de Rivoli, 75004 Paris",
		},
		{
			FirstName: "Antoine",
			LastName:  "Rousseau",
			Email:     "client2@nomorewaste.com",
			Phone:     "0722334455",
			Address:   "5 Avenue de la République, 75011 Paris",
		},
		{
			FirstName: "Élodie",
			LastName:  "Moreau",
			Email:     "client3@nomorewaste.com",
			Phone:     "0733445566",
			Address:   "42 Rue de Clichy, 75009 Paris",
		},
	}

	for _, cData := range clientsData {
		cu := models.User{
			FirstName: cData.FirstName,
			LastName:  cData.LastName,
			Email:     cData.Email,
			Password:  hashedStr,
			Role:      "client",
			Phone:     cData.Phone,
			Address:   cData.Address,
			IsActive:  true,
		}
		if err := db.Where(models.User{Email: cu.Email}).FirstOrCreate(&cu).Error; err == nil {
			log.Printf("[OK] Compte Client : %s (Mot de passe : %s)", cu.Email, defaultPassword)
		}
	}

	log.Println(" Seeding Terminé avec Succès ")
}
