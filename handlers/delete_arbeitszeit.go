// handlers/delete_arbeitszeit.go
package handlers

import (
    "net/http"
    "zeiterfassung-backend/config"
    "zeiterfassung-backend/models"

    "github.com/gin-gonic/gin"
)

// DeleteArbeitszeit löscht einen Arbeitszeiteintrag anhand seiner ID
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
    
    // Lösche die Arbeitszeit
    if err := config.DB.Delete(&arbeitszeit).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Löschen fehlgeschlagen"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Arbeitszeit erfolgreich gelöscht"})
}