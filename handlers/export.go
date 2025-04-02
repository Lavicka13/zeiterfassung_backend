package handlers

import (
    "bytes"
    "encoding/csv"
    "fmt"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/phpdave11/gofpdf" // <--- PDF Bibliothek
    "zeiterfassung-backend/service"
)

// ------------------------
// Export Monatsbericht CSV
// ------------------------

func ExportMonat(c *gin.Context) {
    yearStr := c.Query("year")
    monthStr := c.Query("month")
    userStr := c.Query("user")
    nachname := c.Query("nachname")

    if yearStr == "" || monthStr == "" || userStr == "" || nachname == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter year, month, user oder nachname fehlt"})
        return
    }

    year, _ := strconv.Atoi(yearStr)
    month, _ := strconv.Atoi(monthStr)
    userID, _ := strconv.Atoi(userStr)

    records, summe, err := service.GetMonthlyReportData(year, month, uint(userID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Abrufen der Reportdaten"})
        return
    }

    // CSV erzeugen
    csvContent, err := generateCSV(records, summe)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Generieren der CSV"})
        return
    }

    filename := fmt.Sprintf("Monatsbericht_%s_%02d_%d.csv", nachname, month, year)
    c.Header("Content-Type", "text/csv")
    c.Header("Content-Disposition", "attachment; filename="+filename)
    c.Data(http.StatusOK, "text/csv", csvContent)
}

// ------------------------
// Export Monatsbericht PDF
// ------------------------


func ExportMonatPDF(c *gin.Context) {
    yearStr := c.Query("year")
    monthStr := c.Query("month")
    userStr := c.Query("user")
    nachname := c.Query("nachname")

    if yearStr == "" || monthStr == "" || userStr == "" || nachname == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter year, month, user oder nachname fehlt"})
        return
    }

    year, _ := strconv.Atoi(yearStr)
    month, _ := strconv.Atoi(monthStr)
    userID, _ := strconv.Atoi(userStr)

    records, summe, err := service.GetMonthlyReportData(year, month, uint(userID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Abrufen der Reportdaten"})
        return
    }

    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.AddUTF8Font("DejaVu", "", "./fonts/DejaVuSans.ttf")
    pdf.SetFont("DejaVu", "", 12)
    pdf.AddPage()

    // Titel
    pdf.CellFormat(0, 10, fmt.Sprintf("Monatsbericht für %s - %02d/%d", nachname, month, year), "", 1, "C", false, 0, "")
    pdf.Ln(4)

    // Tabellenkopf
    pdf.SetFont("DejaVu", "B", 11)
    pdf.CellFormat(35, 8, "Datum", "1", 0, "C", false, 0, "")
    pdf.CellFormat(25, 8, "Start", "1", 0, "C", false, 0, "")
    pdf.CellFormat(25, 8, "Ende", "1", 0, "C", false, 0, "")
    pdf.CellFormat(35, 8, "Arbeitszeit (h)", "1", 0, "C", false, 0, "")
    pdf.CellFormat(25, 8, "Pause (Min.)", "1", 1, "C", false, 0, "")

    // Tabelleninhalt
    pdf.SetFont("DejaVu", "", 10)
    for _, r := range records {
        pdf.CellFormat(35, 8, r.Datum.Format("02.01.2006"), "1", 0, "C", false, 0, "")
        pdf.CellFormat(25, 8, r.Anfangszeit, "1", 0, "C", false, 0, "")
        pdf.CellFormat(25, 8, r.Endzeit, "1", 0, "C", false, 0, "")
        pdf.CellFormat(35, 8, fmt.Sprintf("%.2f", r.Arbeitszeit), "1", 0, "C", false, 0, "")
        pdf.CellFormat(25, 8, fmt.Sprintf("%d", r.Pause), "1", 1, "C", false, 0, "")
    }

    // Summenzeile
    pdf.Ln(5)
    pdf.SetFont("DejaVu", "B", 11)
    pdf.CellFormat(115, 8, "Gesamtstunden", "1", 0, "R", false, 0, "")
    pdf.CellFormat(35, 8, fmt.Sprintf("%.2f Stunden", summe), "1", 1, "C", false, 0, "")

    // PDF zurückgeben
    c.Header("Content-Type", "application/pdf")
    c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=Monatsbericht_%s_%02d_%d.pdf", nachname, month, year))
    _ = pdf.Output(c.Writer)
}




