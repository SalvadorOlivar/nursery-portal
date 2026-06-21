package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	cmdplanif "github.com/tuusuario/nursery-portal/internal/application/commands/planificacion"
	cmdturno "github.com/tuusuario/nursery-portal/internal/application/commands/turno"
	qryplanif "github.com/tuusuario/nursery-portal/internal/application/queries/planificacion"
	"github.com/tuusuario/nursery-portal/internal/application/services"
	"github.com/tuusuario/nursery-portal/internal/domain/employee"
	"github.com/tuusuario/nursery-portal/internal/domain/planificacion"
	"github.com/tuusuario/nursery-portal/internal/domain/turno"
)

type PlanificacionHandler struct {
	svc         *services.PlanificacionService
	employeeSvc *services.EmployeeService
}

func NewPlanificacionHandler(svc *services.PlanificacionService, employeeSvc *services.EmployeeService) *PlanificacionHandler {
	return &PlanificacionHandler{svc: svc, employeeSvc: employeeSvc}
}

func (h *PlanificacionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Semana int    `json:"semana"`
		Anio   int    `json:"anio"`
		Nombre string `json:"nombre"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	p, err := h.svc.Create(r.Context(), cmdplanif.CreatePlanificacionCommand{
		Semana: req.Semana,
		Anio:   req.Anio,
		Nombre: req.Nombre,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, toPlanificacionResponse(p))
}

func (h *PlanificacionHandler) List(w http.ResponseWriter, r *http.Request) {
	plans, err := h.svc.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	items := make([]planificacionResponse, len(plans))
	for i, p := range plans {
		items[i] = toPlanificacionResponse(p)
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": items})
}

func (h *PlanificacionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	result, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "planificacion not found")
		return
	}

	employees, err := h.employeeSvc.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, toPlanificacionDetailResponse(result, employees))
}

func (h *PlanificacionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req struct {
		Nombre string `json:"nombre"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.svc.Update(r.Context(), cmdplanif.UpdatePlanificacionCommand{
		ID:     id,
		Nombre: req.Nombre,
	}); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PlanificacionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PlanificacionHandler) Publicar(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.Publicar(r.Context(), id); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PlanificacionHandler) Cerrar(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.Cerrar(r.Context(), id); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PlanificacionHandler) CreateTurno(w http.ResponseWriter, r *http.Request) {
	planificacionID := chi.URLParam(r, "id")

	var req struct {
		EmpleadoID string `json:"empleado_id"`
		DiaSemana  int    `json:"dia_semana"`
		Tipo       string `json:"tipo"`
		Sector     string `json:"sector"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	t, err := h.svc.CreateTurno(r.Context(), cmdturno.CreateTurnoCommand{
		PlanificacionID: planificacionID,
		EmpleadoID:      req.EmpleadoID,
		DiaSemana:       req.DiaSemana,
		Tipo:            req.Tipo,
		Sector:          req.Sector,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, toTurnoResponse(t))
}

func (h *PlanificacionHandler) DeleteTurno(w http.ResponseWriter, r *http.Request) {
	turnoID := chi.URLParam(r, "turnoId")
	if err := h.svc.DeleteTurno(r.Context(), turnoID); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PlanificacionHandler) GetStaffingRequirements(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	dotacion, err := h.svc.GetStaffingRequirements(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "planificacion not found")
		return
	}

	items := make([]dotacionItemResponse, len(dotacion))
	for i, d := range dotacion {
		items[i] = toDotacionItemResponse(d)
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": items})
}

func (h *PlanificacionHandler) GetPlanLeaves(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	leaves, err := h.svc.GetPlanLeaves(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "planificacion not found")
		return
	}

	items := make([]leaveRequestResponse, len(leaves))
	for i, lr := range leaves {
		items[i] = toLeaveRequestResponse(lr)
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": items})
}

func (h *PlanificacionHandler) GetSectores(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sectores, err := h.svc.GetSectores(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "planificacion not found")
		return
	}

	items := make([]sectorResponse, len(sectores))
	for i, s := range sectores {
		items[i] = sectorResponse{ID: s.ID, Nombre: s.Nombre}
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": items})
}

func (h *PlanificacionHandler) UpdateSectores(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req struct {
		Sectores []string `json:"sectores"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.svc.UpdateSectores(r.Context(), cmdplanif.UpdateSectoresCommand{
		PlanificacionID: id,
		Sectores:        req.Sectores,
	}); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PlanificacionHandler) UpdateDotacion(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req struct {
		Items []dotacionUpdateInput `json:"items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	items := make([]cmdplanif.DotacionItemInput, len(req.Items))
	for i, in := range req.Items {
		items[i] = cmdplanif.DotacionItemInput{
			Sector:         in.Sector,
			TipoEmpleado:   in.TipoEmpleado,
			Turno:          in.Turno,
			CantidadMinima: in.CantidadMinima,
		}
	}

	if err := h.svc.UpdateDotacion(r.Context(), cmdplanif.UpdateDotacionCommand{
		PlanificacionID: id,
		Items:           items,
	}); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type planificacionResponse struct {
	ID        string `json:"id"`
	Semana    int    `json:"semana"`
	Anio      int    `json:"anio"`
	Nombre    string `json:"nombre"`
	Estado    string `json:"estado"`
	Dias      int    `json:"dias"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type turnoResponse struct {
	ID              string `json:"id"`
	PlanificacionID string `json:"planificacion_id"`
	EmpleadoID      string `json:"empleado_id"`
	DiaSemana       int    `json:"dia_semana"`
	Tipo            string `json:"tipo"`
	Sector          string `json:"sector"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type planificacionDetailResponse struct {
	planificacionResponse
	Turnos    []turnoResponse    `json:"turnos"`
	Employees []employeeResponse `json:"employees"`
}

func toPlanificacionResponse(p *planificacion.Planificacion) planificacionResponse {
	return planificacionResponse{
		ID:        p.ID,
		Semana:    p.Semana,
		Anio:      p.Anio,
		Nombre:    p.Nombre,
		Estado:    string(p.Estado),
		Dias:      p.Dias(),
		CreatedAt: p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toTurnoResponse(t *turno.Turno) turnoResponse {
	return turnoResponse{
		ID:              t.ID,
		PlanificacionID: t.PlanificacionID,
		EmpleadoID:      t.EmpleadoID,
		DiaSemana:       t.DiaSemana,
		Tipo:            string(t.Tipo),
		Sector:          t.Sector,
		CreatedAt:       t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       t.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toPlanificacionDetailResponse(detail *qryplanif.PlanificacionConTurnos, employees []*employee.Employee) planificacionDetailResponse {
	turnos := make([]turnoResponse, len(detail.Turnos))
	for i, t := range detail.Turnos {
		turnos[i] = toTurnoResponse(t)
	}
	employeeItems := make([]employeeResponse, len(employees))
	for i, emp := range employees {
		employeeItems[i] = toResponse(emp)
	}

	return planificacionDetailResponse{
		planificacionResponse: toPlanificacionResponse(detail.Planificacion),
		Turnos:                turnos,
		Employees:             employeeItems,
	}
}

type dotacionItemResponse struct {
	TipoEmpleado   string `json:"tipo_empleado"`
	Turno          string `json:"turno"`
	CantidadMinima int    `json:"cantidad_minima"`
	Sector         string `json:"sector"`
}

func toDotacionItemResponse(d planificacion.DotacionItem) dotacionItemResponse {
	return dotacionItemResponse{
		TipoEmpleado:   string(d.TipoEmpleado),
		Turno:          d.Turno,
		CantidadMinima: d.CantidadMinima,
		Sector:         d.Sector,
	}
}

type sectorResponse struct {
	ID     string `json:"id"`
	Nombre string `json:"nombre"`
}

type dotacionUpdateInput struct {
	Sector         string `json:"sector"`
	TipoEmpleado   string `json:"tipo_empleado"`
	Turno          string `json:"turno"`
	CantidadMinima int    `json:"cantidad_minima"`
}
