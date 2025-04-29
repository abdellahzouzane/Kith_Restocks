package main

import (
	"log"
	"time"
	"kith-monitor/scraper"
	"kith-monitor/database"
)

func main() {
	dbConfig := database.Config{
		Host:     "127.0.0.1", // ou "localhost"
		Port:     "5432",
		User:     "abdellah_user",
		Password: "password",
		DBName:   "db_user",
		SSLMode:  "disable",
	}

	db, err := database.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Erreur de connexion à la base de données: %v", err)
	}
	defer db.Close()

	kithScraper := scraper.NewKithScraper(db)
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Scraping en cours...")
		if err := kithScraper.ScrapeKithProducts(); err != nil {
			log.Printf("Erreur: %v", err)
		}
	}
}