'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Textarea } from '@/components/ui/textarea'
import { useEmployees } from '@/features/employees/hooks/use-employees'
import { useCreateLeaveRequest } from '@/features/ausencia/hooks/use-ausencia'
import { useMe } from '@/features/auth/hooks/use-auth'
import { toast } from 'sonner'
import type { LeaveType } from '@/types/ausencia'

const leaveTypes: { value: LeaveType; label: string }[] = [
  { value: 'VACACIONES', label: 'Vacaciones' },
  { value: 'ENFERMEDAD', label: 'Enfermedad' },
  { value: 'PERSONAL', label: 'Personal' },
]

export function LeaveRequestForm() {
  const router = useRouter()
  const { data: meData } = useMe()
  const { data: employeesData } = useEmployees()
  const createMutation = useCreateLeaveRequest()
  const user = meData?.user
  const isAdmin = user?.role === 'ADMIN' || user?.role === 'SUPERVISOR'
  const employees = employeesData?.data ?? []

  const [employeeId, setEmployeeId] = useState('')
  const [fechaInicio, setFechaInicio] = useState('')
  const [fechaFin, setFechaFin] = useState('')
  const [tipo, setTipo] = useState<LeaveType>('VACACIONES')
  const [motivo, setMotivo] = useState('')

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()

    if (!fechaInicio || !fechaFin) {
      toast.error('Completa todas las fechas')
      return
    }

    const targetEmployeeId = isAdmin ? employeeId : (user?.employee_id ?? '')

    if (!targetEmployeeId) {
      toast.error('Selecciona un empleado')
      return
    }

    try {
      await createMutation.mutateAsync({
        employee_id: targetEmployeeId,
        fecha_inicio: fechaInicio,
        fecha_fin: fechaFin,
        tipo,
        motivo,
      })
      toast.success('Solicitud de licencia creada')
      router.push('/leave-requests')
    } catch (err: any) {
      toast.error(err?.message || 'Error al crear licencia')
    }
  }

  return (
    <Card className="max-w-lg mx-auto">
      <CardHeader>
        <CardTitle>Nueva solicitud de licencia</CardTitle>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          {isAdmin && (
            <div className="space-y-2">
              <Label htmlFor="employee">Empleado</Label>
              <Select value={employeeId} onValueChange={(v) => v && setEmployeeId(v)}>
                <SelectTrigger id="employee">
                  <SelectValue placeholder="Seleccionar empleado" />
                </SelectTrigger>
                <SelectContent>
                  {employees.filter(e => e.activo).map((emp) => (
                    <SelectItem key={emp.id} value={emp.id}>
                      {emp.apellido}, {emp.nombre}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          )}

          <div className="space-y-2">
            <Label htmlFor="tipo">Tipo de licencia</Label>
            <Select value={tipo} onValueChange={(v) => v && setTipo(v as LeaveType)}>
              <SelectTrigger id="tipo">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {leaveTypes.map((lt) => (
                  <SelectItem key={lt.value} value={lt.value}>
                    {lt.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="fecha_inicio">Fecha inicio</Label>
              <Input
                id="fecha_inicio"
                type="date"
                value={fechaInicio}
                onChange={(e) => setFechaInicio(e.target.value)}
                required
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="fecha_fin">Fecha fin</Label>
              <Input
                id="fecha_fin"
                type="date"
                value={fechaFin}
                onChange={(e) => setFechaFin(e.target.value)}
                required
              />
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="motivo">Motivo</Label>
            <Textarea
              id="motivo"
              value={motivo}
              onChange={(e) => setMotivo(e.target.value)}
              placeholder="Opcional: describí el motivo de la licencia"
              rows={3}
            />
          </div>

          <div className="flex gap-2">
            <Button type="submit" disabled={createMutation.isPending}>
              {createMutation.isPending ? 'Guardando...' : 'Solicitar'}
            </Button>
            <Button type="button" variant="outline" onClick={() => router.back()}>
              Cancelar
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  )
}
