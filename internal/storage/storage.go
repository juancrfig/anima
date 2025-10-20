package storage

import (
	"database/sql"
	"time"
    "fmt"

    "anima/internal/auth"
    "anima/internal/config"
    "anima/internal/crypto"

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
	db               *sql.DB
    auth             *auth.Manager
    cryptoParams     *crypto.Params
}


func (s *Storage) DB() *sql.DB {
    return s.db
}


func (s *Storage) getKey() ([]byte, error) {
    key, err := s.auth.GetPassword()
    if err != nil {
        return nil, fmt.Errorf("Not authenticated: %w", err)
    }
    return key, nil
}


func (s *Storage) SetCryptoParamsForTesting(params *crypto.Params) {
    s.cryptoParams = params
}


// New initializes the database connection and creates necessary tables.
func New(dbPath string, cfg *config.Config, authMgr *auth.Manager) (*Storage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS entries (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		content BLOB NOT NULL,
		location TEXT,
		created_at DATETIME NOT NULL,
        date DATE  NOT NULL
	);
    CREATE INDEX IF NOT EXISTS idx_entries_date ON entries(date);
    `

	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, err
	}

    params, err := cfg.CryptoParams()
    if err != nil {
        return nil, fmt.Errorf("could not load crypto params from config: %w", err)
    }

	return &Storage{
        db:             db,
        auth:           authMgr,
        cryptoParams:   params,
    }, nil
}

// Close closes the database connection.
func (s *Storage) Close() {
	s.db.Close()
}

// CreateEntry inserts a new journal entry into the database.
func (s *Storage) CreateEntry(content, location string, entryDate time.Time) (*Entry, error) {

    // Get session key
    key, err := s.getKey()
    if err != nil {
        return nil, err
    }

    // Encrypt content
    encryptedContent, err := crypto.Encrypt([]byte(content), key, s.cryptoParams)
    if err != nil {
        return nil, fmt.Errorf("Could not encrypt entry: %w", err)
    }

	now := time.Now().UTC().Truncate(time.Minute)
    dateOnly := entryDate.UTC().Truncate(24 * time.Hour)

	stmt, err := s.db.Prepare("INSERT INTO entries(content, location, created_at, date) VALUES(?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(encryptedContent, location, now, dateOnly)
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
        Date: dateOnly,
	}, nil
}


// A helper to scan a row and decrypt the content
func (s *Storage) scanEntry(row *sql.Row) (*Entry, error) {
    var entry Entry
    var encryptedContent []byte

    err := row.Scan(&entry.ID, &encryptedContent, &entry.Location, &entry.CreatedAt, &entry.Date)
    if err != nil {
        return nil, err
    }

    key, err := s.getKey()
    if err != nil {
        return nil, err
    }

    plaintext, err := crypto.Decrypt(encryptedContent, key)
    if err != nil {
        return nil, fmt.Errorf("Could not decrypt entry: %w", err)
    }
    entry.Content = string(plaintext)

    entry.CreatedAt = entry.CreatedAt.UTC().Truncate(time.Minute)
    return &entry, nil
}


// GetEntry retrieves and decrypts a single entry.
func (s *Storage) GetEntry(id int64) (*Entry, error) {
	row := s.db.QueryRow("SELECT id, content, location, created_at, date FROM entries WHERE id = ?", id)
    return s.scanEntry(row)
}


func (s *Storage) GetEntryByDate(date time.Time) (*Entry, error) {
	// Standardize the lookup date to UTC before truncating
	targetDate := date.UTC().Truncate(24 * time.Hour)
    row := s.db.QueryRow("SELECT id, content, location, created_at, date FROM entries WHERE date = ?", targetDate)
    return s.scanEntry(row)
}


func (s *Storage) UpdateEntryContent(id int64, content string) error {
    key, err := s.getKey()
    if err != nil {
        return err
    }

    encryptedContent, err := crypto.Encrypt([]byte(content), key, s.cryptoParams)
    if err != nil {
        return fmt.Errorf("Could not encrypt entry for update: %w", err)
    }

	stmt, err := s.db.Prepare("UPDATE entries SET content = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(encryptedContent, id)
	return err
}


func (s *Storage) GetOrCreateEntryByDate(date time.Time, location string) (*Entry, bool, error) {
    targetDate := date.Truncate(24 * time.Hour)

    entry, err := s.GetEntryByDate(targetDate)
    if err == nil {
        return entry, false, nil
    }

    if err != sql.ErrNoRows {
        return nil, false, err
    }

    newEntry, err := s.CreateEntry("", location, targetDate)
    if err != nil {
        return nil, false, err
    }
    return newEntry, true, nil
}
