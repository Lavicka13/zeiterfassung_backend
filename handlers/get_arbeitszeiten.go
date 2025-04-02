package handlers

import (
    "zeiterfassung-backend/config"
    "zeiterfassung-backend/models"
    "github.com/gin-gonic/gin"
)

func GetArbeitszeiten(c *gin.Context) {
    id := c.Param("id")
    monat := c.Query("monat") // optionaler Monatsfilter z.B. ?monat=2025-04

    var zeiten []models.Arbeitszeit
    query := config.DB.Where("nutzer_id = ?", id)
    if monat != "" {
        query = query.Where("DATE_FORMAT(datum, '%Y-%m') = ?", monat)
    }

    if err := query.Find(&zeiten).Error; err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, zeiten)
}
