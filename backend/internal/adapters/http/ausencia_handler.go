package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	cmdcomp "github.com/tuusuario/nursery-portal/internal/application/commands/compensatory"
	cmdleave "github.com/tuusuario/nursery-portal/internal/application/commands/leave"
	"github.com/tuusuario/nursery-portal/internal/application/services"
	"github.com/tuusuario/nursery-portal/internal/domain/ausencia"
)

type AusenciaHandler struct {
	svc *services.AusenciaService
}

func NewAusenciaHandler(svc *services.AusenciaService) *AusenciaHandler {
	return &AusenciaHandler{svc: svc}
}

func (h *AusenciaHandler) CreateLeaveRequest(w http.ResponseWriter, r *http.Request) {
	var req struct {
		EmployeeID  string `json:"employee_id"`
		FechaInicio string `json:"fecha_inicio"`
		FechaFin    string `json:"fecha_fin"`
		Tipo        string `json:"tipo"`
		Motivo      string `json:"motivo"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	fechaInicio, err := time.Parse("2006-01-02", req.FechaInicio)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid fecha_inicio format (expected YYYY-MM-DD)")
		return
	}
	fechaFin, err := time.Parse("2006-01-02", req.FechaFin)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid fecha_fin format (expected YYYY-MM-DD)")
		return
	}

	lr, err := h.svc.CreateLeaveRequest(r.Context(), cmdleave.CreateLeaveRequestCommand{
		EmployeeID:  req.EmployeeID,
		FechaInicio: fechaInicio,
		FechaFin:    fechaFin,
		Tipo:        req.Tipo,
		Motivo:      req.Motivo,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, toLeaveRequestResponse(lr))
}

func (h *AusenciaHandler) ListLeaveRequests(w http.ResponseWriter, r *http.Request) {
	employeeID := r.URL.Query().Get("employee_id")
	requests, err := h.svc.ListLeaveRequests(r.Context(), employeeID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	items := make([]leaveRequestResponse, len(requests))
	for i, lr := range requests {
		items[i] = toLeaveRequestResponse(lr)
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": items})
}

func (h *AusenciaHandler) GetLeaveRequest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	lr, err := h.svc.GetLeaveRequest(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "leave request not found")
		return
	}

	writeJSON(w, http.StatusOK, toLeaveRequestResponse(lr))
}

func (h *AusenciaHandler) ApproveLeaveRequest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, ok := authUserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.svc.ApproveLeaveRequest(r.Context(), id, user.ID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AusenciaHandler) RejectLeaveRequest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, ok := authUserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.svc.RejectLeaveRequest(r.Context(), id, user.ID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AusenciaHandler) CreateCompensatoryDay(w http.ResponseWriter, r *http.Request) {
	var req struct {
		EmployeeID  string  `json:"employee_id"`
		FechaOrigen string  `json:"fecha_origen"`
		Motivo      string  `json:"motivo"`
		TurnoID     *string `json:"turno_id,omitempty"`
		Descripcion string  `json:"descripcion"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	fechaOrigen, err := time.Parse("2006-01-02", req.FechaOrigen)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid fecha_origen format (expected YYYY-MM-DD)")
		return
	}

	cd, err := h.svc.CreateCompensatoryDay(r.Context(), cmdcomp.CreateCompensatoryDayCommand{
		EmployeeID:  req.EmployeeID,
		FechaOrigen: fechaOrigen,
		Motivo:      req.Motivo,
		TurnoID:     req.TurnoID,
		Descripcion: req.Descripcion,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, toCompensatoryDayResponse(cd))
}

func (h *AusenciaHandler) UseCompensatoryDay(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		FechaUso string `json:"fecha_uso"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	fechaUso, err := time.Parse("2006-01-02", req.FechaUso)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid fecha_uso format (expected YYYY-MM-DD)")
		return
	}

	if err := h.svc.UseCompensatoryDay(r.Context(), id, fechaUso); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AusenciaHandler) ListCompensatoryDays(w http.ResponseWriter, r *http.Request) {
	employeeID := chi.URLParam(r, "employeeId")
	result, err := h.svc.ListCompensatoryDays(r.Context(), employeeID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	items := make([]compensatoryDayResponse, len(result.Items))
	for i, cd := range result.Items {
		items[i] = toCompensatoryDayResponse(cd)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data":           items,
		"available_count": result.AvailableCount,
	})
}

type leaveRequestResponse struct {
	ID          string  `json:"id"`
	EmployeeID  string  `json:"employee_id"`
	FechaInicio string  `json:"fecha_inicio"`
	FechaFin    string  `json:"fecha_fin"`
	Tipo        string  `json:"tipo"`
	Estado      string  `json:"estado"`
	Motivo      string  `json:"motivo"`
	AprobadoPor *string `json:"aprobado_por,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type compensatoryDayResponse struct {
	ID          string  `json:"id"`
	EmployeeID  string  `json:"employee_id"`
	FechaOrigen string  `json:"fecha_origen"`
	Motivo      string  `json:"motivo"`
	TurnoID     *string `json:"turno_id,omitempty"`
	Descripcion string  `json:"descripcion"`
	Utilizado   bool    `json:"utilizado"`
	FechaUso    *string `json:"fecha_uso,omitempty"`
	CreatedAt   string  `json:"created_at"`
}

func toLeaveRequestResponse(lr *ausencia.LeaveRequest) leaveRequestResponse {
	r := leaveRequestResponse{
		ID:          lr.ID,
		EmployeeID:  lr.EmployeeID,
		FechaInicio: lr.FechaInicio.Format("2006-01-02"),
		FechaFin:    lr.FechaFin.Format("2006-01-02"),
		Tipo:        string(lr.Tipo),
		Estado:      string(lr.Estado),
		Motivo:      lr.Motivo,
		AprobadoPor: lr.AprobadoPor,
		CreatedAt:   lr.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   lr.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if lr.AprobadoPor != nil {
		r.AprobadoPor = lr.AprobadoPor
	}
	return r
}

func toCompensatoryDayResponse(cd *ausencia.CompensatoryDay) compensatoryDayResponse {
	r := compensatoryDayResponse{
		ID:          cd.ID,
		EmployeeID:  cd.EmployeeID,
		FechaOrigen: cd.FechaOrigen.Format("2006-01-02"),
		Motivo:      string(cd.Motivo),
		TurnoID:     cd.TurnoID,
		Descripcion: cd.Descripcion,
		Utilizado:   cd.Utilizado,
		CreatedAt:   cd.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if cd.FechaUso != nil {
		s := cd.FechaUso.Format("2006-01-02")
		r.FechaUso = &s
	}
	return r
}
