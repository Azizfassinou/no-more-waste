package main

import (
	"net/http"

	"github.com/FASSINOU/no-more-waste-api/internal/database"
	"github.com/FASSINOU/no-more-waste-api/internal/handlers"
	"github.com/FASSINOU/no-more-waste-api/internal/middlewares"
	"github.com/FASSINOU/no-more-waste-api/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()

	services.StartRenewalChecker()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API no More Waste en ligne avec SQLite Opérationnelle !",
		})
	})

	r.POST("/login", handlers.Login)
	r.POST("/register/client", handlers.RegisterClient)
	r.POST("/register/merchant", handlers.RegisterMerchant)
	r.GET("/Profile/:id", handlers.GetMyProfile)

	merchantGroup := r.Group("/merchant")
	merchantGroup.Use(middlewares.AuthMiddleware(), middlewares.MerchantOnly())
	{
		merchantGroup.POST("/product", handlers.CreateProduct)
		merchantGroup.PUT("/:id", handlers.UpdateMyMerchantProfile)
	}

	staffGroup := r.Group("/staff")
	staffGroup.Use(middlewares.AuthMiddleware(), middlewares.StaffOnly())
	{
		staffGroup.GET("/merchants/pending", handlers.GetPendingMerchants)
		staffGroup.PUT("/merchants/:id/approve", handlers.ApproveMerchant)
		staffGroup.POST("/volunteers", handlers.RegisterVolunteer)
		staffGroup.POST("/missions", handlers.CreateCollectionMission)
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
		adminGroup.POST("/missions", handlers.CreateCollectionMission)
	}
	r.Run(":8080")
}
