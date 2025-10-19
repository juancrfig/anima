package auth

import (
	"bytes"
	"testing"
	"time"
)


// The test implementation of the Clock interface
// It allows us to control time in our tests
type mockClock struct {
	currentTime time.Time
}


func (m *mockClock) Now() time.Time {
	return m.currentTime
}


func (m *mockClock) Advance(d time.Duration) {
	m.currentTime = m.currentTime.Add(d)
}


// --- Test Cases ---

func TestAuthManager_InitialState(t *testing.T) {
	mc := &mockClock{currentTime: time.Now()}
	mgr := NewManager(mc) // Use the test constructor

	if mgr.IsAuthenticated() {
		t.Fatal("manager should not be authenticated on init")
	}

	_, err := mgr.GetPassword()
	if err == nil {
		t.Fatal("GetPassword() should return an error on init")
	}
	if err != ErrSessionExpired {
		t.Fatalf("got %v, want %v", err, ErrSessionExpired)
	}
	t.Logf("Got expected error: %v", err)
}

func TestAuthManager_SetPassword(t *testing.T) {
	mc := &mockClock{currentTime: time.Now()}
	mgr := NewManager(mc)

	password := []byte("my-session-password")
	duration := 10 * time.Minute

	mgr.SetPassword(password, duration)

	if !mgr.IsAuthenticated() {
		t.Fatal("manager should be authenticated after setting password")
	}

	retrieved, err := mgr.GetPassword()
	if err != nil {
		t.Fatalf("GetPassword() returned an error: %v", err)
	}

	if !bytes.Equal(password, retrieved) {
		t.Fatal("GetPassword() returned wrong password")
	}
}


func TestAuthManager_SessionExpiry(t *testing.T) {
	startTime := time.Now()
	mc := &mockClock{currentTime: startTime}
	mgr := NewManager(mc)

	password := []byte("my-session-password")
	duration := 30 * time.Minute

	// 1. Set the password
	mgr.SetPassword(password, duration)
	if !mgr.IsAuthenticated() {
		t.Fatal("manager should be authenticated")
	}

	// 2. Advance time by 29 minutes
	mc.Advance(29 * time.Minute)
	if !mgr.IsAuthenticated() {
		t.Fatal("manager should still be authenticated after 29m")
	}

	// 3. Advance time past the expiry (total 31m)
	mc.Advance(2 * time.Minute) // Total 31m
	if mgr.IsAuthenticated() {
		t.Fatal("manager should be expired after 31m")
	}

	// 4. GetPassword should now fail
	_, err := mgr.GetPassword()
	if err == nil {
		t.Fatal("GetPassword() should fail after session expiry")
	}
	if err != ErrSessionExpired {
		t.Fatalf("got %v, want %v", err, ErrSessionExpired)
	}
	t.Logf("Got expected expiry error: %v", err)

	// 5. Check that password was cleared from memory
	if mgr.password != nil {
		t.Fatal("manager did not clear expired password from memory")
	}
}


func TestAuthManager_Clear(t *testing.T) {
	mc := &mockClock{currentTime: time.Now()}
	mgr := NewManager(mc)

	mgr.SetPassword([]byte("pass"), 10*time.Minute)
	if !mgr.IsAuthenticated() {
		t.Fatal("manager should be authenticated")
	}

	mgr.Clear()

	if mgr.IsAuthenticated() {
		t.Fatal("manager should not be authenticated after Clear()")
	}
	_, err := mgr.GetPassword()
	if err == nil {
		t.Fatal("GetPassword() should fail after Clear()")
	}
}


// TestAuthManager_ZeroDuration proves that a 0 duration
// creates a non-expiring session.
func TestAuthManager_ZeroDuration(t *testing.T) {
	startTime := time.Now()
	mc := &mockClock{currentTime: startTime}
	mgr := NewManager(mc)

	password := []byte("my-session-password")
	duration := 0 * time.Minute // 0 duration means "no expiry"

	mgr.SetPassword(password, duration)

	// Advance time by 1000 hours; it should not expire
	mc.Advance(1000 * time.Hour)
	if !mgr.IsAuthenticated() {
		t.Fatal("manager should not expire with 0 duration")
	}

	retrieved, err := mgr.GetPassword()
	if err != nil {
		t.Fatalf("GetPassword() failed for 0 duration: %v", err)
	}
	if !bytes.Equal(password, retrieved) {
		t.Fatal("GetPassword() returned wrong password")
	}
}
