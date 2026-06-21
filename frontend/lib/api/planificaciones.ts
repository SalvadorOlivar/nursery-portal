import { api } from './client'
import type {
  Planificacion,
  PlanificacionDetail,
  CreatePlanificacionPayload,
  CreateTurnoPayload,
  Turno,
  DotacionItem,
  SectorItem,
  UpdateSectoresPayload,
  UpdateDotacionPayload,
} from '@/types/planificacion'

export const planificacionesApi = {
  list: () => api.get<{ data: Planificacion[] }>('/planificaciones'),

  getById: (id: string) => api.get<PlanificacionDetail>(`/planificaciones/${id}`),

  create: (payload: CreatePlanificacionPayload) =>
    api.post<Planificacion>('/planificaciones', payload),

  update: (id: string, nombre: string) =>
    api.put<void>(`/planificaciones/${id}`, { nombre }),

  delete: (id: string) => api.delete<void>(`/planificaciones/${id}`),

  publicar: (id: string) => api.post<void>(`/planificaciones/${id}/publicar`, {}),

  cerrar: (id: string) => api.post<void>(`/planificaciones/${id}/cerrar`, {}),

  createTurno: (planificacionId: string, payload: CreateTurnoPayload) =>
    api.post<Turno>(`/planificaciones/${planificacionId}/turnos`, payload),

  deleteTurno: (planificacionId: string, turnoId: string) =>
    api.delete<void>(`/planificaciones/${planificacionId}/turnos/${turnoId}`),

  getDotacion: (id: string) =>
    api.get<{ data: DotacionItem[] }>(`/planificaciones/${id}/requirements`),

  getSectores: (id: string) =>
    api.get<{ data: SectorItem[] }>(`/planificaciones/${id}/sectores`),

  getLeaves: (id: string) =>
    api.get<{ data: import('@/types/ausencia').LeaveRequest[] }>(`/planificaciones/${id}/leaves`),

  updateSectores: (id: string, payload: UpdateSectoresPayload) =>
    api.put<void>(`/planificaciones/${id}/sectores`, payload),

  updateDotacion: (id: string, payload: UpdateDotacionPayload) =>
    api.put<void>(`/planificaciones/${id}/dotacion`, payload),
}
