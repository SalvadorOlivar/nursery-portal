'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { intercambioApi } from '@/lib/api/intercambio'
import type { CreateSwapRequestPayload } from '@/types/intercambio'

export const SWAP_REQUESTS_KEY = ['swap-requests'] as const
export const SWAP_HISTORY_KEY = ['swap-history'] as const

export function useSwapRequests() {
  return useQuery({
    queryKey: SWAP_REQUESTS_KEY,
    queryFn: () => intercambioApi.list(),
  })
}

export function useSwapRequest(id: string) {
  return useQuery({
    queryKey: [...SWAP_REQUESTS_KEY, id],
    queryFn: () => intercambioApi.getById(id),
    enabled: !!id,
  })
}

export function useCreateSwapRequest() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (payload: CreateSwapRequestPayload) =>
      intercambioApi.create(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: SWAP_REQUESTS_KEY })
    },
  })
}

export function useAcceptSwapRequest() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => intercambioApi.accept(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: SWAP_REQUESTS_KEY })
    },
  })
}

export function useRejectSwapRequest() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => intercambioApi.reject(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: SWAP_REQUESTS_KEY })
    },
  })
}

export function useApproveSwapRequest() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => intercambioApi.approve(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: SWAP_REQUESTS_KEY })
    },
  })
}

export function useCancelSwapRequest() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => intercambioApi.cancel(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: SWAP_REQUESTS_KEY })
    },
  })
}

export function useSwapHistory(id: string) {
  return useQuery({
    queryKey: [...SWAP_HISTORY_KEY, id],
    queryFn: () => intercambioApi.history(id),
    enabled: !!id,
  })
}
