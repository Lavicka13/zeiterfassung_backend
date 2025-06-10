// handlers/delete_arbeitszeit.go - Aktualisierte Version
package handlers

import (
    "net/http"
    "zeiterfassung-backend/config"
    "zeiterfassung-backend/models"
    "zeiterfassung-backend/utils"

    "github.com/gin-gonic/gin"
)

// DeleteArbeitszeit löscht einen Arbeitszeiteintrag anhand seiner ID
// Berücksichtigt dabei das 3-Monats-Limit
func DeleteArbeitszeit(c *gin.Context) {
    id := c.Param("id")
    
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Keine ID angegeben"})
        return
    }
    
    // Prüfe, ob die Arbeitszeit existiert
    var arbeitszeit models.Arbeitszeit
    if err := config.DB.First(&arbeitszeit, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Arbeitszeit nicht gefunden"})
        return
    }
    
    // Prüfen, ob der Eintrag noch gelöscht werden darf (3-Monats-Limit)
    if !utils.IsEditAllowed(arbeitszeit.Datum) {
        c.JSON(http.StatusForbidden, gin.H{
            "error": "Dieser Eintrag ist älter als 3 Monate und kann nicht mehr gelöscht werden.",
        })
        return
    }
    
    // Beginne eine Transaktion
    tx := config.DB.Begin()
    
    // Lösche zuerst alle zugehörigen Audit-Einträge
    if err := tx.Exec("DELETE FROM arbeitszeiten_audit WHERE arbeitszeit_id = ?", id).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Löschen der Audit-Einträge"})
        return
    }
    
    // Lösche dann die Arbeitszeit
    if err := tx.Delete(&arbeitszeit).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Löschen fehlgeschlagen"})
        return
    }
    
    // Schließe die Transaktion erfolgreich ab
    if err := tx.Commit().Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaktion konnte nicht abgeschlossen werden"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Arbeitszeit erfolgreich gelöscht"})
}