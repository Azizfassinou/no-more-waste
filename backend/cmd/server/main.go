package main

import (
	"net/http"

	"github.com/FASSINOU/no-more-waste-api/internal/database"
	"github.com/FASSINOU/no-more-waste-api/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()
	r := gin.Default()
	/* http.HandleFunc(("/ping"), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "API no More Waste en ligne"}`))
	})
	fmt.Println("Serveur démarré sur :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Erreur au démarrage du serveur : %v\n", err) */

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API no More Waste en ligne avec SQLite Opérationnelle !",
		})
	})

	r.GET("/merchants", handlers.GetMerchants)
	r.POST("/merchants", handlers.CreateMerchant)
	r.PUT("/merchants/:id", handlers.UpdateMerchant)
	r.DELETE("/merchants/:id", handlers.DeleteMerchant)

	r.Run(":8080")
}
