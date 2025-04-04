package main

import (
    "zeiterfassung-backend/config"
    "zeiterfassung-backend/handlers"
    "zeiterfassung-backend/middleware"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "log"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Println("WARNUNG: Keine .env gefunden")
    }

    config.ConnectDB()

    r := gin.Default()
    r.Use(middleware.CORSMiddleware())

    r.POST("/register", handlers.Register)
    r.POST("/login", handlers.Login)

    // 🔐 Öffentlich zugänglicher Passwort-Vergessen-Endpunkt
    r.POST("/api/passwort-vergessen", handlers.PasswortVergessenHandler)

    // 🔒 Geschützte Routen
    protected := r.Group("/api")
    protected.Use(middleware.AuthMiddleware())
    {
        protected.GET("/me", handlers.Me)
        protected.GET("/mitarbeiter", handlers.GetMitarbeiter)
        protected.GET("/arbeitszeiten/:id", handlers.GetArbeitszeiten)
        protected.POST("/arbeitszeiten", handlers.CreateArbeitszeit)
        protected.PUT("/arbeitszeiten", handlers.UpdateArbeitszeit)

        protected.GET("/export/monat", handlers.ExportMonat)
        protected.GET("/export/jahr", handlers.ExportJahr)
        protected.GET("/export/monat/pdf", handlers.ExportMonatPDF)
        protected.GET("/export/jahr/pdf", handlers.ExportJahrPDF)
    }

    r.Run(":8080")
}
