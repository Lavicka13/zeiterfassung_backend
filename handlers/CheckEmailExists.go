package handlers

import (
	"net/http"
	"zeiterfassung-backend/config"
	"zeiterfassung-backend/models"
	
	"github.com/gin-gonic/gin"
)

// CheckEmailExists prüft, ob ein Nutzer mit der angegebenen E-Mail existiert
func CheckEmailExists(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Keine E-Mail angegeben"})
		return
	}

	var nutzer models.Nutzer
	result := config.DB.Where("email = ?", email).First(&nutzer)
	
	// Prüfen, ob ein Nutzer gefunden wurde
	exists := result.Error == nil
	
	c.JSON(http.StatusOK, gin.H{"exists": exists})
}