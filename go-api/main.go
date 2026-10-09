package main

import (
	"crypto/rand"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sort"
	"sync"
)

type Link struct {
	Code   string `json:"code"`
	URL    string `json:"url"`
	Clicks int    `json:"clicks"`
}

type Store struct {
	mu    sync.Mutex
	links map[string]*Link
}

func NewStore() *Store {
	return &Store{links: make(map[string]*Link)}
}

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func newCode() string {
	b := make([]byte, 7)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	for i := range b {
		b[i] = alphabet[int(b[i])%len(alphabet)]
	}
	return string(b)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (s *Store) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /links", s.create)
	mux.HandleFunc("GET /links", s.list)
	mux.HandleFunc("GET /links/{code}/stats", s.stats)
	mux.HandleFunc("GET /{code}", s.follow)
	return mux
}

func (s *Store) create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	code := newCode()
	for s.links[code] != nil {
		code = newCode()
	}
	l := &Link{Code: code, URL: body.URL}
	s.links[code] = l
	writeJSON(w, http.StatusCreated, l)
}

func (s *Store) list(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]*Link, 0, len(s.links))
	for _, l := range s.links {
		out = append(out, l)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Code < out[j].Code })
	writeJSON(w, http.StatusOK, out)
}

func (s *Store) stats(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	l := s.links[r.PathValue("code")]
	if l == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	writeJSON(w, http.StatusOK, l)
}

func (s *Store) follow(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	l := s.links[r.PathValue("code")]
	if l != nil {
		l.Clicks++
	}
	s.mu.Unlock()
	if l == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	http.Redirect(w, r, l.URL, http.StatusTemporaryRedirect)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("go-api listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, NewStore().Handler()))
}
