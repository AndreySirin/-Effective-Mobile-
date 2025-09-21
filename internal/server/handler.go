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
	lg := s.lg.With("handler", "CreateSubs")
	lg.Info("received create subscription request")

	var req entity.SubsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lg.Error("failed to decode request body", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	subs, err := entity.SubsToDataBase(lg, req)
	if err != nil {
		lg.Error("failed to convert request to subscription entity", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := s.storage.CreateSubs(r.Context(), &subs)
	if err != nil {
		lg.Error("failed to create subscription in storage",
			"user_id", subs.UserId,
			"service_name", subs.ServiceName,
			"err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	endDateStr := "nil"
	if subs.EndDate != nil {
		endDateStr = subs.EndDate.Format("2006-01")
	}

	lg.Info("subscription created successfully",
		"id", id,
		"user_id", subs.UserId,
		"service_name", subs.ServiceName,
		"price", subs.Price,
		"start_date", subs.StartDate.Format("2006-01"),
		"end_date", endDateStr,
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := map[string]string{"id": id.String()}
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		lg.Error("failed to encode response", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) ReadSubs(w http.ResponseWriter, r *http.Request) {
	lg := s.lg.With("handler", "ReadSubs")

	subsID := chi.URLParam(r, "id")
	lg.Info("received read subscription request", "id", subsID)

	id, err := uuid.Parse(subsID)
	if err != nil {
		lg.Error("failed to parse subscription id", "id", subsID, "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if id == uuid.Nil {
		lg.Warn("subscription id is nil", "id", subsID)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	subs, err := s.storage.ReadSubs(r.Context(), id)
	if errors.Is(err, storage.ErrNotFound) {
		lg.Info("subscription not found", "id", id)
		http.Error(w, "subscription not found", http.StatusNotFound)
		return
	} else if err != nil {
		lg.Error("failed to read subscription from storage", "id", id, "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	lg.Info("subscription retrieved successfully",
		"id", subs.SubsID,
		"user_id", subs.UserId,
		"service_name", subs.ServiceName,
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(subs); err != nil {
		lg.Error("failed to encode response", "id", subs.SubsID, "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) UpdateSubs(w http.ResponseWriter, r *http.Request) {
	lg := s.lg.With("handler", "UpdateSubs")
	lg.Info("received update subscription request")

	var req entity.SubsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lg.Error("failed to decode request body", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	subs, err := entity.SubsToDataBase(lg, req)
	if err != nil {
		lg.Error("failed to convert request to subscription entity", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	subsID := chi.URLParam(r, "id")
	id, err := uuid.Parse(subsID)
	if err != nil {
		lg.Error("failed to parse subscription id", "id", subsID, "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if id == uuid.Nil {
		lg.Warn("subscription id is nil", "id", subsID)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = s.storage.UpdateSubs(r.Context(), id, &subs)
	if errors.Is(err, storage.ErrNotFound) {
		lg.Info("subscription not found in storage", "id", id)
		http.Error(w, "subscription not found", http.StatusNotFound)
		return
	} else if err != nil {
		lg.Error("failed to update subscription in storage", "id", id, "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// безопасная обработка EndDate
	endDateStr := "nil"
	if subs.EndDate != nil {
		endDateStr = subs.EndDate.Format("2006-01")
	}

	lg.Info("subscription updated successfully",
		"id", id,
		"user_id", subs.UserId,
		"service_name", subs.ServiceName,
		"price", subs.Price,
		"start_date", subs.StartDate.Format("2006-01"),
		"end_date", endDateStr,
	)

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) DeleteSubs(w http.ResponseWriter, r *http.Request) {
	lg := s.lg.With("handler", "DeleteSubs")

	subsID := chi.URLParam(r, "id")
	lg.Info("received delete subscription request", "id", subsID)

	id, err := uuid.Parse(subsID)
	if err != nil {
		lg.Error("failed to parse subscription id", "id", subsID, "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.storage.DeleteSubs(r.Context(), id)
	if errors.Is(err, storage.ErrNotFound) {
		lg.Info("subscription not found in storage", "id", id)
		http.Error(w, "subscription not found", http.StatusNotFound)
		return
	} else if err != nil {
		lg.Error("failed to delete subscription from storage", "id", id, "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	lg.Info("subscription deleted successfully", "id", id)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) ListSubs(w http.ResponseWriter, r *http.Request) {
	lg := s.lg.With("handler", "ListSubs")
	lg.Info("received list subscriptions request")

	date := r.URL.Query().Get("point_of_reference")
	if date == "" {
		lg.Warn("missing query parameter", "point_of_reference", date)
		http.Error(w, "missing point_of_reference query parameter", http.StatusBadRequest)
		return
	}

	pointOfReference, err := time.Parse("01-2006", date)
	if err != nil {
		lg.Error("failed to parse date", "date", date, "err", err)
		http.Error(w, "error parsing start date", http.StatusBadRequest)
		return
	}

	subs, err := s.storage.ListSubs(r.Context(), pointOfReference)
	if err != nil {
		lg.Error("failed to list subscriptions from storage", "date", pointOfReference, "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	lg.Info("subscriptions retrieved successfully", "count", len(subs), "date", pointOfReference)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(subs); err != nil {
		lg.Error("failed to encode response", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) TotalCost(w http.ResponseWriter, r *http.Request) {
	lg := s.lg.With("handler", "TotalCost")
	lg.Info("received total cost request")

	var req entity.TotalCostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lg.Error("failed to decode request body", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	request, err := entity.TotalCostToDataBase(lg, req)
	if err != nil {
		lg.Error("failed to convert request to database model", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	totalCost, err := s.storage.TotalCost(r.Context(), request)
	if err != nil {
		lg.Error("failed to calculate total cost from storage", "user_id", request.UserId, "service_name", request.ServiceName, "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lg.Info("total cost calculated successfully",
		"user_id", request.UserId,
		"service_name", request.ServiceName,
		"date1", request.Date1.Format("2006-01"),
		"date2", request.Date2.Format("2006-01"),
		"total_cost", totalCost,
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(totalCost); err != nil {
		lg.Error("failed to encode response", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
