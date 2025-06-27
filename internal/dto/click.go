package dto



type ClickCountOuput struct {

	LinkID    uint      `gorm:"index"`             // Clé étrangère vers la table 'links', indexée pour des requêtes efficaces

	Timestamp time.Time // Horodatage précis du clic
	UserAgent string    `gorm:"size:255"` // User-Agent de l'utilisateur qui a cliqué (informations sur le navigateur/OS)
	IPAddress string    `gorm:"size:50"`  // Adresse IP de l'utilisateur
}

