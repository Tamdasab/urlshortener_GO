package api

import (
    "github.com/axellelanca/urlshortener/internal/services"
    "github.com/gin-gonic/gin"
)

// SetupRouter configure et retourne le routeur Gin avec toutes les routes
func SetupRouter(linkService *services.LinkService, clickService *services.ClickService) *gin.Engine {
    router := gin.Default()

    router.GET("/api/stats/:shortcode", GetLinkStatsHandler(linkService, clickService))

    return router
}
