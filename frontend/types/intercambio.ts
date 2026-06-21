export type SwapRequestStatus = 'PENDIENTE_RESPUESTA' | 'PENDIENTE_APROBACION' | 'APROBADO' | 'RECHAZADO' | 'CANCELADO'

export interface ShiftSwapRequest {
  id: string
  planificacion_id: string
  turno_solicitante_id: string
  turno_destino_id: string
  solicitante_id: string
  destino_id: string
  estado: SwapRequestStatus
  aprobado_por?: string
  created_at: string
  updated_at: string
}

export interface CreateSwapRequestPayload {
  planificacion_id: string
  turno_solicitante_id: string
  turno_destino_id: string
  destino_id: string
}

export interface SwapHistoryEntry {
  id: string
  swap_request_id: string
  accion: string
  actor_id: string
  detalle?: string
  created_at: string
}
