// utils/edit_time_limit.go
package utils

import (
    "time"
)

// IsEditAllowed prüft, ob ein Arbeitszeit-Eintrag noch bearbeitet werden darf
// Einträge können nur bis zu 3 Monate nach dem Datum bearbeitet werden
func IsEditAllowed(datum time.Time) bool {
    // Aktuelles Datum
    now := time.Now()
    
    // Datum vor 3 Monaten berechnen
    threeMonthsAgo := now.AddDate(0, -3, 0)
    
    // Prüfen, ob das Datum nicht älter als 3 Monate ist
    return datum.After(threeMonthsAgo) || datum.Equal(threeMonthsAgo.Truncate(24*time.Hour))
}

// GetEditTimeLimit gibt das früheste Datum zurück, für das noch Bearbeitungen erlaubt sind
func GetEditTimeLimit() time.Time {
    return time.Now().AddDate(0, -3, 0)
}