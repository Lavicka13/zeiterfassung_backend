// service/export_service.go - Korrigierte Version
package service

import (
    "zeiterfassung-backend/config"
    "zeiterfassung-backend/models"
    "time"
)

// ReportRecord für den CSV-Export
type ReportRecord struct {
    Datum        time.Time
    Anfangszeit  string
    Endzeit      string
    Arbeitszeit  float64
    Pause        int
}

// calculateWorkingHours berechnet die Arbeitszeit korrekt, auch über Mitternacht hinaus
func calculateWorkingHours(start, end time.Time, pause int) float64 {
    // Wenn Endzeit vor Startzeit liegt, füge einen Tag hinzu
    if end.Before(start) {
        end = end.Add(24 * time.Hour)
    }
    
    // Berechne die Gesamtzeit in Stunden
    totalHours := end.Sub(start).Hours()
    
    // Ziehe die Pause in Stunden ab
    pauseHours := float64(pause) / 60.0
    workingHours := totalHours - pauseHours
    
    // Stelle sicher, dass das Ergebnis nicht negativ ist
    if workingHours < 0 {
        workingHours = 0
    }
    
    return workingHours
}

// GetMonthlyReportData lädt alle Arbeitszeiten eines Monats
func GetMonthlyReportData(year int, month int, userID uint) ([]ReportRecord, float64, error) {
    var arbeitszeiten []models.Arbeitszeit

    // Zeitraum für den Monat
    startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
    endDate := startDate.AddDate(0, 1, 0)

    // Filter nach NutzerID
    err := config.DB.
        Where("datum >= ? AND datum < ? AND nutzer_id = ?", startDate, endDate, userID).
        Find(&arbeitszeiten).Error

    if err != nil {
        return nil, 0, err
    }

    var records []ReportRecord
    var summe float64 = 0

    for _, a := range arbeitszeiten {
        endzeit := "offen"
        dauer := 0.0

        if a.Endzeit != nil {
            endzeit = a.Endzeit.Format("15:04")
            
            // Verwende die neue Berechnungsfunktion
            dauer = calculateWorkingHours(a.Anfangszeit, *a.Endzeit, a.Pause)
            summe += dauer
        }

        record := ReportRecord{
            Datum:       a.Datum,
            Anfangszeit: a.Anfangszeit.Format("15:04"),
            Endzeit:     endzeit,
            Arbeitszeit: dauer,
            Pause:       a.Pause,
        }

        records = append(records, record)
    }

    return records, summe, nil
}

func GetYearlyReportData(year int, userID uint) ([]ReportRecord, float64, error) {
    var arbeitszeiten []models.Arbeitszeit

    startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
    endDate := startDate.AddDate(1, 0, 0)

    // Filter nach Nutzer
    err := config.DB.
        Where("datum >= ? AND datum < ? AND nutzer_id = ?", startDate, endDate, userID).
        Find(&arbeitszeiten).Error

    if err != nil {
        return nil, 0, err
    }

    var records []ReportRecord
    var summe float64 = 0

    for _, a := range arbeitszeiten {
        endzeit := "offen"
        dauer := 0.0

        if a.Endzeit != nil {
            endzeit = a.Endzeit.Format("15:04")
            
            // Verwende die neue Berechnungsfunktion
            dauer = calculateWorkingHours(a.Anfangszeit, *a.Endzeit, a.Pause)
            summe += dauer
        }

        record := ReportRecord{
            Datum:       a.Datum,
            Anfangszeit: a.Anfangszeit.Format("15:04"),
            Endzeit:     endzeit,
            Arbeitszeit: dauer,
            Pause:       a.Pause,
        }

        records = append(records, record)
    }

    return records, summe, nil
}