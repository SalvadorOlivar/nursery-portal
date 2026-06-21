'use client'

import { useMemo, useState, useCallback, useRef, useEffect } from 'react'
import { Button } from '@/components/ui/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
} from '@/components/ui/select'
import { usePlanificacion, useSectores, useCreateTurno, useDeleteTurno, usePlanLeaves } from '@/features/planificaciones/hooks/use-planificaciones'
import type { Turno, TipoTurno } from '@/types/planificacion'
import type { Employee } from '@/types/employee'
import type { LeaveRequest } from '@/types/ausencia'

const tipoLabels: Record<string, string> = {
  SUPERVISOR: 'Supervisor/a',
  NURSE: 'Licenciada/o en Enfermería',
  NURSE_ASSISTANT: 'Enfermera/o',
  AUXILIAR_SERVICIO: 'Auxiliar de Servicio',
}

const turnoLabels: Record<string, string> = {
  MANANA: 'Mañana',
  TARDE: 'Tarde',
  VESPERTINO: 'Vespertino',
  NOCHE: 'Noche',
}

const turnoColors: Record<string, string> = {
  MANANA: 'bg-blue-50 border-blue-200',
  TARDE: 'bg-yellow-50 border-yellow-200',
  VESPERTINO: 'bg-orange-50 border-orange-200',
  NOCHE: 'bg-indigo-50 border-indigo-200',
}

const turnoBadgeColors: Record<string, string> = {
  MANANA: 'bg-blue-200 text-blue-900',
  TARDE: 'bg-yellow-200 text-yellow-900',
  VESPERTINO: 'bg-orange-200 text-orange-900',
  NOCHE: 'bg-indigo-300 text-indigo-950',
}

const tipoOrden = ['SUPERVISOR', 'NURSE', 'NURSE_ASSISTANT', 'AUXILIAR_SERVICIO']
const turnosOrden: TipoTurno[] = ['MANANA', 'TARDE', 'VESPERTINO', 'NOCHE']
const emptyTurnos: Turno[] = []
const emptyEmployees: Employee[] = []
const emptySectores: { id: string; nombre: string }[] = []
const emptyLeaves: LeaveRequest[] = []

const leaveTypeLabels: Record<string, string> = {
  VACACIONES: 'Vac',
  ENFERMEDAD: 'Lic',
  PERSONAL: 'Pers',
  DIA_FAVOR: 'Día F',
}

const diaLabels: Record<number, string> = {
  1: 'Lunes',
  2: 'Martes',
  3: 'Miércoles',
  4: 'Jueves',
  5: 'Viernes',
  6: 'Sábado',
  7: 'Domingo',
}

function isoWeekToDate(anio: number, semana: number, diaSemana: number): Date {
  const jan4 = new Date(anio, 0, 4)
  const daysSinceMonday = jan4.getDay() === 0 ? 6 : jan4.getDay() - 1
  const mondayWeek1 = new Date(anio, 0, 4 - daysSinceMonday)
  return new Date(anio, 0, mondayWeek1.getDate() + (semana - 1) * 7 + (diaSemana - 1))
}

interface PlanillaDiariaProps {
  planificacionId: string
  readonly: boolean
}

function EmployeeSelector({
  employees,
  onSelect,
  onClose,
}: {
  employees: Employee[]
  onSelect: (id: string) => void
  onClose: () => void
}) {
  const ref = useRef<HTMLDivElement>(null)

  useEffect(() => {
    function handleClick(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        onClose()
      }
    }
    document.addEventListener('mousedown', handleClick)
    return () => document.removeEventListener('mousedown', handleClick)
  }, [onClose])

  return (
    <div ref={ref} className="absolute z-10 mt-1 bg-white border rounded-md shadow-lg max-h-40 overflow-y-auto min-w-[160px]">
      {employees.length === 0 ? (
        <div className="px-3 py-2 text-xs text-muted-foreground">Sin empleados disponibles</div>
      ) : (
        employees.map((emp) => (
          <button
            key={emp.id}
            type="button"
            onClick={() => onSelect(emp.id)}
            className="w-full text-left px-3 py-1.5 text-xs hover:bg-muted transition-colors truncate"
          >
            {emp.apellido}, {emp.nombre}
          </button>
        ))
      )}
    </div>
  )
}

