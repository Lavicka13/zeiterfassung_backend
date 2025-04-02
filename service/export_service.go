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
            dauer = a.Endzeit.Sub(a.Anfangszeit).Hours() - float64(a.Pause)/60
            if dauer < 0 {
                dauer = 0
            }
            summe += dauer // ⬅️ summiert direkt mit
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
            dauer = a.Endzeit.Sub(a.Anfangszeit).Hours() - float64(a.Pause)/60
            if dauer < 0 {
                dauer = 0
            }
            summe += dauer // summiert gleich mit
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

