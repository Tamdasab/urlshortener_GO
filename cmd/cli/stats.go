package cli

import (
	"fmt"
	"log"
	"os"

	cmd2 "github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/repository"
	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/spf13/cobra"
	"github.com/axellelanca/urlshortener/internal/models"
	"gorm.io/driver/sqlite" // Driver SQLite pour GORM
	"gorm.io/gorm"
)

// TODO : variable shortCodeFlag qui stockera la valeur du flag --code
var shortCodeFlag string

// StatsCmd représente la commande 'stats'
var StatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Affiche les statistiques (nombre de clics) pour un lien court.",
	Long: `Cette commande permet de récupérer et d'afficher le nombre total de clics
pour une URL courte spécifique en utilisant son code.

Exemple:
  url-shortener stats --code="xyz123"`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO : Valider que le flag --code a été fourni.
		// os.Exit(1) si erreur
		if shortCodeFlag == "" {
            fmt.Println("Erreur: Le flag --code est requis")
            os.Exit(1)
        }

		// TODO : Charger la configuration chargée globalement via cmd.cfg
		cfg := cmd2.GetConfig()
		

		// TODO 3: Initialiser la connexion à la base de données SQLite avec GORM.
		// log.Fatalf si erreur
		db, err := gorm.Open(sqlite.Open(cfg.Database.Path), &gorm.Config{})
        if err != nil {
            log.Fatalf("FATAL: Échec de la connexion à la base de données: %v", err)
        }

		sqlDB, err := db.DB()
			if err != nil {
				log.Fatalf("FATAL: Échec de l'obtention de la base de données SQL sous-jacente: %v", err)
			}

		// TODO S'assurer que la connexion est fermée à la fin de l'exécution de la commande
		// fermeture de la connexion après la fin de l'éxecution de la fonction
		defer sqlDB.Close()

		// TODO : Initialiser les repositories et services nécessaires NewLinkRepository & NewLinkService
        linkRepo := repository.NewGormLinkRepository(db)
        linkService := services.NewLinkService(linkRepo)

		clickRepo := repository.NewGormClickRepository(db)
        clickService := services.NewClickService(clickRepo)

		// TODO 5: Appeler GetLinkStats pour récupérer le lien et ses statistiques.
		// Attention, la fonction retourne 3 valeurs
		// Pour l'erreur, utilisez gorm.ErrRecordNotFound
		// Si erreur, os.Exit(1)
        link, totalClicks, err := GetLinkStats(linkService, clickService, shortCodeFlag)
        if err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                fmt.Printf("Erreur: Aucun lien trouvé pour le code: %s\n", shortCodeFlag)
            } else {
                fmt.Printf("Erreur: %v\n", err)
            }
            os.Exit(1)
        }

		fmt.Printf("Statistiques pour le code court: %s\n", link.ShortCode)
		fmt.Printf("URL longue: %s\n", link.LongURL)
		fmt.Printf("Total de clics: %d\n", totalClicks)
	},
}

// init() s'exécute automatiquement lors de l'importation du package.
// Il est utilisé pour définir les flags que cette commande accepte.


// GetLinkStats récupère un lien et ses statistiques de clics
func GetLinkStats(linkService *services.LinkService, clickService *services.ClickService, shortCode string) (*models.Link, int, error) {
    // Récupérer le lien
    link, err := linkService.GetLinkByShortCode(shortCode)
    if err != nil {
        return nil, 0, err
    }

    // Récupérer le nombre de clics
    totalClicks, err := clickService.GetClicksCountByLinkID(link.ID)
    if err != nil {
        return nil, 0, err
    }

    return link, totalClicks, nil
}

func init() {
    // Ajouter le flag --code à la commande
    StatsCmd.Flags().StringVar(&shortCodeFlag, "code", "", "Code court du lien à analyser (requis)")
    StatsCmd.MarkFlagRequired("code")
}
