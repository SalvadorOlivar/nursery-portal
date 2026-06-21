'use client'

import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import {
  usePlanificacion,
  usePublicarPlanificacion,
  useCerrarPlanificacion,
} from '@/features/planificaciones/hooks/use-planificaciones'

import { PlanillaDiaria } from './planilla-diaria'
import { SectorManager } from './sector-manager'
import { toast } from 'sonner'
import type { PlanificacionDetail } from '@/types/planificacion'
import { useMe } from '@/features/auth/hooks/use-auth'

const estadoColors: Record<string, 'default' | 'secondary' | 'outline'> = {
  BORRADOR: 'secondary',
  PUBLICADO: 'default',
  CERRADO: 'outline',
}

interface CalendarioProps {
  planificacionId: string
}

export function PlanificacionCalendario({ planificacionId }: CalendarioProps) {
  const { data: planifData, isLoading: planifLoading, isError: planifError } = usePlanificacion(planificacionId)
  const { data: meData } = useMe()
  const publicarMutation = usePublicarPlanificacion()
  const cerrarMutation = useCerrarPlanificacion()

  const [vista, setVista] = useState<'planilla' | 'configuracion'>('planilla')

  if (planifLoading) {
    return <div className="text-center py-8 text-muted-foreground">Cargando planificación...</div>
  }

  if (planifError || !planifData) {
    return <div className="text-center py-8 text-destructive">Error al cargar planificación</div>
  }

  const planificacion: PlanificacionDetail = planifData
  const employees = planificacion.employees ?? []
  const activeEmployees = employees.filter((e) => e.activo)
  const canEdit = meData?.user.role === 'ADMIN' || meData?.user.role === 'SUPERVISOR'
  const readonly = !canEdit || planificacion.estado !== 'BORRADOR'

  async function handlePublicar() {
    try {
      await publicarMutation.mutateAsync(planificacionId)
      toast.success('Planificación publicada')
    } catch {
      toast.error('Error al publicar planificación')
    }
  }

  async function handleCerrar() {
    try {
      await cerrarMutation.mutateAsync(planificacionId)
      toast.success('Planificación cerrada')
    } catch {
      toast.error('Error al cerrar planificación')
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="space-y-1">
          <h2 className="text-xl font-semibold">{planificacion.nombre}</h2>
          <div className="flex items-center gap-2">
            <Badge variant={estadoColors[planificacion.estado] ?? 'secondary'}>
              {planificacion.estado === 'BORRADOR' ? 'Borrador' :
               planificacion.estado === 'PUBLICADO' ? 'Publicado' : 'Cerrado'}
            </Badge>
            <span className="text-sm text-muted-foreground">
              {activeEmployees.length} empleados activos · 7 días
            </span>
          </div>
        </div>

        <div className="flex gap-2">
          {canEdit && planificacion.estado === 'BORRADOR' && (
            <Button onClick={handlePublicar} disabled={publicarMutation.isPending}>
              Publicar
            </Button>
          )}
          {canEdit && planificacion.estado === 'PUBLICADO' && (
            <Button variant="outline" onClick={handleCerrar} disabled={cerrarMutation.isPending}>
              Cerrar planificación
            </Button>
          )}
        </div>
      </div>

      <div className="flex gap-1 border-b">
        <button
          type="button"
          onClick={() => setVista('planilla')}
          className={`px-4 py-2 text-sm font-medium border-b-2 transition-colors ${
            vista === 'planilla'
              ? 'border-primary text-primary'
              : 'border-transparent text-muted-foreground hover:text-foreground'
          }`}
        >
          Vista planilla
        </button>
        {canEdit && (
          <button
            type="button"
            onClick={() => setVista('configuracion')}
            className={`px-4 py-2 text-sm font-medium border-b-2 transition-colors ${
              vista === 'configuracion'
                ? 'border-primary text-primary'
                : 'border-transparent text-muted-foreground hover:text-foreground'
            }`}
          >
            Configuración
          </button>
        )}
      </div>

      {vista === 'planilla' ? (
        <PlanillaDiaria planificacionId={planificacionId} readonly={readonly} />
      ) : (
        <div className="space-y-6">
          <SectorManager planificacionId={planificacionId} readonly={readonly} />
        </div>
      )}
    </div>
  )
}
