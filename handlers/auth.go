package handlers

import (
    "net/http"
    "zeiterfassung-backend/config"
    "zeiterfassung-backend/models"
    "zeiterfassung-backend/utils"

    "github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
    var input struct {
        Vorname   string `json:"vorname"`
        Nachname  string `json:"nachname"`
        Email     string `json:"email"`
        Passwort  string `json:"passwort"`
        RechteID  uint   `json:"rechte_id"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Ungültige Daten"})
        return
    }

    hash, _ := utils.HashPassword(input.Passwort)
    nutzer := models.Nutzer{
        Vorname:  input.Vorname,
        Nachname: input.Nachname,
        Email:    input.Email,
        PwHash:   hash,
        RechteID: input.RechteID,
    }

    if err := config.DB.Create(&nutzer).Error; err != nil {
        c.JSON(http.StatusConflict, gin.H{"error": "E-Mail schon vergeben"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Registrierung erfolgreich"})
}

func Login(c *gin.Context) {
    var input struct {
        Email    string `json:"email"`
        Passwort string `json:"passwort"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Ungültige Daten"})
        return
    }

    var nutzer models.Nutzer
    if err := config.DB.Where("email = ?", input.Email).First(&nutzer).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Benutzer nicht gefunden"})
        return
    }

    if !utils.CheckPasswordHash(input.Passwort, nutzer.PwHash) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Falsches Passwort"})
        return
    }

    token, _ := utils.GenerateJWT(nutzer.ID, nutzer.Email, nutzer.RechteID)

    c.JSON(http.StatusOK, gin.H{
        "token":     token,
        "nutzer_id": nutzer.ID,
        "vorname":   nutzer.Vorname,
        "nachname":  nutzer.Nachname,
        "rechte_id": nutzer.RechteID,
    })
}

func Me(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "nutzer_id": c.GetUint("nutzer_id"),
        "rechte_id": c.GetUint("rechte_id"),
    })
}
