package database

import (
	"log"

	"github.com/FASSINOU/no-more-waste-api/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("no-more-waste.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Echec de la connexion à la base de données : %v", err)
	}
	err = DB.AutoMigrate(&models.User{}, &models.Merchant{}, &models.Volunteer{}, &models.Staff{}, &models.Order{}, &models.CollectionMission{})
	if err != nil {
		log.Fatalf("Echec de la migration de la base de données : %v", err)
	}
	log.Println("Base de données SQLite initialisée avec succès !")
}
