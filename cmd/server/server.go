package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	cmd2 "github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/api"
	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/repository"
	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// RunServerCmd représente la commande 'run-server' de Cobra.
var RunServerCmd = &cobra.Command{
	Use:   "run-server",
	Short: "Lance le serveur API de raccourcissement d'URLs",
	Long:  `Cette commande initialise la base de données, configure les APIs et lance le serveur HTTP.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Charger la configuration globale
		cfg := cmd2.Cfg
		if cfg == nil {
			log.Fatalf("ERREUR: Configuration non chargée")
		}

		// Initialiser la connexion à la base de données SQLite
		db, err := gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
		if err != nil {
			log.Fatalf("ERREUR: Impossible de se connecter à la base de données: %v", err)
		}

		// Auto-migrer les modèles GORM
		err = db.AutoMigrate(&models.Link{}, &models.Click{})
		if err != nil {
			log.Fatalf("ERREUR: Échec de la migration automatique: %v", err)
		}
		log.Println(" Migration automatique des modèles terminée avec succès")

		// Initialiser les repositories
		linkRepo := repository.NewLinkRepository(db)
		log.Println(" Repositories initialisés")

		// Initialiser les services métiers
		linkService := services.NewLinkService(linkRepo)
		log.Println(" Services métiers initialisés")

		// Configurer le routeur Gin et les routes
		router := gin.Default()
		api.SetupRoutes(router, linkService)
		log.Println(" Routes API configurées")

		// Créer le serveur HTTP
		serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
		srv := &http.Server{
			Addr:    serverAddr,
			Handler: router,
		}

		// Démarrer le serveur dans une goroutine
		go func() {
			log.Printf(" Serveur démarré sur le port %d", cfg.Server.Port)
			log.Printf(" API disponible sur: http://localhost:%d/api/v1/links", cfg.Server.Port)
			log.Printf("  Health check: http://localhost:%d/health", cfg.Server.Port)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("ERREUR: Impossible de démarrer le serveur: %v", err)
			}
		}()

		// Gestion de l'arrêt gracieux
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// Bloquer jusqu'à réception d'un signal d'arrêt
		<-quit
		log.Println(" Signal d'arrêt reçu. Arrêt du serveur...")

		// Arrêt propre avec timeout
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("ERREUR lors de l'arrêt du serveur: %v", err)
		}

		log.Println(" Serveur arrêté proprement")
	},
}

func init() {
	cmd2.RootCmd.AddCommand(RunServerCmd)
}
