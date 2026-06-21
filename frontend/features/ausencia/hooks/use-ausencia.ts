'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { ausenciaApi } from '@/lib/api/ausencia'
import type { CreateLeaveRequestPayload } from '@/types/ausencia'

export const LEAVE_REQUESTS_KEY = ['leave-requests'] as const
export const COMPENSATORY_DAYS_KEY = ['compensatory-days'] as const

export function useLeaveRequests(employeeId?: string) {
  return useQuery({
    queryKey: [...LEAVE_REQUESTS_KEY, { employeeId }],
    queryFn: () => ausenciaApi.leaveRequests.list(employeeId),
  })
}

export function useLeaveRequest(id: string) {
  return useQuery({
    queryKey: [...LEAVE_REQUESTS_KEY, id],
    queryFn: () => ausenciaApi.leaveRequests.getById(id),
    enabled: !!id,
  })
}

export function useCreateLeaveRequest() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (payload: CreateLeaveRequestPayload) =>
      ausenciaApi.leaveRequests.create(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: LEAVE_REQUESTS_KEY })
    },
  })
}

export function useApproveLeaveRequest() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => ausenciaApi.leaveRequests.approve(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: LEAVE_REQUESTS_KEY })
    },
  })
}

export function useRejectLeaveRequest() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => ausenciaApi.leaveRequests.reject(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: LEAVE_REQUESTS_KEY })
    },
  })
}

export function useCompensatoryDays(employeeId: string) {
  return useQuery({
    queryKey: [...COMPENSATORY_DAYS_KEY, employeeId],
    queryFn: () => ausenciaApi.compensatoryDays.listByEmployee(employeeId),
    enabled: !!employeeId,
  })
}

export function useCreateCompensatoryDay() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (payload: any) => ausenciaApi.compensatoryDays.create(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: COMPENSATORY_DAYS_KEY })
    },
  })
}

export function useUseCompensatoryDay() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, fechaUso }: { id: string; fechaUso: string }) =>
      ausenciaApi.compensatoryDays.use(id, fechaUso),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: COMPENSATORY_DAYS_KEY })
    },
  })
}
