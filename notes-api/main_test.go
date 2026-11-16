package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// Lab-Werte, keine echten Secrets.
const labKey = "lab-api-key-123"

func TestLoadAPIKeyMissing(t *testing.T) {
	t.Setenv("NOTES_API_KEY_FILE", "")
	t.Setenv("NOTES_API_KEY", "")
	if _, _, err := loadAPIKey(); err == nil {
		t.Fatal("Fehler erwartet")
	}
}

func TestLoadFaults(t *testing.T) {
	t.Setenv("NOTES_DELAY", "300ms")
	t.Setenv("NOTES_ERROR_RATE", "1")
	f, err := loadFaults()
	if err != nil || f.delay.Milliseconds() != 300 || !f.active() {
		t.Fatalf("got %+v, %v", f, err)
	}
	f.delay = 0
	if f.inject() == nil {
		t.Fatal("Rate 1 muss immer scheitern")
	}
	t.Setenv("NOTES_DELAY", "kurz")
	if _, err := loadFaults(); err == nil {
		t.Fatal("Fehler erwartet")
	}
}

func TestPostNotesRequiresKey(t *testing.T) {
	h := newMux(newMemStore(), labKey, newMetrics("memory"), faults{})
	cases := []struct {
		key  string
		want int
	}{
		{"", http.StatusUnauthorized},
		{"falsch", http.StatusUnauthorized},
		{labKey, http.StatusCreated},
	}
	for _, c := range cases {
		req := httptest.NewRequest(http.MethodPost, "/notes",
			strings.NewReader(`{"text":"test"}`))
		req.Header.Set("X-API-Key", c.key)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != c.want {
			t.Errorf("key %q: status %d, want %d", c.key, rec.Code, c.want)
		}
	}
}

func TestErrorRate(t *testing.T) {
	h := newMux(newMemStore(), labKey, newMetrics("memory"), faults{errorRate: 1})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/notes", nil))
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status %d", rec.Code)
	}
}

func TestMetrics(t *testing.T) {
	h := newMux(newMemStore(), labKey, newMetrics("memory"), faults{})
	req := httptest.NewRequest(http.MethodPost, "/notes", strings.NewReader(`{"text":"x"}`))
	req.Header.Set("X-API-Key", labKey)
	h.ServeHTTP(httptest.NewRecorder(), req)
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/gibtsnicht", nil))

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	body := rec.Body.String()
	for _, want := range []string{
		`notes_build_info{store="memory",version="dev"} 1`,
		`notes_http_requests_total{code="201",method="POST",path="/notes"} 1`,
		`notes_http_requests_total{code="404",method="GET",path="other"} 1`,
		`notes_http_request_duration_seconds_count{method="POST",path="/notes"} 1`,
		`notes_created_total 1`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("fehlt: %s\n%s", want, body)
		}
	}
	if strings.Contains(body, "notes_db_up") {
		t.Error("notes_db_up nur im Postgres-Betrieb")
	}
}

// Integrationstest gegen echtes Postgres, nur wenn NOTES_TEST_DB_URL gesetzt ist.
func TestPostgresStore(t *testing.T) {
	url := os.Getenv("NOTES_TEST_DB_URL")
	if url == "" {
		t.Skip("NOTES_TEST_DB_URL nicht gesetzt")
	}
	ctx := context.Background()
	db, err := openDB(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := migrate(ctx, db); err != nil {
		t.Fatal(err)
	}
	st := &pgStore{db: db}
	n, err := st.Add(ctx, "integration")
	if err != nil {
		t.Fatal(err)
	}
	notes, err := st.List(ctx)
	if err != nil || len(notes) == 0 || notes[0].ID != n.ID {
		t.Fatalf("got %v, %v", notes, err)
	}
}
