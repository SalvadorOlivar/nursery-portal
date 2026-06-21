'use client'

import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { usePlanificaciones, useDeletePlanificacion } from '@/features/planificaciones/hooks/use-planificaciones'
import { toast } from 'sonner'
import type { Planificacion } from '@/types/planificacion'
import { useMe } from '@/features/auth/hooks/use-auth'

const estadoColors: Record<string, 'default' | 'secondary' | 'outline'> = {
  BORRADOR: 'secondary',
  PUBLICADO: 'default',
  CERRADO: 'outline',
}

const estadoLabels: Record<string, string> = {
  BORRADOR: 'Borrador',
  PUBLICADO: 'Publicado',
  CERRADO: 'Cerrado',
}

function PlanificacionRow({ planificacion, canEdit }: { planificacion: Planificacion; canEdit: boolean }) {
  const deleteMutation = useDeletePlanificacion()

  async function handleDelete() {
    if (!confirm('¿Eliminar esta planificación? Se borrarán todos los turnos asociados.')) return
    try {
      await deleteMutation.mutateAsync(planificacion.id)
      toast.success('Planificación eliminada')
    } catch {
      toast.error('Error al eliminar planificación')
    }
  }

  return (
    <TableRow>
      <TableCell className="font-medium">{planificacion.nombre}</TableCell>
      <TableCell>Semana {planificacion.semana} de {planificacion.anio}</TableCell>
      <TableCell>{planificacion.dias} días</TableCell>
      <TableCell>
        <Badge variant={estadoColors[planificacion.estado] ?? 'secondary'}>
          {estadoLabels[planificacion.estado] ?? planificacion.estado}
        </Badge>
      </TableCell>
      <TableCell className="text-right">
        <div className="flex justify-end gap-2">
          <Link href={`/planificaciones/${planificacion.id}`}>
            <Button variant="outline" size="sm">Ver</Button>
          </Link>
          {canEdit && planificacion.estado === 'BORRADOR' && (
            <Button
              variant="outline"
              size="sm"
              className="text-destructive"
              onClick={handleDelete}
              disabled={deleteMutation.isPending}
            >
              Eliminar
            </Button>
          )}
        </div>
      </TableCell>
    </TableRow>
  )
}

export function PlanificacionList() {
  const { data, isLoading, isError } = usePlanificaciones()
  const { data: meData } = useMe()
  const canEdit = meData?.user.role === 'ADMIN' || meData?.user.role === 'SUPERVISOR'

  if (isLoading) {
    return <div className="text-center py-8 text-muted-foreground">Cargando planificaciones...</div>
  }

  if (isError) {
    return (
      <div className="text-center py-8 text-destructive">
        Error al cargar planificaciones. Verifica que el servidor esté corriendo.
      </div>
    )
  }

  const plans = data?.data ?? []

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between">
        <CardTitle>Planificaciones</CardTitle>
        {canEdit && (
          <Link href="/planificaciones/new">
            <Button>Nueva planificación</Button>
          </Link>
        )}
      </CardHeader>
      <CardContent>
        {plans.length === 0 ? (
          <div className="text-center py-8 text-muted-foreground">
            No hay planificaciones. Crea la primera.
          </div>
        ) : (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Nombre</TableHead>
                <TableHead>Período</TableHead>
                <TableHead>Días</TableHead>
                <TableHead>Estado</TableHead>
                <TableHead className="text-right">Acciones</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {plans.map((plan) => (
                <PlanificacionRow key={plan.id} planificacion={plan} canEdit={canEdit} />
              ))}
            </TableBody>
          </Table>
        )}
      </CardContent>
    </Card>
  )
}
