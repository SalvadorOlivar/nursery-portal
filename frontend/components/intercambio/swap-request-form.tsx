'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { usePlanificaciones } from '@/features/planificaciones/hooks/use-planificaciones'
import { usePlanificacion } from '@/features/planificaciones/hooks/use-planificaciones'
import { useEmployees } from '@/features/employees/hooks/use-employees'
import { useMe } from '@/features/auth/hooks/use-auth'
import { useCreateSwapRequest } from '@/features/intercambio/hooks/use-intercambio'
import { toast } from 'sonner'
import { ApiError } from '@/lib/api/client'

const turnoLabels: Record<string, string> = {
  MANANA: 'Mañana',
  TARDE: 'Tarde',
  VESPERTINO: 'Vespertino',
  NOCHE: 'Noche',
}

export function SwapRequestForm() {
  const router = useRouter()
  const { data: meData } = useMe()
  const { data: planificacionesData } = usePlanificaciones()
  const { data: employeesData, isPending: loadingEmployees, error: employeesError } = useEmployees()
  const createMutation = useCreateSwapRequest()

  const [selectedPlanifID, setSelectedPlanifID] = useState('')
  const [turnoSolicitanteID, setTurnoSolicitanteID] = useState('')
  const [destinoID, setDestinoID] = useState('')
  const [turnoDestinoID, setTurnoDestinoID] = useState('')

  const { data: planifDetail } = usePlanificacion(selectedPlanifID)

  const planificaciones = planificacionesData?.data ?? []
  const employees = employeesData?.data ?? []
  const userEmployeeID = meData?.user?.employee_id

  const myTurnos = planifDetail?.turnos?.filter(
    (t) => t.empleado_id === userEmployeeID
  ) ?? []

  const targetEmployeeTurnos = planifDetail?.turnos?.filter(
    (t) => t.empleado_id === destinoID
  ) ?? []

  const currentEmployee = employees.find((e) => e.id === userEmployeeID)
  const otherEmployees = employees.filter(
    (e) => e.id !== userEmployeeID && e.activo && e.tipo === currentEmployee?.tipo
  )

  function resetForm() {
    setTurnoSolicitanteID('')
    setDestinoID('')
    setTurnoDestinoID('')
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()

    if (!selectedPlanifID || !turnoSolicitanteID || !destinoID || !turnoDestinoID) {
      toast.error('Completa todos los campos')
      return
    }

    try {
      await createMutation.mutateAsync({
        planificacion_id: selectedPlanifID,
        turno_solicitante_id: turnoSolicitanteID,
        turno_destino_id: turnoDestinoID,
        destino_id: destinoID,
      })
      toast.success('Solicitud de intercambio creada')
      router.push('/intercambio')
    } catch (error) {
      toast.error(error instanceof ApiError ? error.message : 'Error al crear la solicitud')
    }
  }

  return (
    <Card className="max-w-xl mx-auto">
      <CardHeader>
        <CardTitle>Nuevo Intercambio de Turno</CardTitle>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label>Planificación</Label>
            <Select
              value={selectedPlanifID}
              onValueChange={(v) => v && (setSelectedPlanifID(v), resetForm())}
            >
              <SelectTrigger>
                <SelectValue placeholder="Seleccionar planificación" />
              </SelectTrigger>
              <SelectContent>
                  {planificaciones.map((p) => (
                  <SelectItem key={p.id} value={p.id}>
                    Sem {p.semana} / {p.anio} — {p.nombre}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          {selectedPlanifID && (
            <>
              <div className="space-y-2">
                <Label>Mi turno a intercambiar</Label>
                <Select value={turnoSolicitanteID} onValueChange={(v) => v && setTurnoSolicitanteID(v)}>
                  <SelectTrigger>
                    <SelectValue placeholder="Seleccionar mi turno" />
                  </SelectTrigger>
                  <SelectContent>
                    {myTurnos.length === 0 && (
                      <SelectItem value="__none__" disabled>
                        No tienes turnos en esta planificación
                      </SelectItem>
                    )}
                    {myTurnos.map((t) => (
                      <SelectItem key={t.id} value={t.id}>
                        {['LUN','MAR','MIE','JUE','VIE','SAB','DOM'][t.dia_semana - 1]} - {turnoLabels[t.tipo] ?? t.tipo} ({t.sector || 'Sin sector'})
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-2">
                <Label>Empleado destino</Label>
                {loadingEmployees ? (
                  <p className="text-sm text-muted-foreground">Cargando empleados...</p>
                ) : employeesError ? (
                  <p className="text-sm text-destructive">
                    Error al cargar empleados: {employeesError.message}
                  </p>
                ) : employees.length === 0 ? (
                  <p className="text-sm text-muted-foreground">No hay empleados registrados</p>
                ) : !currentEmployee ? (
                  <p className="text-sm text-destructive">
                    No se encontró tu empleado vinculado. Contacta al administrador.
                  </p>
                ) : (
                  <Select value={destinoID} onValueChange={(v) => v && (setDestinoID(v), setTurnoDestinoID(''))}>
                    <SelectTrigger>
                      <SelectValue placeholder="Seleccionar empleado" />
                    </SelectTrigger>
                    <SelectContent>
                      {otherEmployees.length === 0 ? (
                        <SelectItem value="__none__" disabled>
                          No hay empleados activos de tipo {currentEmployee.tipo}
                        </SelectItem>
                      ) : (
                        otherEmployees.map((e) => (
                          <SelectItem key={e.id} value={e.id}>
                            {e.apellido}, {e.nombre}
                          </SelectItem>
                        ))
                      )}
                    </SelectContent>
                  </Select>
                )}
              </div>

              {destinoID && (
                <div className="space-y-2">
                  <Label>Turno del empleado destino</Label>
                  <Select value={turnoDestinoID} onValueChange={(v) => v && setTurnoDestinoID(v)}>
                    <SelectTrigger>
                      <SelectValue placeholder="Seleccionar turno destino" />
                    </SelectTrigger>
                    <SelectContent>
                      {targetEmployeeTurnos.length === 0 && (
                        <SelectItem value="__none__" disabled>
                          El empleado no tiene turnos en esta planificación
                        </SelectItem>
                      )}
                      {targetEmployeeTurnos.map((t) => (
                        <SelectItem key={t.id} value={t.id}>
                          {['LUN','MAR','MIE','JUE','VIE','SAB','DOM'][t.dia_semana - 1]} - {turnoLabels[t.tipo] ?? t.tipo} ({t.sector || 'Sin sector'})
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
              )}
            </>
          )}

          <div className="flex gap-2 pt-4">
            <Button type="submit" disabled={createMutation.isPending}>
              {createMutation.isPending ? 'Creando...' : 'Solicitar intercambio'}
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
