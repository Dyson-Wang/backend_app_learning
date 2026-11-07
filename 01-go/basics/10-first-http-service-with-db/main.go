package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Store 抽象数据层：换成 PostgreSQL 时只换实现，不动 handler
type Store interface {
	List() ([]User, error)
	Create(name string) (User, error)
	Get(id int64) (User, error)
}

type MemoryStore struct {
	mu     sync.RWMutex
	nextID int64
	users  map[int64]User
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{nextID: 1, users: make(map[int64]User)}
}

func (s *MemoryStore) List() ([]User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]User, 0, len(s.users))
	for _, u := range s.users {
		out = append(out, u)
	}
	return out, nil
}

func (s *MemoryStore) Create(name string) (User, error) {
	if strings.TrimSpace(name) == "" {
		return User{}, errors.New("name is required")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	u := User{ID: s.nextID, Name: name}
	s.users[u.ID] = u
	s.nextID++
	return u, nil
}

func (s *MemoryStore) Get(id int64) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[id]
	if !ok {
		return User{}, errors.New("user not found")
	}
	return u, nil
}

type API struct {
	store Store
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type response struct {
	Data  interface{} `json:"data,omitempty"`
	Error *apiError   `json:"error,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, body response) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("write response: %v", err)
	}
}

func (a *API) usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		users, err := a.store.List()
		if err != nil {
			writeJSON(w, http.StatusInternalServerError,
				response{Error: &apiError{Code: "INTERNAL", Message: err.Error()}})
			return
		}
		writeJSON(w, http.StatusOK, response{Data: users})
	case http.MethodPost:
		var in struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			writeJSON(w, http.StatusBadRequest,
				response{Error: &apiError{Code: "BAD_JSON", Message: err.Error()}})
			return
		}
		u, err := a.store.Create(in.Name)
		if err != nil {
			writeJSON(w, http.StatusBadRequest,
				response{Error: &apiError{Code: "VALIDATION", Message: err.Error()}})
			return
		}
		writeJSON(w, http.StatusCreated, response{Data: u})
	default:
		w.Header().Set("Allow", "GET, POST")
		writeJSON(w, http.StatusMethodNotAllowed,
			response{Error: &apiError{Code: "METHOD_NOT_ALLOWED", Message: "use GET or POST"}})
	}
}

func (a *API) userByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed,
			response{Error: &apiError{Code: "METHOD_NOT_ALLOWED", Message: "use GET"}})
		return
	}
	raw := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest,
			response{Error: &apiError{Code: "BAD_ID", Message: "id must be an integer"}})
		return
	}
	u, err := a.store.Get(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound,
			response{Error: &apiError{Code: "NOT_FOUND", Message: err.Error()}})
		return
	}
	writeJSON(w, http.StatusOK, response{Data: u})
}

// withRequestID：日志中间件 + trace id + panic 兜底
func withRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = strconv.FormatInt(time.Now().UnixNano(), 36)
		}
		w.Header().Set("X-Request-ID", id)
		start := time.Now()
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic=%v request_id=%s", rec, id)
				writeJSON(w, http.StatusInternalServerError,
					response{Error: &apiError{Code: "INTERNAL", Message: "internal error"}})
				return
			}
			log.Printf("method=%s path=%s latency=%s request_id=%s",
				r.Method, r.URL.Path, time.Since(start), id)
		}()
		next.ServeHTTP(w, r)
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	api := &API{store: NewMemoryStore()}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ready")
	})
	mux.HandleFunc("/users", api.usersHandler)
	mux.HandleFunc("/users/", api.userByIDHandler)

	// 所有超时都显式设置：不设 ReadHeaderTimeout 会有 Slowloris 风险
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           withRequestID(mux),
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("SIGTERM received, draining")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
	log.Println("stopped")
}
