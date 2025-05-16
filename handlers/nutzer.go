package handlers

import (
	"net/http"
	"strconv"
	"zeiterfassung-backend/config"
	"zeiterfassung-backend/models"
	"zeiterfassung-backend/utils"

	"github.com/gin-gonic/gin"
)

// CreateNutzer erstellt einen neuen Nutzer
func CreateNutzer(c *gin.Context) {
	var input struct {
		Vorname  string `json:"Vorname" binding:"required"`
		Nachname string `json:"Nachname" binding:"required"`
		Email    string `json:"Email" binding:"required,email"`
		Passwort string `json:"Passwort" binding:"required,min=6"`
		Rolle    string `json:"Rolle" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ermittle die Rechte-ID basierend auf dem Rollennamen
	var rechteID uint = 1 // Standard: Mitarbeiter
	if input.Rolle == "admin" {
		rechteID = 3
	} else if input.Rolle == "vorgesetzter" {
		rechteID = 2
	}

	// Hash das Passwort
	hash, err := utils.HashPassword(input.Passwort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Passwort-Fehler"})
		return
	}

	user := models.Nutzer{
		Vorname:  input.Vorname,
		Nachname: input.Nachname,
		Email:    input.Email,
		PwHash:   hash,
		RechteID: rechteID,
	}

	// Speichere den neuen Nutzer
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nutzer konnte nicht erstellt werden, Email möglicherweise schon vergeben"})
		return
	}

	// Gib den neuen Nutzer zurück, ohne PwHash
	user.PwHash = ""
	c.JSON(http.StatusCreated, user)
}

// UpdateNutzer aktualisiert die Daten eines Nutzers
func UpdateNutzer(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ungültige ID"})
		return
	}

	var input struct {
		Vorname  string `json:"Vorname"`
		Nachname string `json:"Nachname"`
		Email    string `json:"Email"`
		Rolle    string `json:"Rolle"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.Nutzer
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nutzer nicht gefunden"})
		return
	}

	// Aktualisiere nur die angegebenen Felder
	updates := make(map[string]interface{})
	if input.Vorname != "" {
		updates["vorname"] = input.Vorname
	}
	if input.Nachname != "" {
		updates["nachname"] = input.Nachname
	}
	if input.Email != "" {
		updates["email"] = input.Email
	}
	if input.Rolle != "" {
		var rechteID uint = 1
		if input.Rolle == "admin" {
			rechteID = 3
		} else if input.Rolle == "vorgesetzter" {
			rechteID = 2
		}
		updates["rechte_id"] = rechteID
	}

	if err := config.DB.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Aktualisierung fehlgeschlagen"})
		return
	}

	// Lade den aktualisierten Nutzer
	config.DB.First(&user, userID)

	// Gib den aktualisierten Nutzer zurück, ohne PwHash
	user.PwHash = ""
	c.JSON(http.StatusOK, user)
}

// DeleteNutzer löscht einen Nutzer und alle verknüpften Daten
func DeleteNutzer(c *gin.Context) {
	id := c.Param("id")
	
	// Prüfe, ob der Nutzer existiert
	var user models.Nutzer
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nutzer nicht gefunden"})
		return
	}
	
	// Beginne eine Transaktion
	tx := config.DB.Begin()
	
	// Lösche zuerst alle Audit-Einträge für die Arbeitszeiten des Nutzers
	if err := tx.Exec("DELETE FROM arbeitszeiten_audit WHERE arbeitszeit_id IN (SELECT id FROM arbeitszeiten WHERE nutzer_id = ?)", id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Löschen der Audit-Einträge"})
		return
	}
	
	// Lösche alle Arbeitszeiten des Nutzers
	if err := tx.Exec("DELETE FROM arbeitszeiten WHERE nutzer_id = ?", id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Löschen der Arbeitszeiten"})
		return
	}
	
	// Lösche dann den Nutzer
	if err := tx.Delete(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Löschen fehlgeschlagen"})
		return
	}
	
	// Schließe die Transaktion erfolgreich ab
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaktion konnte nicht abgeschlossen werden"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Nutzer erfolgreich gelöscht"})
}

// ResetPassword setzt das Passwort eines Nutzers zurück
func ResetPassword(c *gin.Context) {
	id := c.Param("id")
	
	var input struct {
		Passwort string `json:"passwort" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ungültiges Passwort"})
		return
	}

	// Prüfe, ob der Nutzer existiert
	var user models.Nutzer
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nutzer nicht gefunden"})
		return
	}

	// Hash das neue Passwort
	hash, err := utils.HashPassword(input.Passwort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Passwort-Hashing fehlgeschlagen"})
		return
	}

	// Aktualisiere das Passwort
	if err := config.DB.Model(&user).Update("pw_hash", hash).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Passwort-Update fehlgeschlagen"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Passwort erfolgreich zurückgesetzt"})
}