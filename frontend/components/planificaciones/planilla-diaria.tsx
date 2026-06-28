'use client'

import { useMemo, useState, useCallback, useRef, useEffect } from 'react'
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

const tipoShortLabels: Record<string, string> = {
  SUPERVISOR: 'SUP',
  NURSE: 'ENF',
  NURSE_ASSISTANT: 'AUX',
  AUXILIAR_SERVICIO: 'APO',
}

const turnoLabels: Record<string, string> = {
  MANANA: 'Mañana',
  TARDE: 'Tarde',
  VESPERTINO: 'Vespertino',
  NOCHE: 'Noche',
}

const turnoTime: Record<string, string> = {
  MANANA: '06:00-12:00',
  TARDE: '12:00-18:00',
  VESPERTINO: '18:00-00:00',
  NOCHE: '00:00-06:00',
}

const tipoOrden = ['SUPERVISOR', 'NURSE', 'NURSE_ASSISTANT', 'AUXILIAR_SERVICIO']
const turnosOrden: TipoTurno[] = ['MANANA', 'TARDE', 'VESPERTINO', 'NOCHE']
const tiposPorSector = ['NURSE', 'NURSE_ASSISTANT'] as const
const emptyTurnos: Turno[] = []
const emptyEmployees: Employee[] = []
const emptySectores: { id: string; nombre: string }[] = []
const emptyLeaves: LeaveRequest[] = []

const diaShortLabels: Record<number, string> = {
  1: 'Lun', 2: 'Mar', 3: 'Mié', 4: 'Jue', 5: 'Vie', 6: 'Sáb', 7: 'Dom',
}

import { isoWeekToDate } from '@/lib/utils'

const monthsEs = ['enero', 'febrero', 'marzo', 'abril', 'mayo', 'junio', 'julio', 'agosto', 'septiembre', 'octubre', 'noviembre', 'diciembre']

function initials(nombre?: string, apellido?: string): string {
  return ((nombre?.[0] ?? '') + (apellido?.[0] ?? '')).toUpperCase() || '??'
}

