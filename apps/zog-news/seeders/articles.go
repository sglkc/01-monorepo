package seeders

import (
	"database/sql"
)

func SeedArticles(db *sql.DB) error {
    // Insert users with password hashes
    _, err := db.Exec(`
        INSERT INTO articles (title, content, author, status) VALUES
        ('Great news title', 'Artikel ini membahas tentang perubahan teknologi...', 'Seya', 'published'),
        ('Bad clickbait', 'Perkembangan teknologi yang pesat...', 'Cikal', 'draft')
        ON CONFLICT DO NOTHING;
    `)

    return err
}
