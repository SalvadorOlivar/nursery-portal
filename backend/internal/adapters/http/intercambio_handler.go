package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	cmdinter "github.com/tuusuario/nurse-portal/internal/application/commands/intercambio"
	"github.com/tuusuario/nurse-portal/internal/application/services"
	"github.com/tuusuario/nurse-portal/internal/domain/intercambio"
)

type IntercambioHandler struct {
	svc *services.IntercambioService
}

func NewIntercambioHandler(svc *services.IntercambioService) *IntercambioHandler {
	return &IntercambioHandler{svc: svc}
}

func (h *IntercambioHandler) CreateSwapRequest(w http.ResponseWriter, r *http.Request) {
	user, ok := authUserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if user.EmployeeID == nil {
		writeError(w, http.StatusForbidden, "user must be linked to an employee")
		return
	}

	var req struct {
		PlanificacionID    string `json:"planificacion_id"`
		TurnoSolicitanteID string `json:"turno_solicitante_id"`
		TurnoDestinoID     string `json:"turno_destino_id"`
		DestinoID          string `json:"destino_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	swapReq, err := h.svc.CreateSwapRequest(r.Context(), cmdinter.CreateSwapRequestCommand{
		PlanificacionID:    req.PlanificacionID,
		TurnoSolicitanteID: req.TurnoSolicitanteID,
		TurnoDestinoID:     req.TurnoDestinoID,
		SolicitanteID:      *user.EmployeeID,
		DestinoID:          req.DestinoID,
		ActorID:            user.ID,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, toSwapRequestResponse(swapReq))
}

func (h *IntercambioHandler) ListSwapRequests(w http.ResponseWriter, r *http.Request) {
	user, ok := authUserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var employeeID string
	if user.EmployeeID != nil {
		employeeID = *user.EmployeeID
	}

	requests, err := h.svc.ListSwapRequests(r.Context(), employeeID, string(user.Role))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	items := make([]swapRequestResponse, len(requests))
	for i, req := range requests {
		items[i] = toSwapRequestResponse(req)
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": items})
}

func (h *IntercambioHandler) GetSwapRequest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	req, err := h.svc.GetSwapRequest(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "swap request not found")
		return
	}

	writeJSON(w, http.StatusOK, toSwapRequestResponse(req))
}

func (h *IntercambioHandler) AcceptSwapRequest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, ok := authUserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if user.EmployeeID == nil {
		writeError(w, http.StatusForbidden, "user must be linked to an employee")
		return
	}

	if err := h.svc.AcceptSwapRequest(r.Context(), id, user.ID, *user.EmployeeID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *IntercambioHandler) RejectSwapRequest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, ok := authUserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if user.EmployeeID == nil {
		writeError(w, http.StatusForbidden, "user must be linked to an employee")
		return
	}

	if err := h.svc.RejectSwapRequest(r.Context(), id, user.ID, *user.EmployeeID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *IntercambioHandler) ApproveSwapRequest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, ok := authUserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.svc.ApproveSwapRequest(r.Context(), id, user.ID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *IntercambioHandler) CancelSwapRequest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, ok := authUserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if user.EmployeeID == nil {
		writeError(w, http.StatusForbidden, "user must be linked to an employee")
		return
	}

	if err := h.svc.CancelSwapRequest(r.Context(), id, user.ID, *user.EmployeeID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *IntercambioHandler) GetSwapHistory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	entries, err := h.svc.GetSwapHistory(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "swap request not found")
		return
	}

	items := make([]swapHistoryResponse, len(entries))
	for i, e := range entries {
		items[i] = toSwapHistoryResponse(e)
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": items})
}

type swapRequestResponse struct {
	ID                 string  `json:"id"`
	PlanificacionID    string  `json:"planificacion_id"`
	TurnoSolicitanteID string  `json:"turno_solicitante_id"`
	TurnoDestinoID     string  `json:"turno_destino_id"`
	SolicitanteID      string  `json:"solicitante_id"`
	DestinoID          string  `json:"destino_id"`
	Estado             string  `json:"estado"`
	AprobadoPor        *string `json:"aprobado_por,omitempty"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
}

func toSwapRequestResponse(req *intercambio.ShiftSwapRequest) swapRequestResponse {
	r := swapRequestResponse{
		ID:                 req.ID,
		PlanificacionID:    req.PlanificacionID,
		TurnoSolicitanteID: req.TurnoSolicitanteID,
		TurnoDestinoID:     req.TurnoDestinoID,
		SolicitanteID:      req.SolicitanteID,
		DestinoID:          req.DestinoID,
		Estado:             string(req.Estado),
		CreatedAt:          req.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:          req.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if req.AprobadoPor != nil {
		r.AprobadoPor = req.AprobadoPor
	}
	return r
}

type swapHistoryResponse struct {
	ID            string  `json:"id"`
	SwapRequestID string  `json:"swap_request_id"`
	Accion        string  `json:"accion"`
	ActorID       string  `json:"actor_id"`
	Detalle       *string `json:"detalle,omitempty"`
	CreatedAt     string  `json:"created_at"`
}

func toSwapHistoryResponse(e *intercambio.ShiftSwapHistoryEntry) swapHistoryResponse {
	r := swapHistoryResponse{
		ID:            e.ID,
		SwapRequestID: e.SwapRequestID,
		Accion:        string(e.Accion),
		ActorID:       e.ActorID,
		Detalle:       e.Detalle,
		CreatedAt:     e.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	return r
}
