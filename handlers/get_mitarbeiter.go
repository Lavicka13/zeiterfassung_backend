package handlers

import (
    "zeiterfassung-backend/config"
    "zeiterfassung-backend/models"
    "github.com/gin-gonic/gin"
)

func GetMitarbeiter(c *gin.Context) {
    var mitarbeiter []models.Nutzer
    if err := config.DB.Find(&mitarbeiter).Error; err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, mitarbeiter)
}
