'use client'

import { useMemo } from 'react'
import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { useLeaveRequests, useApproveLeaveRequest, useRejectLeaveRequest } from '@/features/ausencia/hooks/use-ausencia'
import { useEmployees } from '@/features/employees/hooks/use-employees'
import { useMe } from '@/features/auth/hooks/use-auth'
import { toast } from 'sonner'
import type { LeaveRequest } from '@/types/ausencia'

const tipoLabels: Record<string, string> = {
  VACACIONES: 'Vacaciones',
  ENFERMEDAD: 'Enfermedad',
  PERSONAL: 'Personal',
  DIA_FAVOR: 'Día a favor',
}

const estadoColors: Record<string, 'default' | 'secondary' | 'outline' | 'destructive'> = {
  PENDIENTE: 'secondary',
  APROBADO: 'default',
  RECHAZADO: 'destructive',
}

function LeaveRequestRow({ lr, canApprove }: { lr: LeaveRequest; canApprove: boolean }) {
  const approveMutation = useApproveLeaveRequest()
  const rejectMutation = useRejectLeaveRequest()
  const { data: employeesData } = useEmployees()
  const employees = employeesData?.data ?? []
  const emp = employees.find((e) => e.id === lr.employee_id)

  async function handleApprove() {
    try {
      await approveMutation.mutateAsync(lr.id)
      toast.success('Licencia aprobada')
    } catch {
      toast.error('Error al aprobar licencia')
    }
  }

  async function handleReject() {
    try {
      await rejectMutation.mutateAsync(lr.id)
      toast.success('Licencia rechazada')
    } catch {
      toast.error('Error al rechazar licencia')
    }
  }

  return (
    <TableRow>
      <TableCell className="font-medium">
        {emp ? `${emp.apellido}, ${emp.nombre}` : lr.employee_id}
      </TableCell>
      <TableCell>{lr.fecha_inicio}</TableCell>
      <TableCell>{lr.fecha_fin}</TableCell>
      <TableCell>
        <Badge variant={estadoColors[lr.estado] ?? 'secondary'}>
          {tipoLabels[lr.tipo] ?? lr.tipo}
        </Badge>
      </TableCell>
      <TableCell>
        <Badge variant={estadoColors[lr.estado] ?? 'secondary'}>
          {lr.estado}
        </Badge>
      </TableCell>
      <TableCell className="text-xs text-muted-foreground max-w-[200px] truncate">
        {lr.motivo}
      </TableCell>
      <TableCell className="text-right">
        {canApprove && lr.estado === 'PENDIENTE' && (
          <div className="flex justify-end gap-1">
            <Button
              variant="outline"
              size="sm"
              className="text-green-600"
              onClick={handleApprove}
              disabled={approveMutation.isPending}
            >
              Aprobar
            </Button>
            <Button
              variant="outline"
              size="sm"
              className="text-destructive"
              onClick={handleReject}
              disabled={rejectMutation.isPending}
            >
              Rechazar
            </Button>
          </div>
        )}
      </TableCell>
    </TableRow>
  )
}

export function LeaveRequestList() {
  const { data: meData } = useMe()
  const { data, isLoading, isError } = useLeaveRequests()
  const canApprove = meData?.user.role === 'ADMIN' || meData?.user.role === 'SUPERVISOR'

  if (isLoading) {
    return <div className="text-center py-8 text-muted-foreground">Cargando licencias...</div>
  }

  if (isError) {
    return (
      <div className="text-center py-8 text-destructive">
        Error al cargar licencias.
      </div>
    )
  }

  const requests = data?.data ?? []

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between">
        <CardTitle>Licencias y Ausencias</CardTitle>
        <Link href="/leave-requests/new">
          <Button>Nueva licencia</Button>
        </Link>
      </CardHeader>
      <CardContent>
        {requests.length === 0 ? (
          <div className="text-center py-8 text-muted-foreground">
            No hay solicitudes de licencia.
          </div>
        ) : (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Empleado</TableHead>
                <TableHead>Desde</TableHead>
                <TableHead>Hasta</TableHead>
                <TableHead>Tipo</TableHead>
                <TableHead>Estado</TableHead>
                <TableHead>Motivo</TableHead>
                <TableHead className="text-right">Acciones</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {requests.map((lr) => (
                <LeaveRequestRow key={lr.id} lr={lr} canApprove={canApprove} />
              ))}
            </TableBody>
          </Table>
        )}
      </CardContent>
    </Card>
  )
}
