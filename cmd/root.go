package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/axellelanca/urlshortener/internal/config"
	"github.com/spf13/cobra"
)

// Cfg est la variable globale qui contient la configuration chargée
var Cfg *config.Config

// RootCmd représente la commande de base
var RootCmd = &cobra.Command{
	Use:   "url-shortener",
	Short: "Un service de raccourcissement d'URLs avec API REST et CLI",
	Long: `'url-shortener' est une application complète pour gérer des URLs courtes.
Elle inclut un serveur API pour le raccourcissement et la redirection,
ainsi qu'une interface en ligne de commande pour l'administration.

Utilisez 'url-shortener [command] --help' pour plus d'informations sur une commande.`,
}

// Execute est le point d'entrée principal pour l'application Cobra
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Erreur lors de l'exécution de la commande: %v\n", err)
		os.Exit(1)
	}
}

// init initialise la configuration globale
func init() {

	// Configurer l'initialisation de la configuration
	cobra.OnInitialize(initConfig)
	// TODO Initialiser la configuration globale avec OnInitialize

	// IMPORTANT : Ici, nous n'appelons PAS RootCmd.AddCommand() directement
	// pour les commandes 'server', 'create', 'stats', 'migrate'.
	// Ces commandes s'enregistreront elles-mêmes via leur propre fonction init().
	//
	//rootCmd.AddCommand(cli.StatsCmd)
	// Assurez-vous que tous les fichiers de commande comme
	// 'cmd/server/server.go' et 'cmd/cli/*.go' aient bien
	// un `import "url-shortener/cmd"`
	// et un `func init() { cmd.RootCmd.AddCommand(MaCommandeCmd) }`
	// C'est ce qui va faire le lien !

}

// initConfig charge la configuration de l'application
func initConfig() {
	var err error
	Cfg, err = config.LoadConfig()
	if err != nil {
		// LoadConfig gère déjà les fichiers manquants avec des valeurs par défaut
		// On log juste un avertissement pour les autres erreurs
		log.Printf("Attention: Problème lors du chargement de la configuration: %v", err)
	}
}
