package web

import (
	"database/sql"
	"fmt"
)

type Download struct {
	ID           int64
	URL          string
	Title        string
	TargetDir    string
	Filename     string
	Format       string
	Status       string
	ErrorMessage string
	Progress     int
	CreatedAt    string
	UpdatedAt    string
}

const downloadColumns = `id, url, title, target_dir, filename, format, status, error_message, progress, created_at, updated_at`

func scanDownload(row interface{ Scan(...any) error }) (*Download, error) {
	var d Download
	err := row.Scan(
		&d.ID, &d.URL, &d.Title, &d.TargetDir, &d.Filename,
		&d.Format, &d.Status, &d.ErrorMessage, &d.Progress,
		&d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func queryDownload(db *sql.DB, id int64) (*Download, error) {
	return scanDownload(db.QueryRow(
		`SELECT `+downloadColumns+` FROM downloads WHERE id=?`, id,
	))
}

func queryHistory(db *sql.DB) ([]Download, error) {
	rows, err := db.Query(
		`SELECT ` + downloadColumns + ` FROM downloads
		 WHERE status NOT IN ('queued', 'downloading')
		 ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var downloads []Download
	for rows.Next() {
		d, err := scanDownload(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning download: %w", err)
		}
		downloads = append(downloads, *d)
	}
	return downloads, rows.Err()
}
