package handlers

import (
	"errors"
	"net/http"
	"time"
	"zeiterfassung-backend/config"
	"zeiterfassung-backend/models"
	"zeiterfassung-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateOrUpdateArbeitszeit(c *gin.Context) {
	var eingabe struct {
		NutzerID    uint       `json:"nutzer_id"`
		Anfangszeit *time.Time `json:"anfangszeit,omitempty"`
		Endzeit     *time.Time `json:"endzeit,omitempty"`
		Datum       string     `json:"datum"`
	}

	if err := c.ShouldBindJSON(&eingabe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ungültige Eingabe"})
		return
	}

	// Datum aus dem Eingabe-String parsen
	datumParsed, err := time.Parse("2006-01-02", eingabe.Datum)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ungültiges Datum, erwartetes Format: YYYY-MM-DD"})
		return
	}

	var arbeitszeit models.Arbeitszeit
	// Vergleiche das Datum, indem du es in das gleiche Format umwandelst
	result := config.DB.Where("nutzer_id = ? AND DATE(datum) = ?", eingabe.NutzerID, datumParsed.Format("2006-01-02")).First(&arbeitszeit)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			neueAZ := models.Arbeitszeit{
				NutzerID:   eingabe.NutzerID,
				Datum:      datumParsed,
				ErstelltAm: time.Now(),
			}
			if eingabe.Anfangszeit != nil {
				neueAZ.Anfangszeit = *eingabe.Anfangszeit
			}
			if eingabe.Endzeit != nil {
				neueAZ.Endzeit = eingabe.Endzeit
				neueAZ.Pause = utils.CalcPause(neueAZ.Anfangszeit, *eingabe.Endzeit)
			}
			if err := config.DB.Create(&neueAZ).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, neueAZ)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Falls ein Eintrag gefunden wurde, aktualisieren wir ihn
	if eingabe.Anfangszeit != nil {
		arbeitszeit.Anfangszeit = *eingabe.Anfangszeit
	}
	if eingabe.Endzeit != nil {
		arbeitszeit.Endzeit = eingabe.Endzeit
		arbeitszeit.Pause = utils.CalcPause(arbeitszeit.Anfangszeit, *eingabe.Endzeit)
	}
	jetzt := time.Now()
	arbeitszeit.Bearbeitet = &jetzt

	if err := config.DB.Save(&arbeitszeit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, arbeitszeit)
}