// ------------------------
// Export Jahresbericht CSV
// ------------------------

func ExportJahr(c *gin.Context) {
    yearStr := c.Query("year")
    userStr := c.Query("user")
    nachname := c.Query("nachname")

    if yearStr == "" || userStr == "" || nachname == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter year, user oder nachname fehlt"})
        return
    }

    year, _ := strconv.Atoi(yearStr)
    userID, _ := strconv.Atoi(userStr)

    records, summe, err := service.GetYearlyReportData(year, uint(userID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Abrufen der Reportdaten"})
        return
    }

    csvContent, err := generateCSV(records, summe)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Generieren der CSV"})
        return
    }

    filename := fmt.Sprintf("Jahresbericht_%s_%d.csv", nachname, year)
    c.Header("Content-Type", "text/csv")
    c.Header("Content-Disposition", "attachment; filename="+filename)
    c.Data(http.StatusOK, "text/csv", csvContent)
}

// ------------------------
// CSV Generator
// ------------------------

func generateCSV(data []service.ReportRecord, summe float64) ([]byte, error) {
    var buf bytes.Buffer
    writer := csv.NewWriter(&buf)

    // Header
    header := []string{"Datum", "Anfangszeit", "Endzeit", "Arbeitszeit (h)", "Pause (Min.)"}
    if err := writer.Write(header); err != nil {
        return nil, err
    }

    // Daten
    for _, rec := range data {
        row := []string{
            rec.Datum.Format("2006-01-02"),
            rec.Anfangszeit,
            rec.Endzeit,
            fmt.Sprintf("%.2f", rec.Arbeitszeit),
            fmt.Sprintf("%d", rec.Pause),
        }
        if err := writer.Write(row); err != nil {
            return nil, err
        }
    }

    // Summenzeile
    writer.Write([]string{})
    writer.Write([]string{"", "", "Summe", fmt.Sprintf("%.2f", summe), ""})

    writer.Flush()
    return buf.Bytes(), writer.Error()
}


func ExportJahrPDF(c *gin.Context) {
    yearStr := c.Query("year")
    userStr := c.Query("user")
    nachname := c.Query("nachname")

    if yearStr == "" || userStr == "" || nachname == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter year, user oder nachname fehlt"})
        return
    }

    year, _ := strconv.Atoi(yearStr)
    userID, _ := strconv.Atoi(userStr)

    records, summe, err := service.GetYearlyReportData(year, uint(userID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Abrufen der Reportdaten"})
        return
    }

    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.AddPage()
    pdf.SetFont("Arial", "", 12)

    pdf.Cell(0, 10, fmt.Sprintf("Jahresbericht %s %d", nachname, year))
    pdf.Ln(12)

    for _, r := range records {
        line := fmt.Sprintf("%s | %s - %s | %.2f Std | %d min Pause",
            r.Datum.Format("02.01.2006"),
            r.Anfangszeit,
            r.Endzeit,
            r.Arbeitszeit,
            r.Pause,
        )
        pdf.Cell(0, 8, line)
        pdf.Ln(8)
    }

    pdf.Ln(10)
    pdf.Cell(0, 10, fmt.Sprintf("Gesamtsumme: %.2f Stunden", summe))

    c.Header("Content-Type", "application/pdf")
    c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=Jahresbericht_%s_%d.pdf", nachname, year))
    _ = pdf.Output(c.Writer)
}
