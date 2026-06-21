'use client'

import { useEmployee, useDeactivateEmployee } from '@/features/employees/hooks/use-employees'
import { useRouter } from 'next/navigation'
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
import { toast } from 'sonner'
import { EmployeeForm } from './employee-form'
import { useMe } from '@/features/auth/hooks/use-auth'
import { useLeaveRequests } from '@/features/ausencia/hooks/use-ausencia'
import { useCompensatoryDays } from '@/features/ausencia/hooks/use-ausencia'

const tipoLabels: Record<string, string> = {
  SUPERVISOR: 'Supervisor/a',
  NURSE: 'Licenciada/o en Enfermería',
  NURSE_ASSISTANT: 'Enfermera/o',
  AUXILIAR_SERVICIO: 'Auxiliar de Servicio',
}

const leaveTypeLabels: Record<string, string> = {
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

const motivoLabels: Record<string, string> = {
  DOBLE_TURNO: 'Doble turno',
  DESCANSO_LABORADO: 'Descanso laborado',
}

export function EmployeeDetail({ id }: { id: string }) {
  const router = useRouter()
  const { data: employee, isLoading, isError } = useEmployee(id)
  const { data: meData } = useMe()
  const deactivateMutation = useDeactivateEmployee()
  const { data: leaveData } = useLeaveRequests(id)
  const { data: compDaysData } = useCompensatoryDays(id)
  const canEdit = meData?.user.role === 'ADMIN'

  if (isLoading) {
    return <div className="text-center py-8 text-muted-foreground">Cargando...</div>
  }

  if (isError || !employee) {
    return (
      <div className="text-center py-8">
        <p className="text-destructive mb-4">Empleado no encontrado</p>
        <Button onClick={() => router.push('/employees')}>Volver</Button>
      </div>
    )
  }

  async function handleDeactivate() {
    try {
      await deactivateMutation.mutateAsync(id)
      toast.success('Empleado desactivado correctamente')
      router.refresh()
    } catch {
      toast.error('Error al desactivar empleado')
    }
  }

  const leaveRequests = leaveData?.data ?? []
  const compDays = compDaysData?.data ?? []
  const availableCompDays = compDaysData?.available_count ?? 0

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">
          {employee.apellido}, {employee.nombre}
        </h1>
        <div className="flex gap-2">
          <Button variant="outline" onClick={() => router.push('/employees')}>
            Volver
          </Button>
          {canEdit && employee.activo && (
            <Button
              variant="outline"
              className="text-destructive"
              onClick={handleDeactivate}
              disabled={deactivateMutation.isPending}
            >
              Desactivar
            </Button>
          )}
        </div>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Información general</CardTitle>
        </CardHeader>
        <CardContent>
          <dl className="grid grid-cols-2 gap-4 text-sm">
            <div>
              <dt className="text-muted-foreground">Tipo</dt>
              <dd className="font-medium">{tipoLabels[employee.tipo]}</dd>
            </div>
            <div>
              <dt className="text-muted-foreground">Estado</dt>
              <dd>
                {employee.activo ? (
                  <Badge variant="default" className="bg-green-600">Activo</Badge>
                ) : (
                  <Badge variant="secondary">Inactivo</Badge>
                )}
              </dd>
            </div>
            <div>
              <dt className="text-muted-foreground">Horas mínimas</dt>
              <dd className="font-medium">{employee.horas_minimas}h / mes</dd>
            </div>
            <div>
              <dt className="text-muted-foreground">Horas máximas</dt>
              <dd className="font-medium">{employee.horas_maximas}h / mes</dd>
            </div>
            <div>
              <dt className="text-muted-foreground">Patrón de trabajo</dt>
              <dd className="font-medium">{employee.work_days}x{employee.rest_days}</dd>
            </div>
            <div>
              <dt className="text-muted-foreground">Días a favor disponibles</dt>
              <dd className="font-medium">{availableCompDays}</dd>
            </div>
          </dl>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Licencias y ausencias</CardTitle>
        </CardHeader>
        <CardContent>
          {leaveRequests.length === 0 ? (
            <div className="text-center py-4 text-sm text-muted-foreground">
              Sin solicitudes de licencia
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Desde</TableHead>
                  <TableHead>Hasta</TableHead>
                  <TableHead>Tipo</TableHead>
                  <TableHead>Estado</TableHead>
                  <TableHead>Motivo</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {leaveRequests.slice(0, 10).map((lr) => (
                  <TableRow key={lr.id}>
                    <TableCell>{lr.fecha_inicio}</TableCell>
                    <TableCell>{lr.fecha_fin}</TableCell>
                    <TableCell>
                      <Badge variant="outline">{leaveTypeLabels[lr.tipo] ?? lr.tipo}</Badge>
                    </TableCell>
                    <TableCell>
                      <Badge variant={estadoColors[lr.estado] ?? 'secondary'}>{lr.estado}</Badge>
                    </TableCell>
                    <TableCell className="text-xs text-muted-foreground max-w-[200px] truncate">
                      {lr.motivo}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Días a favor ({availableCompDays} disponibles)</CardTitle>
        </CardHeader>
        <CardContent>
          {compDays.length === 0 ? (
            <div className="text-center py-4 text-sm text-muted-foreground">
              Sin días a favor registrados
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Fecha origen</TableHead>
                  <TableHead>Motivo</TableHead>
                  <TableHead>Descripción</TableHead>
                  <TableHead>Estado</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {compDays.map((cd) => (
                  <TableRow key={cd.id}>
                    <TableCell>{cd.fecha_origen}</TableCell>
                    <TableCell>{motivoLabels[cd.motivo] ?? cd.motivo}</TableCell>
                    <TableCell className="text-xs text-muted-foreground max-w-[250px] truncate">
                      {cd.descripcion}
                    </TableCell>
                    <TableCell>
                      {cd.utilizado ? (
                        <Badge variant="secondary">Usado ({cd.fecha_uso})</Badge>
                      ) : (
                        <Badge variant="default" className="bg-green-600">Disponible</Badge>
                      )}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>

      {canEdit && (
        <Card>
          <CardHeader>
            <CardTitle>Editar empleado</CardTitle>
          </CardHeader>
          <CardContent>
            <EmployeeForm employee={employee} />
          </CardContent>
        </Card>
      )}
    </div>
  )
}
