import { api } from './client'
import type { ShiftSwapRequest, CreateSwapRequestPayload, SwapHistoryEntry } from '@/types/intercambio'

export const intercambioApi = {
  list: () =>
    api.get<{ data: ShiftSwapRequest[] }>('/swap-requests'),

  getById: (id: string) =>
    api.get<ShiftSwapRequest>(`/swap-requests/${id}`),

  create: (payload: CreateSwapRequestPayload) =>
    api.post<ShiftSwapRequest>('/swap-requests', payload),

  accept: (id: string) =>
    api.post<void>(`/swap-requests/${id}/accept`, {}),

  reject: (id: string) =>
    api.post<void>(`/swap-requests/${id}/reject`, {}),

  approve: (id: string) =>
    api.post<void>(`/swap-requests/${id}/approve`, {}),

  cancel: (id: string) =>
    api.post<void>(`/swap-requests/${id}/cancel`, {}),

  history: (id: string) =>
    api.get<{ data: SwapHistoryEntry[] }>(`/swap-requests/${id}/history`),
}
