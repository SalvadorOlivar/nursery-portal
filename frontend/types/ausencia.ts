export type LeaveType = 'VACACIONES' | 'ENFERMEDAD' | 'PERSONAL' | 'DIA_FAVOR'
export type LeaveStatus = 'PENDIENTE' | 'APROBADO' | 'RECHAZADO'
export type CompensatoryMotivo = 'DOBLE_TURNO' | 'DESCANSO_LABORADO'

export interface LeaveRequest {
  id: string
  employee_id: string
  fecha_inicio: string
  fecha_fin: string
  tipo: LeaveType
  estado: LeaveStatus
  motivo: string
  aprobado_por?: string
  created_at: string
  updated_at: string
}

export interface CreateLeaveRequestPayload {
  employee_id: string
  fecha_inicio: string
  fecha_fin: string
  tipo: LeaveType
  motivo: string
}

export interface CompensatoryDay {
  id: string
  employee_id: string
  fecha_origen: string
  motivo: CompensatoryMotivo
  turno_id?: string
  descripcion: string
  utilizado: boolean
  fecha_uso?: string
  created_at: string
}

export interface CompensatoryDaysResponse {
  data: CompensatoryDay[]
  available_count: number
}

export interface CreateCompensatoryDayPayload {
  employee_id: string
  fecha_origen: string
  motivo: CompensatoryMotivo
  turno_id?: string
  descripcion: string
}
