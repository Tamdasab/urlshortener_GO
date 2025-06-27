package models

import (
    "gorm.io/gorm"
)

// TODO : Créer la struct Link
// Link représente un lien raccourci dans la base de données.
// Les tags `gorm:"..."` définissent comment GORM doit mapper cette structure à une table SQL.
// ID qui est une primaryKey
// Shortcode : doit être unique, indexé pour des recherches rapide (voir doc), taille max 10 caractères
// LongURL : doit pas être null
// CreateAt : Horodatage de la créatino du lien

type Link struct {
	gorm.Model
	ShortCode string `json:"short_code" gorm:"unique;not null"`
    LongURL   string `json:"long_url" gorm:"not null"`
}