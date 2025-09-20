package server

import (
	"encoding/json"
	"errors"
	"github.com/AndreySirin/-Effective-Mobile-/internal/entity"
	"github.com/AndreySirin/-Effective-Mobile-/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (s *Server) CreateSubs(w http.ResponseWriter, r *http.Request) {
	var req entity.SubsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	subs, err := entity.SubsToDataBase(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := s.storage.CreateSubs(r.Context(), &subs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) ReadSubs(w http.ResponseWriter, r *http.Request) {
	subsID := chi.URLParam(r, "id")
	id, err := uuid.Parse(subsID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if id == uuid.Nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	subs, err := s.storage.ReadSubs(r.Context(), id)
	if errors.Is(err, storage.ErrNotFound) {
		http.Error(w, "subscription not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(subs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) UpdateSubs(w http.ResponseWriter, r *http.Request) {
	var req entity.SubsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	subs, err := entity.SubsToDataBase(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	subsID := chi.URLParam(r, "id")
	id, err := uuid.Parse(subsID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if id == uuid.Nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	err = s.storage.UpdateSubs(r.Context(), id, &subs)
	if errors.Is(err, storage.ErrNotFound) {
		http.Error(w, "subscription not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) DeleteSubs(w http.ResponseWriter, r *http.Request) {
	subsID := chi.URLParam(r, "id")
	id, err := uuid.Parse(subsID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.storage.DeleteSubs(r.Context(), id)
	if errors.Is(err, storage.ErrNotFound) {
		http.Error(w, "subscription not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) ListSubs(w http.ResponseWriter, r *http.Request) {
	var date string
	if err := json.NewDecoder(r.Body).Decode(&date); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pointOfReference, err := time.Parse("01-2006", date)
	if err != nil {
		http.Error(w, "error parsing start date", http.StatusBadRequest)
		return
	}
	subs, err := s.storage.ListSubs(r.Context(), pointOfReference)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(subs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) TotalCost(w http.ResponseWriter, r *http.Request) {
	var req entity.TotalCostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	request, err := entity.TotalCostToDataBase(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	totalCost, err := s.storage.TotalCost(r.Context(), request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(totalCost); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
