package api

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type jobStore struct {
	mu            sync.Mutex
	items         map[string]*jobItem
	ttl           time.Duration
	evictInterval time.Duration
}

type jobItem struct {
	data      []byte
	expiresAt time.Time
}

func newJobStore(ttl, evictInterval time.Duration) *jobStore {
	s := &jobStore{
		items:         make(map[string]*jobItem),
		ttl:           ttl,
		evictInterval: evictInterval,
	}
	go s.evictLoop()
	return s
}

func (s *jobStore) save(data []byte) string {
	id := fmt.Sprintf("%d", time.Now().UnixNano())
	s.mu.Lock()
	s.items[id] = &jobItem{data: data, expiresAt: time.Now().Add(s.ttl)}
	s.mu.Unlock()
	return id
}

func (s *jobStore) get(id string) ([]byte, bool) {
	s.mu.Lock()
	item, ok := s.items[id]
	s.mu.Unlock()
	if !ok {
		return nil, false
	}
	return item.data, true
}

func (s *jobStore) evictLoop() {
	t := time.NewTicker(s.evictInterval)
	defer t.Stop()
	for range t.C {
		now := time.Now()
		s.mu.Lock()
		for id, item := range s.items {
			if now.After(item.expiresAt) {
				delete(s.items, id)
			}
		}
		s.mu.Unlock()
	}
}

func (h *Handler) getResult(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	data, ok := h.store.get(id)
	if !ok {
		http.Error(w, "not found or expired", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", `attachment; filename="music.zip"`)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
