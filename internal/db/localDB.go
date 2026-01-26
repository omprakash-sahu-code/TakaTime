package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Rtarun3606k/TakaTime/internal/types"
	_ "modernc.org/sqlite"
)

func InitSQLite() (*sql.DB, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dbDir := filepath.Join(home, ".takatime")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(dbDir, "taka.db")

	// Open the DB connection
	sqliteDB, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Create table
	query := `
    CREATE TABLE IF NOT EXISTS offline_logs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        data TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	if _, err := sqliteDB.Exec(query); err != nil {
		return nil, err
	}

	// Return the active connection
	return sqliteDB, nil
}

func Enqueue(entry types.LogEntry, db *sql.DB) error {
	jsonData, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO offline_logs (data) VALUES (?)", string(jsonData))
	return err
}

func Flush(uploadFunc func(batch []types.LogEntry) error, db *sql.DB) bool {
	// 1. Fetch oldest 50 logs
	rows, err := db.Query("SELECT id, data FROM offline_logs ORDER BY id ASC LIMIT 50")
	if err != nil {
		log.Printf("SQLite Read Error: %v", err)
		return false
	}
	defer rows.Close()

	var batch []types.LogEntry
	var ids []int

	for rows.Next() {
		var id int
		var rawJSON string
		if err := rows.Scan(&id, &rawJSON); err == nil {
			var entry types.LogEntry
			if json.Unmarshal([]byte(rawJSON), &entry) == nil {
				batch = append(batch, entry)
				ids = append(ids, id)
			}
		}
	}

	// If queue is empty, we are done!
	if len(batch) == 0 {
		return false
	}

	// 2. ATTEMPT UPLOAD (Critical Step)
	log.Printf("Attempting to sync %d logs...", len(batch))
	err = uploadFunc(batch)

	if err != nil {
		// FAIL: Keep data, stop processing
		log.Printf("Upload failed (Data kept safely offline): %v", err)
		return false // Stop loop
	}

	// 3. SUCCESS: Delete these specific IDs
	deleteBatch(ids, db)
	log.Printf("Synced & Cleaned %d logs.", len(batch))

	// Return true to say "We processed a batch, maybe check if there is more?"
	return true
}

func deleteBatch(ids []int, db *sql.DB) {
	if len(ids) == 0 {
		return
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf("DELETE FROM offline_logs WHERE id IN (%s)", strings.Join(placeholders, ","))
	_, err := db.Exec(query, args...)
	if err != nil {
		log.Printf("Failed to delete synced logs: %v", err)
	}
}

func SyncQueue(mongoURI string, db *sql.DB) {
	for {
		// Flush calls our closure for every batch of 50 logs
		moreData := Flush(func(batch []types.LogEntry) error {
			return AddEntryToMongo(batch, mongoURI)
		}, db)

		if !moreData {
			break // Queue is empty or sync failed (loop stops)
		}
	}
}
