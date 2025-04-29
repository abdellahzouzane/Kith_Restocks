package scraper

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type Product struct {
	ID         string    `json:"id"`    // Toujours traité comme string
	Name       string    `json:"name"`
	URL        string    `json:"url"`
	Price      string    `json:"price"`
	Image      string    `json:"image"`
	InStock    bool      `json:"in_stock"`
	DetectedAt time.Time `json:"detected_at"`
}

type KithScraper struct {
	db *sql.DB
}

func NewKithScraper(db *sql.DB) *KithScraper {
	return &KithScraper{db: db}
}

func (ks *KithScraper) ScrapeKithProducts() error {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get("https://eu.kith.com/collections/all/products.json")
	if err != nil {
		return fmt.Errorf("erreur HTTP: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("statut non-200: %s", resp.Status)
	}

	var data struct {
		Products []struct {
			ID       int64   `json:"id"`
			Title    string  `json:"title"`
			Handle   string  `json:"handle"`
			Variants []struct {
				Price    string `json:"price"`
				Available bool  `json:"available"`
			} `json:"variants"`
			Images []struct {
				Src string `json:"src"`
			} `json:"images"`
		} `json:"products"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("erreur décodage JSON: %v", err)
	}

	for _, p := range data.Products {
		if len(p.Variants) == 0 {
			continue
		}

		product := Product{
			ID:         fmt.Sprintf("%d", p.ID), // Conversion en string
			Name:       strings.TrimSpace(p.Title),
			URL:        fmt.Sprintf("https://eu.kith.com/products/%s", p.Handle),
			Price:      p.Variants[0].Price,
			InStock:    p.Variants[0].Available,
			DetectedAt: time.Now().UTC(), // Utilisation d'UTC pour la cohérence
		}

		if len(p.Images) > 0 {
			product.Image = strings.Replace(p.Images[0].Src, "//", "https://", 1)
		}

		// Vérification d'existence avec gestion d'erreur améliorée
		var exists bool
		err := ks.db.QueryRow(`
			SELECT EXISTS(
				SELECT 1 FROM products 
				WHERE id = $1
			)`, product.ID).Scan(&exists)
		
		if err != nil {
			log.Printf("[WARN] Erreur vérification produit %s: %v", product.ID, err)
			continue
		}

		// Transaction pour garantir l'intégrité des données
		tx, err := ks.db.Begin()
		if err != nil {
			log.Printf("[ERROR] Début transaction échoué: %v", err)
			continue
		}

		if exists {
			_, err = tx.Exec(`
				UPDATE products 
				SET 
					name = $1,
					url = $2,
					price = $3,
					image = $4,
					in_stock = $5, 
					detected_at = $6 
				WHERE id = $7`,
				product.Name,
				product.URL,
				product.Price,
				product.Image,
				product.InStock,
				product.DetectedAt,
				product.ID,
			)
		} else if product.InStock {
			_, err = tx.Exec(`
				INSERT INTO products (
					id, name, url, price, 
					image, in_stock, detected_at
				) VALUES (
					$1, $2, $3, $4, $5, $6, $7
				)`,
				product.ID,
				product.Name,
				product.URL,
				product.Price,
				product.Image,
				product.InStock,
				product.DetectedAt,
			)
		}

		if err != nil {
			tx.Rollback()
			log.Printf("[ERROR] Opération DB échouée: %v", err)
			continue
		}

		if err := tx.Commit(); err != nil {
			log.Printf("[ERROR] Commit transaction échoué: %v", err)
		}
	}

	log.Printf("[INFO] Scraping terminé - %d produits traités", len(data.Products))
	return nil
}