package dto

import "time"

// UrlStatsDTO représente les statistiques d'une URL pour l'affichage
type StatsOuput struct {
    ShortCode  string    `json:"short_code"`
    LongURL    string    `json:"long_url"`
    ClickCount int       `json:"click_count"`
}

