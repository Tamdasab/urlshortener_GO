package api

import (
	"errors"
	"github.com/axellelanca/urlshortener/cmd"
	"log"
	"net/http"
	"time"

	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm" // Pour gérer gorm.ErrRecordNotFound
)


// TODO Créer une variable ClickEventsChannel qui est un chan de type ClickEvent
// ClickEventsChannel est le channel global (ou injecté) utilisé pour envoyer les événements de clic
// aux workers asynchrones. Il est bufferisé pour ne pas bloquer les requêtes de redirection.
var ClickEventsChannel chan models.ClickEvent


// SetupRoutes configure toutes les routes de l'API Gin et injecte les dépendances nécessaires
func SetupRoutes(router *gin.Engine, linkService *services.LinkService) {
	// Le channel est initialisé ici.
	if ClickEventsChannel == nil {
		ClickEventsChannel = make(chan models.ClickEvent, cmd.Cfg.Analytics.BufferSize)
	}


	router.GET("/health", HealthCheckHandler)

	// Routes de l'API au format /api/v1/
	api := router.Group("/api/v1")
	{
		// POST /links
		api.POST("/links", CreateShortLinkHandler(linkService))
		
		// GET /links/:shortCode/stats
		api.GET("/links/:shortCode/stats", GetLinkStatsHandler(linkService))

		api.GET("/links/:shortCode", (linkService))
	}


	// TODO : Routes de l'API
	// Doivent être au format /api/v1/
	// POST /links
	// GET /links/:shortCode/stats
	router.GET("/api/v1/links/:shortCode", HealthCheckHandler)


	// Route de Redirection (au niveau racine pour les short codes)
	router.GET("/:shortCode", RedirectHandler(linkService))
}

// HealthCheckHandler gère la route /health pour vérifier l'état du service.
func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// CreateLinkRequest représente le corps de la requête JSON pour la création d'un lien.
type CreateLinkRequest struct {
	LongURL string `json:"long_url" binding:"required,url"` // 'binding:required' pour validation, 'url' pour format URL
}

// CreateShortLinkHandler gère la création d'une URL courte.
func CreateShortLinkHandler(linkService *services.LinkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateLinkRequest
		
		// Tente de lier le JSON de la requête à la structure CreateLinkRequest.
		// Gin gère la validation 'binding'.
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Appeler le LinkService (CreateLink) pour créer le nouveau lien.
		link, err := linkService.CreateLink(req.LongURL)
		if err != nil {
			log.Printf("Error creating link for URL %s: %v", req.LongURL, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
			return
		}

		// Retourne le code court et l'URL longue dans la réponse JSON.
		// Code HTTP 201 Created (nouvelle ressource)
		c.JSON(http.StatusCreated, gin.H{
			"short_code":     link.ShortCode,
			"long_url":       link.LongURL,
			"full_short_url": cmd.Cfg.Server.BaseURL + "/" + link.ShortCode, // Utiliser cfg.Server.BaseURL
		})
	}
}
// RedirectHandler gère la redirection d'une URL courte vers l'URL longue et l'enregistrement asynchrone des clics.
func RedirectHandler(linkService *services.LinkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Récupère le shortCode de l'URL avec c.Param
		shortCode :=c.Param("shortCode")

		link, err := linkService.GetLinkByShortCode(shortCode)
		if err != nil {
			// Si le lien n'est pas trouvé, retourner HTTP 404 Not Found.
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Short link not found"})
				return
			}
			log.Printf("Error retrieving link for %s: %v", shortCode, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		clickEvent := models.ClickEvent{
			LinkID:    link.ID,
			TimesTamp: time.Now(),
			UserAgent: c.GetHeader("User-Agent"),
			IPAddress:  c.ClientIP(),
		}

		select {
		case ClickEventsChannel <- clickEvent:
		default:
			log.Printf("Warning: ClickEventsChannel is full, dropping click event for %s.", shortCode)
		}

		c.Redirect(http.StatusFound, link.LongURL)

	}
}

// GetLinkStatsHandler gère la récupération des statistiques pour un lien spécifique.
func GetLinkStatsHandler(linkService *services.LinkService, , clickService *services.ClickService) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortCode := c.Param("shortcode")
		// TODO 6: Appeler le LinkService pour obtenir le lien et le nombre total de clics.

		link, err := linkService.GetLinkByShortCode(shortCode)
        if err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                c.JSON(http.StatusNotFound, gin.H{"error":  err.Error()})
                return
            }
            // Gérer d'autres erreurs
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
            return
        }

		// Récupération du total click de notre link
		totalClicks, err = clickService.GetClicksCountByLinkID(shortCode)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve click count"})
            return
        }

		}

		// Retourne les statistiques dans la réponse JSON.
		c.JSON(http.StatusOK, gin.H{
			"short_code":   link.ShortCode,
			"long_url":     link.LongURL,
			"total_clicks": totalClicks,
		})
	}
}