export function PlanillaDiaria({ planificacionId, readonly }: PlanillaDiariaProps) {
  const { data: planifData } = usePlanificacion(planificacionId)
  const { data: sectoresData } = useSectores(planificacionId)
  const { data: leavesData } = usePlanLeaves(planificacionId)
  const createTurnoMutation = useCreateTurno()
  const deleteTurnoMutation = useDeleteTurno()

  const [dia, setDia] = useState(1)
  const [addingTo, setAddingTo] = useState<{ sector: string; tipo: string; turno: string } | null>(null)

  const employees = planifData?.employees ?? emptyEmployees
  const activeEmployees = employees.filter((e) => e.activo)
  const turnos = planifData?.turnos ?? emptyTurnos

  const planif = planifData

  const leaves = leavesData?.data ?? emptyLeaves
  const empleadosConLicencia = useMemo(() => {
    const set = new Set<string>()
    if (!planif) return set
    for (const lr of leaves) {
      const fechaDia = isoWeekToDate(planif.anio, planif.semana, dia)
      const inicioParts = lr.fecha_inicio.split('-').map(Number)
      const finParts = lr.fecha_fin.split('-').map(Number)
      const fechaInicio = new Date(inicioParts[0], inicioParts[1] - 1, inicioParts[2])
      const fechaFin = new Date(finParts[0], finParts[1] - 1, finParts[2])
      if (lr.estado === 'APROBADO' && fechaDia >= fechaInicio && fechaDia <= fechaFin) {
        set.add(lr.employee_id)
      }
    }
    return set
  }, [leaves, dia, planif])

  const turnosDelDia = useMemo(() => {
    return turnos.filter((t) => t.dia_semana === dia)
  }, [turnos, dia])

  const empleadosAsignadosPorTipo = useMemo(() => {
    const set = new Set<string>()
    for (const t of turnosDelDia) {
      set.add(`${t.empleado_id}|${t.tipo}`)
    }
    return set
  }, [turnosDelDia])

  const employeesByTipo = useMemo(() => {
    const map: Record<string, Employee[]> = {}
    for (const tipo of tipoOrden) {
      map[tipo] = []
    }
    for (const emp of activeEmployees) {
      if (map[emp.tipo]) {
        map[emp.tipo].push(emp)
      }
    }
    return map
  }, [activeEmployees])

  const turnosPorTipoYTurno = useMemo(() => {
    const map: Record<string, Record<string, Turno[]>> = {}
    for (const tipo of tipoOrden) {
      map[tipo] = {}
      for (const turno of turnosOrden) {
        map[tipo][turno] = []
      }
    }
    for (const t of turnosDelDia) {
      const emp = activeEmployees.find((e) => e.id === t.empleado_id)
      if (emp && map[emp.tipo]) {
        map[emp.tipo][t.tipo].push(t)
      }
    }
    return map
  }, [turnosDelDia, activeEmployees])

  const turnosPorTipoSectorYTurno = useMemo(() => {
    const map: Record<string, Record<string, Record<string, Turno[]>>> = {}
    for (const t of turnosDelDia) {
      const emp = activeEmployees.find((e) => e.id === t.empleado_id)
      if (!emp || !tipoOrden.includes(emp.tipo)) continue
      const sec = t.sector || ''
      if (!map[emp.tipo]) map[emp.tipo] = {}
      if (!map[emp.tipo][sec]) {
        map[emp.tipo][sec] = {}
        for (const turno of turnosOrden) {
          map[emp.tipo][sec][turno] = []
        }
      }
      map[emp.tipo][sec][t.tipo].push(t)
    }
    return map
  }, [turnosDelDia, activeEmployees])

  const sectores = sectoresData?.data ?? emptySectores

  const allSectores = useMemo(() => {
    return sectores.map((s) => s.nombre).sort((a, b) => {
      const na = parseInt(a.split('-')[0], 10)
      const nb = parseInt(b.split('-')[0], 10)
      return na - nb
    })
  }, [sectores])

  const tiposPorSector = ['NURSE', 'NURSE_ASSISTANT']

  const handleEmployeeClick = useCallback(async (turno: Turno) => {
    if (readonly) return
    await deleteTurnoMutation.mutateAsync({
      planificacionId,
      turnoId: turno.id,
    })
  }, [planificacionId, readonly, deleteTurnoMutation])

  const handleAddEmployee = useCallback(async (empleadoId: string, sector: string, tipo: string, turno: string) => {
    setAddingTo(null)
    if (readonly) return
    await createTurnoMutation.mutateAsync({
      planificacionId,
      payload: { empleado_id: empleadoId, dia_semana: dia, tipo: turno as TipoTurno, sector },
    })
  }, [dia, planificacionId, readonly, createTurnoMutation])

  const getAvailableEmployees = useCallback((_sector: string, tipo: string, turno: string) => {
    return (employeesByTipo[tipo] ?? []).filter(
      (e) => !empleadosAsignadosPorTipo.has(`${e.id}|${turno}`) && !empleadosConLicencia.has(e.id)
    )
  }, [employeesByTipo, empleadosAsignadosPorTipo, empleadosConLicencia])

  const isLoading = createTurnoMutation.isPending || deleteTurnoMutation.isPending

  function renderCell(
    sector: string,
    tipo: string,
    turno: string,
    emps: Turno[],
  ) {
    const isOpen = addingTo?.sector === sector && addingTo?.tipo === tipo && addingTo?.turno === turno

    return (
      <td
        key={turno}
        className={`px-3 py-2 border-b align-top ${turnoColors[turno]} relative`}
      >
        <div className="flex flex-col gap-1 min-h-[40px]">
          {emps.map((t) => {
            const emp = activeEmployees.find((e) => e.id === t.empleado_id)
            return (
              <button
                key={t.id}
                type="button"
                disabled={isLoading || readonly}
                onClick={() => handleEmployeeClick(t)}
                title={
                  readonly
                    ? `${emp?.apellido}, ${emp?.nombre}`
                    : `${emp?.apellido}, ${emp?.nombre} - Click para cambiar`
                }
                className={`text-[11px] px-1.5 py-0.5 rounded ${
                  turnoBadgeColors[t.tipo]
                } ${
                  readonly ? 'cursor-default' : 'cursor-pointer hover:opacity-80'
                } transition-opacity text-left truncate max-w-[130px]`}
              >
                {emp?.apellido}, {emp?.nombre}
              </button>
            )
          })}
          {!readonly && (
            <div className="relative">
              <button
                type="button"
                onClick={() => setAddingTo(isOpen ? null : { sector, tipo, turno })}
                disabled={isLoading}
                className="text-[11px] w-full py-0.5 rounded border border-dashed border-muted-foreground/30 text-muted-foreground/60 hover:border-muted-foreground/60 hover:text-muted-foreground transition-colors"
              >
                + Agregar
              </button>
              {isOpen && (
                <EmployeeSelector
                  employees={getAvailableEmployees(sector, tipo, turno)}
                  onSelect={(empId) => handleAddEmployee(empId, sector, tipo, turno)}
                  onClose={() => setAddingTo(null)}
                />
              )}
            </div>
          )}
        </div>
      </td>
    )
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center gap-3">
        <span className="text-sm font-medium whitespace-nowrap">Día:</span>
        <Select value={String(dia)} onValueChange={(v) => { if (v) setDia(Number(v)) }}>
          <SelectTrigger className="w-48">
            <span className="flex-1 text-left">{diaLabels[dia] ?? 'Seleccionar día'}</span>
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="1">Lunes</SelectItem>
            <SelectItem value="2">Martes</SelectItem>
            <SelectItem value="3">Miércoles</SelectItem>
            <SelectItem value="4">Jueves</SelectItem>
            <SelectItem value="5">Viernes</SelectItem>
            <SelectItem value="6">Sábado</SelectItem>
            <SelectItem value="7">Domingo</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {empleadosConLicencia.size > 0 && (
        <div className="flex items-center gap-2 text-xs text-muted-foreground bg-red-50 border border-red-200 rounded-md px-3 py-2">
          <span className="font-medium text-red-700">Empleados con licencia este día:</span>
          {Array.from(empleadosConLicencia).map((empId) => {
            const emp = activeEmployees.find((e) => e.id === empId)
            return emp ? (
              <span key={empId} className="inline-flex items-center gap-1 bg-red-100 text-red-800 px-2 py-0.5 rounded">
                {emp.apellido}, {emp.nombre}
              </span>
            ) : null
          })}
        </div>
      )}

      <div className="overflow-x-auto pb-4">
        <table className="w-full text-sm border-collapse">
          <thead>
            <tr>
              <th className="text-left px-3 py-2 font-medium text-muted-foreground border-b w-48">
                Cargo
              </th>
              {turnosOrden.map((turno) => (
                <th
                  key={turno}
                  className={`px-3 py-2 font-medium text-center border-b ${turnoColors[turno]}`}
                >
                  {turnoLabels[turno]}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {(() => {
              const tipo = 'SUPERVISOR'
              const supervisorTurnos = turnosPorTipoYTurno[tipo] ?? {}

              return (
                <tr key={tipo}>
                  <td className="px-3 py-2 font-medium border-b align-top">
                    {tipoLabels[tipo] ?? tipo}
                    <div className="text-[10px] text-muted-foreground font-normal">
                      {(employeesByTipo[tipo] ?? []).length} empleados
                    </div>
                  </td>
                  {turnosOrden.map((turno) =>
                    renderCell('', tipo, turno, supervisorTurnos[turno] ?? [])
                  )}
                </tr>
              )
            })()}
            {allSectores.flatMap((sec) => [
              <tr key={`sec-${sec}`}>
                <td className="px-3 py-2 font-medium border-b border-t-2 bg-muted/20" colSpan={5}>
                  Sector {sec}
                </td>
              </tr>,
              ...tiposPorSector.map((tipo) => {
                const sectorTurnos = turnosPorTipoSectorYTurno[tipo]?.[sec] ?? {}
                const totalTipo = (employeesByTipo[tipo] ?? []).length

                return (
                  <tr key={`${sec}-${tipo}`}>
                    <td className="px-3 py-2 text-sm border-b align-top pl-6">
                      {tipoLabels[tipo] ?? tipo}
                      <div className="text-[10px] text-muted-foreground font-normal">
                        {totalTipo} empleados
                      </div>
                    </td>
                    {turnosOrden.map((turno) =>
                      renderCell(sec, tipo, turno, sectorTurnos[turno] ?? [])
                    )}
                  </tr>
                )
              }),
            ])}
            {(() => {
              const tipo = 'AUXILIAR_SERVICIO'
              const auxTurnos = turnosPorTipoYTurno[tipo] ?? {}

              return (
                <tr key={tipo}>
                  <td className="px-3 py-2 font-medium border-b align-top">
                    {tipoLabels[tipo] ?? tipo}
                    <div className="text-[10px] text-muted-foreground font-normal">
                      {(employeesByTipo[tipo] ?? []).length} empleados
                    </div>
                  </td>
                  {turnosOrden.map((turno) =>
                    renderCell('', tipo, turno, auxTurnos[turno] ?? [])
                  )}
                </tr>
              )
            })()}
          </tbody>
        </table>
      </div>

      {!readonly && (
        <div className="flex items-center gap-4 text-xs text-muted-foreground">
          <span>Click en un empleado para eliminar su turno &middot; + Agregar para asignar un nuevo turno</span>
        </div>
      )}
    </div>
  )
}