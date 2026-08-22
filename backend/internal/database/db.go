package database

import (
	"log"
	"os"

	"github.com/FASSINOU/no-more-waste-api/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dbFile := os.Getenv("DB_FILE")
	if dbFile == "" {
		dbFile = "no-more-waste.db"
		if _, err := os.Stat(dbFile); os.IsNotExist(err) {
			if _, err2 := os.Stat("../no-more-waste.db"); err2 == nil {
				dbFile = "../no-more-waste.db"
			} else if _, err3 := os.Stat("../../no-more-waste.db"); err3 == nil {
				dbFile = "../../no-more-waste.db"
			}
		}
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		log.Fatalf("Echec de la connexion à la base de données : %v", err)
	}
	err = DB.AutoMigrate(&models.User{}, &models.Merchant{}, &models.Volunteer{}, &models.Staff{}, &models.Order{}, &models.CollectionMission{}, &models.Product{}, &models.Service{}, &models.ServiceRegistration{}, &models.Skill{}, &models.DistributionRound{}, &models.Delivery{})
	if err != nil {
		log.Fatalf("Echec de la migration de la base de données : %v", err)
	}
	log.Println("Base de données SQLite initialisée avec succès !")
}