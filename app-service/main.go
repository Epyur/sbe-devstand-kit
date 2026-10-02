package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Заготовка локального сервиса разработчика (L1 мини-стенд).
// Показывает обычную схему SBE: JWT приложения `dev`, собственные данные в
// своей БД (таблица items), никаких прав/ролей. Замените своим кодом.
type server struct {
	pool   *pgxpool.Pool
	secret []byte
	appID  string
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}
	secret := os.Getenv("DEV_SERVICE_SECRET")
	if secret == "" {
		log.Fatal("DEV_SERVICE_SECRET is required")
	}
	appID := os.Getenv("DEV_APP_ID")
	if appID == "" {
		appID = "dev"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("pgxpool.New: %v", err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("ping: %v", err)
	}
	if _, err := pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS items (
		id         SERIAL PRIMARY KEY,
		title      TEXT NOT NULL,
		author     TEXT NOT NULL DEFAULT '',
		created_at TIMESTAMPTZ NOT NULL DEFAULT now())`); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	s := &server{pool: pool, secret: []byte(secret), appID: appID}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/dev/health", s.handleHealth)
	mux.HandleFunc("GET /api/dev/items", s.requireAuth(s.handleListItems))
	mux.HandleFunc("POST /api/dev/items", s.requireAuth(s.handleCreateItem))

	log.Printf("app-service listening on :%s", port)
	if err := (&http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}).ListenAndServe(); err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}

func (s *server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if err := s.pool.Ping(r.Context()); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{"status": "error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok"})
}

type item struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	CreatedAt string `json:"created_at"`
}

func (s *server) handleListItems(w http.ResponseWriter, r *http.Request) {
	rows, err := s.pool.Query(r.Context(),
		`SELECT id, title, author, created_at FROM items ORDER BY id DESC LIMIT 200`)
	if err != nil {
		internalError(w, err)
		return
	}
	defer rows.Close()
	list := make([]item, 0)
	for rows.Next() {
		var it item
		var created time.Time
		if err := rows.Scan(&it.ID, &it.Title, &it.Author, &created); err != nil {
			internalError(w, err)
			return
		}
		it.CreatedAt = created.Format(time.RFC3339Nano)
		list = append(list, it)
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": list})
}

func (s *server) handleCreateItem(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid json"})
		return
	}
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "title is required"})
		return
	}
	email, _ := r.Context().Value(emailCtx{}).(string)
	var id int
	if err := s.pool.QueryRow(r.Context(),
		`INSERT INTO items (title, author) VALUES ($1, $2) RETURNING id`,
		req.Title, email).Scan(&id); err != nil {
		internalError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"id": id})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("writeJSON: %v", err)
	}
}

func internalError(w http.ResponseWriter, err error) {
	log.Printf("internal: %v", err)
	writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "internal error"})
}
