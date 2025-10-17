// internal/storage/storage.go
package storage

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)


type Entry struct {
	ID        int64
	Content   string
	Location  string
	CreatedAt time.Time
    Date time.Time
}

// Storage handles all database operations for Anima.
type Storage struct {
	db *sql.DB
}

// New initializes the database connection and creates necessary tables.
func New(dbPath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS entries (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		content TEXT NOT NULL,
		location TEXT,
		created_at DATETIME NOT NULL,
        date DATE  NOT NULL
	);
    CREATE INDEX IF NOT EXISTS idx_entries_date ON entries(date);
    `

	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

// Close closes the database connection.
func (s *Storage) Close() {
	s.db.Close()
}

// CreateEntry inserts a new journal entry into the database.
func (s *Storage) CreateEntry(content, location string) (*Entry, error) {
	now := time.Now().UTC()
    today := now.Truncate(24 * time.Hour)
	stmt, err := s.db.Prepare("INSERT INTO entries(content, location, created_at, date) VALUES(?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(content, location, now, today)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Entry{
		ID:        id,
		Content:   content,
		Location:  location,
		CreatedAt: now.Truncate(time.Second),
        Date: today,
	}, nil
}

// GetEntry retrieves a single entry by its ID.
func (s *Storage) GetEntry(id int64) (*Entry, error) {
	row := s.db.QueryRow("SELECT id, content, location, created_at, date FROM entries WHERE id = ?", id)

	var entry Entry
	err := row.Scan(&entry.ID, &entry.Content, &entry.Location, &entry.CreatedAt, &entry.Date)
	if err != nil {
		// This includes the case where no row is found (sql.ErrNoRows)
		return nil, err
	}
    // Truncate to handle potential database precision differences
    entry.CreatedAt = entry.CreatedAt.Truncate(time.Second)
    return &entry, nil
}

func (s *Storage) GetEntryByDate(date time.Time) (*Entry, error) {
	// Ensure we are only comparing the date part
	targetDate := date.Truncate(24 * time.Hour)
	row := s.db.QueryRow("SELECT id, content, location, created_at, date FROM entries WHERE date = ?", targetDate)

	var entry Entry
	err := row.Scan(&entry.ID, &entry.Content, &entry.Location, &entry.CreatedAt, &entry.Date)
	if err != nil {
		return nil, err // This will be sql.ErrNoRows if not found
	}
	entry.CreatedAt = entry.CreatedAt.Truncate(time.Second)
	return &entry, nil
}


func (s *Storage) UpdateEntryContent(id int64, content string) error {
	stmt, err := s.db.Prepare("UPDATE entries SET content = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(content, id)
	return err
}

