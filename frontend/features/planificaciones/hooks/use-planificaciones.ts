'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { planificacionesApi } from '@/lib/api/planificaciones'
import type { CreatePlanificacionPayload, CreateTurnoPayload } from '@/types/planificacion'

const PLANIFICACIONES_KEY = ['planificaciones'] as const

export function usePlanificaciones() {
  return useQuery({
    queryKey: PLANIFICACIONES_KEY,
    queryFn: () => planificacionesApi.list(),
  })
}

export function usePlanificacion(id: string) {
  return useQuery({
    queryKey: [...PLANIFICACIONES_KEY, id],
    queryFn: () => planificacionesApi.getById(id),
    enabled: !!id,
  })
}

export function useCreatePlanificacion() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (payload: CreatePlanificacionPayload) => planificacionesApi.create(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: PLANIFICACIONES_KEY })
    },
  })
}

export function useUpdatePlanificacion() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, nombre }: { id: string; nombre: string }) => planificacionesApi.update(id, nombre),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: PLANIFICACIONES_KEY })
      queryClient.invalidateQueries({ queryKey: [...PLANIFICACIONES_KEY, variables.id] })
    },
  })
}

export function useDeletePlanificacion() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => planificacionesApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: PLANIFICACIONES_KEY })
    },
  })
}

export function usePublicarPlanificacion() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => planificacionesApi.publicar(id),
    onSuccess: (_data, id) => {
      queryClient.invalidateQueries({ queryKey: PLANIFICACIONES_KEY })
      queryClient.invalidateQueries({ queryKey: [...PLANIFICACIONES_KEY, id] })
    },
  })
}

export function useCerrarPlanificacion() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => planificacionesApi.cerrar(id),
    onSuccess: (_data, id) => {
      queryClient.invalidateQueries({ queryKey: PLANIFICACIONES_KEY })
      queryClient.invalidateQueries({ queryKey: [...PLANIFICACIONES_KEY, id] })
    },
  })
}

export function useCreateTurno() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ planificacionId, payload }: { planificacionId: string; payload: CreateTurnoPayload }) =>
      planificacionesApi.createTurno(planificacionId, payload),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: [...PLANIFICACIONES_KEY, variables.planificacionId] })
    },
  })
}

export function useDeleteTurno() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ planificacionId, turnoId }: { planificacionId: string; turnoId: string }) =>
      planificacionesApi.deleteTurno(planificacionId, turnoId),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: [...PLANIFICACIONES_KEY, variables.planificacionId] })
    },
  })
}

export function useStaffingRequirements(id: string) {
  return useQuery({
    queryKey: [...PLANIFICACIONES_KEY, id, 'requirements'],
    queryFn: () => planificacionesApi.getDotacion(id),
    enabled: !!id,
  })
}

export function usePlanLeaves(id: string) {
  return useQuery({
    queryKey: [...PLANIFICACIONES_KEY, id, 'leaves'],
    queryFn: () => planificacionesApi.getLeaves(id),
    enabled: !!id,
  })
}

export function useSectores(id: string) {
  return useQuery({
    queryKey: [...PLANIFICACIONES_KEY, id, 'sectores'],
    queryFn: () => planificacionesApi.getSectores(id),
    enabled: !!id,
  })
}

export function useUpdateSectores() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ planificacionId, payload }: { planificacionId: string; payload: { sectores: string[] } }) =>
      planificacionesApi.updateSectores(planificacionId, payload),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: [...PLANIFICACIONES_KEY, variables.planificacionId, 'sectores'] })
      queryClient.invalidateQueries({ queryKey: [...PLANIFICACIONES_KEY, variables.planificacionId, 'requirements'] })
    },
  })
}

export function useUpdateDotacion() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ planificacionId, payload }: { planificacionId: string; payload: { items: { sector: string; tipo_empleado: string; turno: string; cantidad_minima: number }[] } }) =>
      planificacionesApi.updateDotacion(planificacionId, payload),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: [...PLANIFICACIONES_KEY, variables.planificacionId, 'requirements'] })
    },
  })
}
