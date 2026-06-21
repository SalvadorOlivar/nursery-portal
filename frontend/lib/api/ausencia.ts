import { api } from './client'
import type { LeaveRequest, CreateLeaveRequestPayload, CompensatoryDaysResponse, CreateCompensatoryDayPayload } from '@/types/ausencia'

export const ausenciaApi = {
  leaveRequests: {
    list: (employeeId?: string) =>
      api.get<{ data: LeaveRequest[] }>(`/leave-requests${employeeId ? `?employee_id=${employeeId}` : ''}`),

    getById: (id: string) =>
      api.get<LeaveRequest>(`/leave-requests/${id}`),

    create: (payload: CreateLeaveRequestPayload) =>
      api.post<LeaveRequest>('/leave-requests', payload),

    approve: (id: string) =>
      api.post<void>(`/leave-requests/${id}/approve`, {}),

    reject: (id: string) =>
      api.post<void>(`/leave-requests/${id}/reject`, {}),
  },

  compensatoryDays: {
    listByEmployee: (employeeId: string) =>
      api.get<CompensatoryDaysResponse>(`/employees/${employeeId}/compensatory-days`),

    create: (payload: CreateCompensatoryDayPayload) =>
      api.post<CompensatoryDayResponse>('/compensatory-days', payload),

    use: (id: string, fechaUso: string) =>
      api.post<void>(`/compensatory-days/${id}/use`, { fecha_uso: fechaUso }),
  },
}

// Local type for create response
interface CompensatoryDayResponse {
  id: string
  employee_id: string
  fecha_origen: string
  motivo: string
  turno_id?: string
  descripcion: string
  utilizado: boolean
  created_at: string
}
