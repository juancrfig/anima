package auth

import (
	"errors"
	"time"
)

// Clock defines an interface for getting the current time.
// This allows us to mock time.Now() in tests.
type Clock interface {
	Now() time.Time
}

// realClock is the standard, production implementation of Clock.
type realClock struct{}

func (realClock) Now() time.Time {
	return time.Now()
}

// ---

var ErrSessionExpired = errors.New("auth: session expired or not authenticated")

// Manager holds the in-memory session password and its expiry.
// It is not thread-safe; synchronization must be handled by the caller.
type Manager struct {
	clock       Clock
	password    []byte
	expiry      time.Time
	hasPassword bool
}

// NewManager creates a new, unauthenticated Manager with a given clock.
// This is used for testing.
func NewManager(c Clock) *Manager {
	return &Manager{
		clock: c,
	}
}

// New creates a production Manager that uses the real system clock.
func New() *Manager {
	return NewManager(realClock{})
}

// SetPassword stores the password and sets its expiry.
// A duration of 0 means the session never expires.
func (m *Manager) SetPassword(password []byte, duration time.Duration) {
	// Create a copy to store, so the caller's buffer can be cleared
	m.password = make([]byte, len(password))
	copy(m.password, password)
	
	m.hasPassword = true

	if duration == 0 {
		// Set to "zero time", our flag for non-expiring
		m.expiry = time.Time{}
	} else {
		m.expiry = m.clock.Now().Add(duration)
	}
}

// Clear removes the password from memory and invalidates the session.
// This is critical for security (e.g., on logout or exit).
func (m *Manager) Clear() {
	// Overwrite the password buffer with zeros before nil-ing it
	for i := range m.password {
		m.password[i] = 0
	}
	m.password = nil
	m.hasPassword = false
	m.expiry = time.Time{}
}

// IsAuthenticated checks if the session is currently valid.
func (m *Manager) IsAuthenticated() bool {
	if !m.hasPassword {
		return false
	}

	// Check for "never expires" flag (zero time)
	if m.expiry.IsZero() {
		return true
	}

	// Check if current time is before the expiry time
	return m.clock.Now().Before(m.expiry)
}

// GetPassword returns the password if the session is valid.
// If the session is expired, it clears the password from memory
// and returns ErrSessionExpired.
func (m *Manager) GetPassword() ([]byte, error) {
	if !m.IsAuthenticated() {
		// If it was authenticated but just expired, clear it.
		if m.hasPassword {
			m.Clear()
		}
		return nil, ErrSessionExpired
	}
	return m.password, nil
}