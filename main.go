package main

import (
    "zeiterfassung-backend/config"
    "zeiterfassung-backend/handlers"
    "zeiterfassung-backend/middleware"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "log"
    "os"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Println("WARNUNG: Keine .env gefunden")
    }

    config.ConnectDB()

    // Bestimme den Laufmodus
    mode := os.Getenv("GIN_MODE")
    if mode == "release" {
        gin.SetMode(gin.ReleaseMode)
    }

    r := gin.Default()
    r.Use(middleware.CORSMiddleware())

    // Öffentliche Routen
    r.POST("/register", handlers.Register)
    r.POST("/login", handlers.Login)

    // Passwort-Vergessen-Endpunkt
    r.POST("/api/passwort-vergessen", handlers.PasswortVergessenHandler)

    // Geschützte Routen
    protected := r.Group("/api")
    protected.Use(middleware.AuthMiddleware())
    {
        protected.GET("/me", handlers.Me)
        
        // Nutzerverwaltung
        protected.GET("/mitarbeiter", handlers.GetMitarbeiter)
        protected.POST("/mitarbeiter", handlers.CreateNutzer)     // Neu
        protected.PUT("/mitarbeiter/:id", handlers.UpdateNutzer)  // Neu
        protected.DELETE("/mitarbeiter/:id", handlers.DeleteNutzer) // Neu
        protected.PUT("/mitarbeiter/:id/passwort", handlers.ResetPassword) // Neu
        
        // Arbeitszeiten
        protected.GET("/arbeitszeiten/:id", handlers.GetArbeitszeiten)
        protected.POST("/arbeitszeiten", handlers.CreateArbeitszeit)
        protected.PUT("/arbeitszeit/update", handlers.UpdateArbeitszeit)

        // Export-Funktionen
        protected.GET("/export/monat", handlers.ExportMonat)
        protected.GET("/export/jahr", handlers.ExportJahr)
        protected.GET("/export/monat/pdf", handlers.ExportMonatPDF)
        protected.GET("/export/jahr/pdf", handlers.ExportJahrPDF)
    }

    // Server starten
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    log.Println("Server startet auf Port", port)
    r.Run(":" + port)
}