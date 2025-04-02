package middleware

import (
    "zeiterfassung-backend/utils"
    "net/http"

    "github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenStr := c.GetHeader("Authorization")
        if tokenStr == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Kein Token"})
            c.Abort()
            return
        }

        claims, err := utils.ParseToken(tokenStr)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Ungültiges Token"})
            c.Abort()
            return
        }

        c.Set("nutzer_id", claims.NutzerID)
        c.Set("rechte_id", claims.RechteID)
        c.Next()
    }
}