function fullName(emp?: { nombre?: string; apellido?: string } | null): string {
  if (!emp) return 'Sin empleado'
  return `${emp.apellido ?? ''}, ${emp.nombre ?? ''}`.replace(/^,\s/, '').replace(/,\s$/, '') || 'Sin nombre'
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

  const weekDates = useMemo(() => {
    if (!planif) return []
    return [1, 2, 3, 4, 5, 6, 7].map((d) => isoWeekToDate(planif.anio, planif.semana, d))
  }, [planif])

  const dateRangeStr = useMemo(() => {
    if (weekDates.length < 7) return ''
    const start = weekDates[0]
    const end = weekDates[6]
    return `Semana del ${start.getDate()} al ${end.getDate()} de ${monthsEs[end.getMonth()]} de ${end.getFullYear()}`
  }, [weekDates])

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

  const isLoading = createTurnoMutation.isPending || deleteTurnoMutation.isPending

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

  function renderCell(
    sector: string,
    tipo: string,
    turno: string,
    emps: Turno[],
  ) {
    const isOpen = addingTo?.sector === sector && addingTo?.tipo === tipo && addingTo?.turno === turno

    return (
      <td key={turno} className="p-2 align-top border border-[var(--border)]">
        <div className="flex flex-col gap-1">
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
                    ? fullName(emp)
                    : `${fullName(emp)} - Click para eliminar`
                }
                className={`inline-flex w-full items-center gap-[6px] px-[10px] py-[5px] rounded-[6px] text-[0.78rem] font-[450] bg-[var(--surface)] border border-[var(--border)] transition-all duration-200 ${
                  readonly ? 'cursor-default' : 'cursor-pointer hover:border-[var(--accent)] hover:shadow-[var(--shadow-sm)] hover:-translate-y-[0.5px]'
                }`}
              >
                <span className="w-[22px] h-[22px] rounded-full bg-[var(--accent-light)] text-[var(--accent-dark)] flex items-center justify-center text-[0.65rem] font-semibold shrink-0">
                  {initials(emp?.nombre, emp?.apellido)}
                </span>
                <span className="min-w-0 flex-1 truncate text-left">{fullName(emp)}</span>
                <span className="text-[0.62rem] font-[510] text-[var(--muted-foreground)] ml-auto tracking-[0.03em] opacity-65 shrink-0">
                  {tipoShortLabels[tipo] ?? tipo}
                </span>
              </button>
            )
          })}
          {!readonly && (
            <div className="relative">
              <button
                type="button"
                onClick={() => setAddingTo(isOpen ? null : { sector, tipo, turno })}
                disabled={isLoading}
                className="flex w-full items-center justify-center gap-[4px] py-[5px] rounded-[6px] border border-dashed border-[var(--border)] bg-transparent text-[var(--muted-foreground)] text-[0.75rem] cursor-pointer transition-all duration-200 hover:border-[var(--accent)] hover:text-[var(--accent)] hover:bg-[var(--accent-light)]"
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
    <div className="space-y-6">
      {planif && dateRangeStr && (
        <p className="text-sm text-[var(--muted-foreground)]">
          {dateRangeStr} &middot; Planificación: {planif.nombre}
        </p>
      )}

      <div className="flex flex-wrap gap-[5px]" role="group" aria-label="Selector de día">
        {[1, 2, 3, 4, 5, 6, 7].map((d) => {
          const date = weekDates[d - 1]
          return (
            <button
              key={d}
              type="button"
              onClick={() => setDia(d)}
              className={`px-4 py-[5px] rounded-full text-[0.78rem] font-[450] border transition-all duration-200 cursor-pointer ${
                dia === d
                  ? 'border-[var(--accent)] bg-[var(--accent)] text-white font-[510]'
                  : 'border-[var(--border)] bg-[var(--surface)] text-[var(--muted-foreground)] hover:border-[var(--accent)] hover:text-[var(--accent)]'
              }`}
            >
              {diaShortLabels[d]}{date ? ` ${date.getDate()}` : ''}
            </button>
          )
        })}
      </div>

      {empleadosConLicencia.size > 0 && (
        <div className="flex items-center gap-3 text-sm text-[var(--danger)] bg-red-50 border border-red-200 rounded-[var(--radius-sm)] px-4 py-3">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="w-[18px] h-[18px] shrink-0">
            <circle cx="12" cy="12" r="10" /><line x1="12" y1="8" x2="12" y2="12" /><line x1="12" y1="16" x2="12.01" y2="16" />
          </svg>
          <span>Personal con licencia hoy: </span>
          {Array.from(empleadosConLicencia).map((empId) => {
            const emp = activeEmployees.find((e) => e.id === empId)
            return emp ? (
              <span key={empId} className="inline-flex px-[10px] py-[3px] rounded-full bg-red-100 text-red-800 text-[0.75rem] font-[510]">
                {emp.apellido}, {emp.nombre}
              </span>
            ) : null
          })}
        </div>
      )}

      <div className="overflow-x-auto pb-2">
        <table className="w-full text-sm border-collapse">
          <thead>
            <tr>
              <th className="text-left px-3 py-2 font-medium text-[var(--muted-foreground)] border-b border-[var(--border)] w-48">
                Cargo
              </th>
              {turnosOrden.map((turno) => (
                <th key={turno} className="px-3 py-2 font-medium text-center border-b border-[var(--border)] text-[var(--fg)]">
                  <div>{turnoLabels[turno]}</div>
                  <div className="text-[0.65rem] font-normal text-[var(--muted-foreground)]">{turnoTime[turno]}</div>
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {(() => {
              const tipo = 'SUPERVISOR'
              const supervisorTurnos = turnosPorTipoSectorYTurno[tipo]?.[''] ?? {}
              const totalTipo = (employeesByTipo[tipo] ?? []).length
              return (
                <tr key={tipo}>
                  <td className="px-3 py-2 font-medium border border-[var(--border)] align-top">
                    {tipoLabels[tipo]}
                    <div className="text-[0.65rem] font-normal text-[var(--muted-foreground)]">{totalTipo} empleados</div>
                  </td>
                  {turnosOrden.map((turno) =>
                    renderCell('', tipo, turno, supervisorTurnos[turno] ?? [])
                  )}
                </tr>
              )
            })()}
            {allSectores.flatMap((sec) => [
              <tr key={`sec-${sec}`}>
                <td className="px-3 py-2 font-medium border border-[var(--border)] bg-[var(--bg)] text-sm" colSpan={5}>
                  <span className="inline-block border-b-2 border-[var(--accent)] pb-0.5 text-[0.82rem] font-[510] uppercase tracking-[0.06em]">
                    {sec}
                  </span>
                </td>
              </tr>,
              ...tiposPorSector.map((tipo) => {
                const sectorTurnos = turnosPorTipoSectorYTurno[tipo]?.[sec] ?? {}
                const totalTipo = (employeesByTipo[tipo] ?? []).length
                return (
                  <tr key={`${sec}-${tipo}`}>
                    <td className="px-3 py-2 text-sm border border-[var(--border)] align-top pl-6">
                      {tipoLabels[tipo]}
                      <div className="text-[0.65rem] font-normal text-[var(--muted-foreground)]">{totalTipo} empleados</div>
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
              const auxTurnos = turnosPorTipoSectorYTurno[tipo]?.[''] ?? {}
              const totalTipo = (employeesByTipo[tipo] ?? []).length
              return (
                <tr key={tipo}>
                  <td className="px-3 py-2 font-medium border border-[var(--border)] align-top">
                    {tipoLabels[tipo]}
                    <div className="text-[0.65rem] font-normal text-[var(--muted-foreground)]">{totalTipo} empleados</div>
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
        <div className="flex items-center gap-4 text-xs text-[var(--muted-foreground)]">
          <span>Click en un empleado para eliminar su turno &middot; + Agregar para asignar un nuevo turno</span>
        </div>
      )}
    </div>
  )
}
