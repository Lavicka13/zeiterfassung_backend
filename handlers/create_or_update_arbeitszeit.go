// handlers/create_or_update_arbeitszeit.go - Aktualisierte Version
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

// --------------------------
// Create
// --------------------------
func CreateArbeitszeit(c *gin.Context) {
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

    datumParsed, err := time.Parse("2006-01-02", eingabe.Datum)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Ungültiges Datum, Format: YYYY-MM-DD"})
        return
    }

    // Prüfen, ob das Datum nicht älter als 3 Monate ist
    if !utils.IsEditAllowed(datumParsed) {
        c.JSON(http.StatusForbidden, gin.H{
            "error": "Einträge können nur bis zu 3 Monate rückwirkend erstellt werden.",
        })
        return
    }

    var arbeitszeit models.Arbeitszeit
    result := config.DB.Where("nutzer_id = ? AND DATE(datum) = ?", eingabe.NutzerID, datumParsed.Format("2006-01-02")).First(&arbeitszeit)
    if result.Error == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Es existiert bereits ein Eintrag für diesen Tag"})
        return
    }

    if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

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
}

// --------------------------
// Update + Protokoll + 3-Monats-Limit
// --------------------------
func UpdateArbeitszeit(c *gin.Context) {
    var eingabe struct {
        ID          uint       `json:"id"`
        Anfangszeit *time.Time `json:"anfangszeit"`
        Endzeit     *time.Time `json:"endzeit"`
        Bearbeiter  uint       `json:"bearbeiter_id"`
    }

    if err := c.ShouldBindJSON(&eingabe); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Ungültige Eingabe"})
        return
    }

    var arbeitszeit models.Arbeitszeit
    if err := config.DB.First(&arbeitszeit, eingabe.ID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Eintrag nicht gefunden"})
        return
    }

    // Prüfen, ob der Eintrag noch bearbeitet werden darf (3-Monats-Limit)
    if !utils.IsEditAllowed(arbeitszeit.Datum) {
        c.JSON(http.StatusForbidden, gin.H{
            "error": "Dieser Eintrag ist älter als 3 Monate und kann nicht mehr bearbeitet werden.",
        })
        return
    }

    jetzt := time.Now()
    changes := []string{}

    if eingabe.Anfangszeit != nil && !arbeitszeit.Anfangszeit.Equal(*eingabe.Anfangszeit) {
        changes = append(changes, "Anfangszeit von "+arbeitszeit.Anfangszeit.Format("15:04")+" auf "+eingabe.Anfangszeit.Format("15:04"))
        arbeitszeit.Anfangszeit = *eingabe.Anfangszeit
    }

    if eingabe.Endzeit != nil && (arbeitszeit.Endzeit == nil || !arbeitszeit.Endzeit.Equal(*eingabe.Endzeit)) {
        old := "-"
        if arbeitszeit.Endzeit != nil { old = arbeitszeit.Endzeit.Format("15:04") }
        changes = append(changes, "Endzeit von "+old+" auf "+eingabe.Endzeit.Format("15:04"))
        arbeitszeit.Endzeit = eingabe.Endzeit
    }

    arbeitszeit.Bearbeitet = &jetzt

    if arbeitszeit.Endzeit != nil {
        arbeitszeit.Pause = utils.CalcPause(arbeitszeit.Anfangszeit, *arbeitszeit.Endzeit)
    }

    if err := config.DB.Save(&arbeitszeit).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Änderung protokollieren
    for _, cng := range changes {
        config.DB.Exec(`INSERT INTO arbeitszeiten_audit (arbeitszeit_id, nutzer_id, feld, alter_wert, neuer_wert, bearbeitet_am) VALUES (?, ?, ?, ?, ?, ?)`,
            arbeitszeit.ID, eingabe.Bearbeiter, "manuell", cng, "", jetzt)
    }

    c.JSON(http.StatusOK, gin.H{"message": "Arbeitszeit aktualisiert", "changes": changes})}