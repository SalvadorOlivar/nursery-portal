package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/tuusuario/nursery-portal/internal/domain/auth"
)

func NewRouter(authHandler *AuthHandler, authMiddleware *AuthMiddleware, employeeHandler *EmployeeHandler, planificacionHandler *PlanificacionHandler, ausenciaHandler *AusenciaHandler, intercambioHandler *IntercambioHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", authHandler.Login)
			r.Post("/set-password", authHandler.SetPassword)
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireAuth)
				r.Get("/me", authHandler.Me)
				r.Post("/logout", authHandler.Logout)
			})
		})

		r.Route("/employees", func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)
			r.With(authMiddleware.RequireRoles(auth.RoleAdmin)).Post("/", employeeHandler.Create)
			r.With(authMiddleware.RequireRoles(auth.RoleAdmin, auth.RoleSupervisor)).Get("/", employeeHandler.List)
			r.Route("/{id}", func(r chi.Router) {
				r.With(authMiddleware.RequireRoles(auth.RoleAdmin, auth.RoleSupervisor, auth.RoleEmployee)).Get("/", employeeHandler.GetByID)
				r.With(authMiddleware.RequireRoles(auth.RoleAdmin)).Put("/", employeeHandler.Update)
				r.With(authMiddleware.RequireRoles(auth.RoleAdmin)).Delete("/", employeeHandler.Deactivate)
			})
		})

		r.Route("/leave-requests", func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)
			r.Post("/", ausenciaHandler.CreateLeaveRequest)
			r.Get("/", ausenciaHandler.ListLeaveRequests)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", ausenciaHandler.GetLeaveRequest)
				r.With(authMiddleware.RequireRoles(auth.RoleAdmin, auth.RoleSupervisor)).Post("/approve", ausenciaHandler.ApproveLeaveRequest)
				r.With(authMiddleware.RequireRoles(auth.RoleAdmin, auth.RoleSupervisor)).Post("/reject", ausenciaHandler.RejectLeaveRequest)
			})
		})

		r.Route("/employees/{employeeId}/compensatory-days", func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)
			r.Get("/", ausenciaHandler.ListCompensatoryDays)
		})

		r.Route("/compensatory-days", func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)
			r.With(authMiddleware.RequireRoles(auth.RoleAdmin, auth.RoleSupervisor)).Post("/", ausenciaHandler.CreateCompensatoryDay)
			r.Route("/{id}", func(r chi.Router) {
				r.With(authMiddleware.RequireRoles(auth.RoleAdmin, auth.RoleSupervisor)).Post("/use", ausenciaHandler.UseCompensatoryDay)
			})
		})

		r.Route("/planificaciones", func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)
			r.With(authMiddleware.RequireRoles(auth.RoleAdmin, auth.RoleSupervisor)).Post("/", planificacionHandler.Create)
			r.Get("/", planificacionHandler.List)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", planificacionHandler.GetByID)
				r.With(authMiddleware.RequireRoles(auth.RoleAdmin, auth.RoleSupervisor)).Put("/", planificacionHandler.Update)
				r.With(authMiddleware.RequireRoles(auth.RoleAdmin, auth.RoleSupervisor)).Delete("/", planificacionHandler.Delete)
				r.With(authMiddleware.RequireRoles(auth.RoleAdmin, auth.RoleSupervisor)).Post("/publicar", planificacionHandler.Publicar)
				r.With(authMiddleware.RequireRoles(auth.RoleAdmin, auth.RoleSupervisor)).Post("/cerrar", planificacionHandler.Cerrar)
				r.With(authMiddleware.RequireRoles(auth.RoleAdmin, auth.RoleSupervisor)).Post("/turnos", planificacionHandler.CreateTurno)
				r.With(authMiddleware.RequireRoles(auth.RoleAdmin, auth.RoleSupervisor)).Delete("/turnos/{turnoId}", planificacionHandler.DeleteTurno)
				r.Get("/requirements", planificacionHandler.GetStaffingRequirements)
			r.Get("/leaves", planificacionHandler.GetPlanLeaves)
				r.Get("/sectores", planificacionHandler.GetSectores)
				r.With(authMiddleware.RequireRoles(auth.RoleAdmin, auth.RoleSupervisor)).Put("/sectores", planificacionHandler.UpdateSectores)
				r.With(authMiddleware.RequireRoles(auth.RoleAdmin, auth.RoleSupervisor)).Put("/dotacion", planificacionHandler.UpdateDotacion)
			})
		})

		r.Route("/swap-requests", func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)
			r.Post("/", intercambioHandler.CreateSwapRequest)
			r.Get("/", intercambioHandler.ListSwapRequests)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", intercambioHandler.GetSwapRequest)
				r.Post("/accept", intercambioHandler.AcceptSwapRequest)
				r.Post("/reject", intercambioHandler.RejectSwapRequest)
				r.With(authMiddleware.RequireRoles(auth.RoleAdmin, auth.RoleSupervisor)).Post("/approve", intercambioHandler.ApproveSwapRequest)
				r.Post("/cancel", intercambioHandler.CancelSwapRequest)
				r.Get("/history", intercambioHandler.GetSwapHistory)
			})
		})
	})

	return r
}
