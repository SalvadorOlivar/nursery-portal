package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	cmd "github.com/tuusuario/nursery-portal/internal/application/commands/employee"
	"github.com/tuusuario/nursery-portal/internal/application/services"
	"github.com/tuusuario/nursery-portal/internal/domain/auth"
	"github.com/tuusuario/nursery-portal/internal/domain/employee"
)

type EmployeeHandler struct {
	svc *services.EmployeeService
}

func NewEmployeeHandler(svc *services.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{svc: svc}
}

func (h *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Nombre       string `json:"nombre"`
		Apellido     string `json:"apellido"`
		Tipo         string `json:"tipo"`
		Sector       string `json:"sector"`
		HorasMinimas int    `json:"horas_minimas"`
		HorasMaximas int    `json:"horas_maximas"`
		WorkDays     *int   `json:"work_days,omitempty"`
		RestDays     *int   `json:"rest_days,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	emp, initialPassword, err := h.svc.Create(r.Context(), cmd.CreateEmployeeCommand{
		Nombre:       req.Nombre,
		Apellido:     req.Apellido,
		Tipo:         req.Tipo,
		Sector:       req.Sector,
		HorasMinimas: req.HorasMinimas,
		HorasMaximas: req.HorasMaximas,
		WorkDays:     req.WorkDays,
		RestDays:     req.RestDays,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := toResponse(emp)
	resp.InitialPassword = initialPassword
	writeJSON(w, http.StatusCreated, resp)
}

func (h *EmployeeHandler) List(w http.ResponseWriter, r *http.Request) {
	employees, err := h.svc.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	items := make([]employeeResponse, len(employees))
	for i, emp := range employees {
		items[i] = toResponse(emp)
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": items})
}

func (h *EmployeeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, ok := authUserFromContext(r.Context())
	if ok && user.Role == auth.RoleEmployee {
		if user.EmployeeID == nil || *user.EmployeeID != id {
			writeError(w, http.StatusForbidden, "forbidden")
			return
		}
	}

	emp, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "employee not found")
		return
	}

	writeJSON(w, http.StatusOK, toResponse(emp))
}

func (h *EmployeeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req struct {
		Nombre       string `json:"nombre"`
		Apellido     string `json:"apellido"`
		Tipo         string `json:"tipo"`
		Sector       string `json:"sector"`
		HorasMinimas int    `json:"horas_minimas"`
		HorasMaximas int    `json:"horas_maximas"`
		WorkDays     *int   `json:"work_days,omitempty"`
		RestDays     *int   `json:"rest_days,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	emp, err := h.svc.Update(r.Context(), cmd.UpdateEmployeeCommand{
		ID:           id,
		Nombre:       req.Nombre,
		Apellido:     req.Apellido,
		Tipo:         req.Tipo,
		Sector:       req.Sector,
		HorasMinimas: req.HorasMinimas,
		HorasMaximas: req.HorasMaximas,
		WorkDays:     req.WorkDays,
		RestDays:     req.RestDays,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, toResponse(emp))
}

func (h *EmployeeHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.Deactivate(r.Context(), id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type employeeResponse struct {
	ID              string `json:"id"`
	Nombre          string `json:"nombre"`
	Apellido        string `json:"apellido"`
	Tipo            string `json:"tipo"`
	Sector          string `json:"sector"`
	HorasMinimas    int    `json:"horas_minimas"`
	HorasMaximas    int    `json:"horas_maximas"`
	WorkDays        int    `json:"work_days"`
	RestDays        int    `json:"rest_days"`
	Activo          bool   `json:"activo"`
	InitialPassword string `json:"initial_password,omitempty"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

func toResponse(e *employee.Employee) employeeResponse {
	return employeeResponse{
		ID:            e.ID,
		Nombre:        e.Nombre,
		Apellido:      e.Apellido,
		Tipo:          string(e.Tipo),
		Sector:        e.Sector,
		HorasMinimas:  e.HorasMinimas,
		HorasMaximas:  e.HorasMaximas,
		WorkDays:      e.PatronTrabajo.WorkDays,
		RestDays:      e.PatronTrabajo.RestDays,
		Activo:        e.Activo,
		CreatedAt:     e.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     e.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}


