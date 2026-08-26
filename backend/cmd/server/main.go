package main

import (
	"log"
	"net/http"

	"github.com/FASSINOU/no-more-waste-api/internal/database"
	"github.com/FASSINOU/no-more-waste-api/internal/handlers"
	"github.com/FASSINOU/no-more-waste-api/internal/middlewares"
	"github.com/FASSINOU/no-more-waste-api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		err = godotenv.Load("../.env")
	}
	if err != nil {
		err = godotenv.Load("../../.env")
	}

	if err == nil {
		log.Println("Fichier .env chargé avec succès !")
	}

	database.InitDB()

	services.StartRenewalChecker()

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API no More Waste en ligne avec SQLite Opérationnelle !",
		})
	})
	r.POST("/seed", handlers.SeedDatabase)
	r.GET("/seed", handlers.SeedDatabase)

	r.POST("/login", handlers.Login)
	r.POST("/forgot-password", handlers.ForgotPassword)
	r.POST("/reset-password", handlers.ResetPassword)
	r.POST("/register/client", handlers.RegisterClient)
	r.POST("/register/merchant", handlers.RegisterMerchant)
	r.POST("/register/volunteer", handlers.RegisterVolunteer)
	r.GET("/profile", middlewares.AuthMiddleware(), handlers.GetMyProfile)
	r.PUT("/profile", middlewares.AuthMiddleware(), handlers.UpdateProfile)

	clientGroup := r.Group("/client")
	clientGroup.Use(middlewares.AuthMiddleware())
	{
		clientGroup.GET("/products", handlers.GetAvailableProductsForClients)
		clientGroup.GET("/products/scan", handlers.SearchProductByBarcode)
		clientGroup.GET("/products/search", handlers.SearchProducts)
		clientGroup.POST("/orders", handlers.CreateOrder)
		clientGroup.GET("/my-orders", handlers.GetMyOrders)
		clientGroup.GET("/services", handlers.GetAvailableServices)
		clientGroup.GET("/services/:id", handlers.GetServiceByID)
		clientGroup.POST("/services/:id/register", handlers.RegisterToService)
		clientGroup.DELETE("/services/:id/register", handlers.UnregisterFromService)
		clientGroup.GET("/my-services", handlers.GetMyServiceRegistrations)
	}

	merchantGroup := r.Group("/merchant")
	merchantGroup.Use(middlewares.AuthMiddleware(), middlewares.MerchantOnly())
	{
		merchantGroup.GET("/stats", handlers.GetMerchantDashboardStats)
		merchantGroup.POST("/product", handlers.CreateProduct)
		merchantGroup.POST("/subscription/renew", handlers.RenewMerchantSubscription)
		merchantGroup.PUT("/:id", handlers.UpdateMyMerchantProfile)
		merchantGroup.DELETE("/profile", handlers.DeleteMyMerchantProfile)
	}

	volunteerGroup := r.Group("/volunteer")
	volunteerGroup.Use(middlewares.AuthMiddleware(), middlewares.VolunteerOnly())
	{
		volunteerGroup.GET("/rounds", handlers.GetMyVolunteerRounds)
		volunteerGroup.GET("/missions", handlers.GetMyVolunteerMissions)
		volunteerGroup.PUT("/missions/:id/status", handlers.UpdateVolunteerMissionStatus)
		volunteerGroup.PUT("/deliveries/:delivery_id", handlers.UpdateVolunteerDeliveryStatus)
		volunteerGroup.GET("/rounds/:id/pdf", handlers.GenerateDistributionPDF)
		volunteerGroup.GET("/distribution-rounds/:id/pdf", handlers.GenerateDistributionPDF)
		volunteerGroup.GET("/services", handlers.GetMyVolunteerServices)
		volunteerGroup.POST("/services", handlers.CreateVolunteerService)
		volunteerGroup.PUT("/profile", handlers.UpdateVolunteerProfile)
	}

	staffGroup := r.Group("/staff")
	staffGroup.Use(middlewares.AuthMiddleware(), middlewares.StaffOnly())
	{
		staffGroup.GET("/merchants/pending", handlers.GetPendingMerchants)
		staffGroup.PUT("/merchants/:id/approve", handlers.ApproveMerchant)
		staffGroup.POST("/volunteers", handlers.RegisterVolunteer)
		staffGroup.GET("/volunteers", handlers.GetAllVolunteers)
		staffGroup.GET("/volunteers/:id", handlers.GetVolunteerWithSkills)
		staffGroup.PUT("/volunteers/:id/skills", handlers.AssignSkillsToVolunteer)
		staffGroup.GET("/volunteers/search", handlers.SearchVolunteersBySkill)
		staffGroup.POST("/missions", handlers.CreateCollectionMission)
		staffGroup.POST("/skills", handlers.CreateSkill)
		staffGroup.GET("/skills", handlers.GetAllSkills)
		staffGroup.DELETE("/skills/:id", handlers.DeleteSkill)
		staffGroup.POST("/services", handlers.CreateService)
		staffGroup.GET("/services", handlers.GetAllServices)
		staffGroup.GET("/services/:id", handlers.GetServiceByID)
		staffGroup.PUT("/services/:id", handlers.UpdateService)
		staffGroup.DELETE("/services/:id", handlers.DeleteService)
		staffGroup.GET("/services/export", handlers.ExportServicesCSV)
		staffGroup.POST("/distribution-rounds", handlers.CreateDistributionRound)
		staffGroup.GET("/distribution-rounds", handlers.GetAllDistributionRounds)
		staffGroup.GET("/distribution-rounds/:id", handlers.GetDistributionRoundByID)
		staffGroup.PUT("/distribution-rounds/:id", handlers.UpdateDistributionRoundStatus)
		staffGroup.DELETE("/distribution-rounds/:id", handlers.DeleteDistributionRound)
		staffGroup.POST("/distribution-rounds/:id/deliveries", handlers.AddDeliveryToRound)
		staffGroup.PUT("/deliveries/:delivery_id", handlers.UpdateDeliveryStatus)
		staffGroup.GET("/deliveries", handlers.GetAllDeliveries)
		staffGroup.PUT("/deliveries/:delivery_id/assign", handlers.AssignDeliveryToRound)
		staffGroup.GET("/distribution-rounds/:id/pdf", handlers.GenerateDistributionPDF)
	}

	adminGroup := r.Group("/admin")
	adminGroup.Use(middlewares.AuthMiddleware(), middlewares.AdminOnly())
	{
		adminGroup.GET("/users", handlers.GetUsers)
		adminGroup.PUT("/users/:id", handlers.UpdateUserRole)
		adminGroup.DELETE("/users/:id", handlers.DeleteUser)
		adminGroup.POST("/staff", handlers.CreateStaffProfile)
		adminGroup.GET("/merchants/pending", handlers.GetPendingMerchants)
		adminGroup.PUT("/merchants/:id/approve", handlers.ApproveMerchant)
		adminGroup.POST("/volunteers", handlers.RegisterVolunteer)
		adminGroup.GET("/volunteers", handlers.GetAllVolunteers)
		adminGroup.GET("/volunteers/:id", handlers.GetVolunteerWithSkills)
		adminGroup.PUT("/volunteers/:id/skills", handlers.AssignSkillsToVolunteer)
		adminGroup.GET("/volunteers/search", handlers.SearchVolunteersBySkill)
		adminGroup.POST("/missions", handlers.CreateCollectionMission)
		adminGroup.POST("/skills", handlers.CreateSkill)
		adminGroup.GET("/skills", handlers.GetAllSkills)
		adminGroup.DELETE("/skills/:id", handlers.DeleteSkill)
		adminGroup.POST("/services", handlers.CreateService)
		adminGroup.GET("/services", handlers.GetAllServices)
		adminGroup.GET("/services/:id", handlers.GetServiceByID)
		adminGroup.PUT("/services/:id", handlers.UpdateService)
		adminGroup.DELETE("/services/:id", handlers.DeleteService)
		adminGroup.GET("/services/export", handlers.ExportServicesCSV)
		adminGroup.POST("/distribution-rounds", handlers.CreateDistributionRound)
		adminGroup.GET("/distribution-rounds", handlers.GetAllDistributionRounds)
		adminGroup.GET("/distribution-rounds/:id", handlers.GetDistributionRoundByID)
		adminGroup.PUT("/distribution-rounds/:id", handlers.UpdateDistributionRoundStatus)
		adminGroup.DELETE("/distribution-rounds/:id", handlers.DeleteDistributionRound)
		adminGroup.POST("/distribution-rounds/:id/deliveries", handlers.AddDeliveryToRound)
		adminGroup.PUT("/deliveries/:delivery_id", handlers.UpdateDeliveryStatus)
		adminGroup.GET("/deliveries", handlers.GetAllDeliveries)
		adminGroup.PUT("/deliveries/:delivery_id/assign", handlers.AssignDeliveryToRound)
		adminGroup.GET("/distribution-rounds/:id/pdf", handlers.GenerateDistributionPDF)
	}
	r.Run(":8080")
}
